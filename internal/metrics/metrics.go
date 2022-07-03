package metrics

import (
	"context"
	"strconv"
	"time"

	"github.com/kpetremann/salt-exporter/pkg/events"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/rs/zerolog/log"
)

func ExposeMetrics(ctx context.Context, eventChan chan events.SaltEvent) {
	newJobCounter := promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "salt_new_job_total",
			Help: "Total number of new job processed",
		},
		[]string{"function", "state", "success"},
	)
	responsesCounter := promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "salt_responses_total",
			Help: "Total number of response job processed",
		},
		[]string{"minion", "success"},
	)
	functionResponsesCounter := promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "salt_function_responses_total",
			Help: "Total number of response per function processed",
		},
		[]string{"function", "state", "success"},
	)

	scheduledJobReturnCounter := promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "salt_scheduled_job_return_total",
			Help: "Total number of scheduled job response",
		},
		[]string{"function", "state", "success"},
	)
	expectedResponsesNumber := promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "salt_expected_responses_total",
			Help: "Total number of expected minions responses",
		},
		[]string{"function", "state"},
	)

	for {
		select {
		case <-ctx.Done():
			log.Info().Msg("stopping event listener")
			return
		case event := <-eventChan:
			start := time.Now()

			switch event.Type {
			case "new":
				state := event.ExtractState()
				newJobCounter.WithLabelValues(event.Data.Fun, state, strconv.FormatBool(event.Data.Success)).Inc()
				expectedResponsesNumber.WithLabelValues(event.Data.Fun, state).Add(float64(event.TargetNumber))
			case "ret":
				state := event.ExtractState()
				if event.IsScheduleJob {
					scheduledJobReturnCounter.WithLabelValues(
						event.Data.Fun,
						state,
						strconv.FormatBool(event.Data.Success),
					).Inc()
				} else {
					sucess := strconv.FormatBool(event.Data.Success)
					responsesCounter.WithLabelValues(
						event.Data.Id,
						sucess,
					).Inc()
					functionResponsesCounter.WithLabelValues(
						event.Data.Fun,
						state,
						sucess,
					).Inc()
				}
			}

			elapsed := time.Since(start)
			log.Debug().Str("metric conversion took", elapsed.String()).Send()
		}
	}
}
