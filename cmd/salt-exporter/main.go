package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
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

	listenAddress := flag.String("host", "", "listen address")
	listenPort := flag.Int("port", 2112, "listen port")
	tlsEnabled := flag.Bool("tls", false, "enable TLS")
	tlsCert := flag.String("tls-cert", "", "TLS certificated")
	tlsKey := flag.String("tls-key", "", "TLS private key")
	healthMinions := flag.Bool("health-minions", true, "Enable minion metrics")
	healthFunctionsFilters := flag.String("health-functions-filter", "state.highstate",
		"Apply filter on functions to monitor, separated by a comma")
	healthStatesFilters := flag.String("health-states-filter", "highstate",
		"Apply filter on states to monitor, separated by a comma")
	ignoreTest := flag.Bool("ignore-test", false, "ignore test=True events")
	ignoreMock := flag.Bool("ignore-mock", false, "ignore mock=True events")
	logLevel := flag.String("log-level", "info", "log level (debug, info, warn, error, fatal, panic, disabled)")
	flag.Parse()

	logging.SetLevel(*logLevel)

	if *tlsEnabled {
		missingFlag := false
		if *tlsCert == "" {
			missingFlag = true
			log.Error().Msg("TLS certificate not specified")
		}
		if *tlsCert == "" {
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
		HealthMinions:          *healthMinions,
		HealthFunctionsFilters: strings.Split(*healthFunctionsFilters, ","),
		HealthStatesFilters:    strings.Split(*healthStatesFilters, ","),
		IgnoreTest:             *ignoreTest,
		IgnoreMock:             *ignoreMock,
	}

	if metricsConfig.HealthMinions {
		log.Info().Msg("health-minions: metrics are enabled")
		log.Info().Msgf("health-minions: functions filters: %s", *healthFunctionsFilters)
		log.Info().Msgf("health-minions: states filters: %s", *healthStatesFilters)
	}

	if metricsConfig.IgnoreTest {
		log.Info().Msg("test=True events will be ignored")
		log.Info().Msg("mock=True events will be ignored")
	}

	listenSocket := fmt.Sprint(*listenAddress, ":", *listenPort)

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

		if !*tlsEnabled {
			err = httpServer.ListenAndServe()
		} else {
			err = httpServer.ListenAndServeTLS(*tlsCert, *tlsKey)
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
