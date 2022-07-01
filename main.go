package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/kpetremann/salt-exporter/internal/logging"
	"github.com/kpetremann/salt-exporter/internal/metrics"
	"github.com/kpetremann/salt-exporter/pkg/events"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog/log"
)

func quit() {
	log.Warn().Msg("Bye.")
}

func main() {
	defer quit()
	logging.ConfigureLogging()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	log.Info().Msg("listening for events...")
	eventChan := make(chan events.SaltEvent)

	// listen and expose metric
	eventListener := events.NewEventListener(ctx, eventChan)

	go eventListener.ListenEvents()
	go metrics.ExposeMetrics(ctx, eventChan)

	// start http server
	log.Info().Msg("exposing metrics on :2112/metrics")

	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())
	httpServer := http.Server{Addr: ":2112", Handler: mux}

	go httpServer.ListenAndServe()

	// exiting
	<-ctx.Done()
	httpServer.Shutdown(context.Background())
}
