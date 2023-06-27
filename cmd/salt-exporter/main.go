package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/kpetremann/salt-exporter/internal/logging"
	"github.com/kpetremann/salt-exporter/internal/metrics"
	"github.com/kpetremann/salt-exporter/pkg/event"
	"github.com/kpetremann/salt-exporter/pkg/listener"
	"github.com/kpetremann/salt-exporter/pkg/parser"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog/log"
)

var (
	version = ""
	commit  = ""
	date    = "unknown"
)

func quit() {
	log.Warn().Msg("Bye.")
}

func main() {
	defer quit()
	logging.Configure()

	config, err := ReadConfig()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load settings during initialization")
	}
	log.Fatal().Msg("DEBUG")
	logging.SetLevel(config.LogLevel)

	if config.TLS.Enabled {
		missingFlag := false
		if config.TLS.Key == "" {
			missingFlag = true
			log.Error().Msg("TLS certificate not specified")
		}
		if config.TLS.Certificate == "" {
			missingFlag = true
			log.Error().Msg("TLS private key not specified")
		}
		if missingFlag {
			return
		}
	}

	log.Info().Str("Version", version).Send()
	log.Info().Str("Commit", commit).Send()
	log.Info().Str("Build time", date).Send()

	metricsConfig := metrics.MetricsConfig{
		HealthMinions:          config.HealthMinions,
		HealthFunctionsFilters: config.HealthFunctionsFilter,
		HealthStatesFilters:    config.HealthStatesFilter,
		IgnoreTest:             config.Metrics.Global.Filters.IgnoreTest,
		IgnoreMock:             config.Metrics.Global.Filters.IgnoreMock,
	}

	if metricsConfig.HealthMinions {
		log.Info().Msg("health-minions: metrics are enabled")
		log.Info().Msgf("health-minions: functions filters: %s", config.HealthFunctionsFilter)
		log.Info().Msgf("health-minions: states filters: %s", config.HealthStatesFilter)
	}

	if metricsConfig.IgnoreTest {
		log.Info().Msg("test=True events will be ignored")
		log.Info().Msg("mock=True events will be ignored")
	}

	listenSocket := fmt.Sprint(config.ListenAddress, ":", config.ListenPort)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	log.Info().Msg("listening for events...")
	eventChan := make(chan event.SaltEvent)

	// listen and expose metric
	parser := parser.NewEventParser(false)
	eventListener := listener.NewEventListener(ctx, parser, eventChan)

	go eventListener.ListenEvents()
	go metrics.ExposeMetrics(ctx, eventChan, metricsConfig)

	// start http server
	log.Info().Msg("exposing metrics on " + listenSocket + "/metrics")

	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())
	httpServer := http.Server{Addr: listenSocket, Handler: mux}

	go func() {
		var err error

		if !config.TLS.Enabled {
			err = httpServer.ListenAndServe()
		} else {
			err = httpServer.ListenAndServeTLS(config.TLS.Certificate, config.TLS.Key)
		}

		if err != nil {
			log.Error().Err(err).Send()
			stop()
		}
	}()

	// exiting
	<-ctx.Done()
	if err := httpServer.Shutdown(context.Background()); err != nil {
		log.Error().Err(err).Send()
	}
}
