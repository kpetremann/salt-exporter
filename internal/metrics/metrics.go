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

func boolToFloat64(b bool) float64 {
	if b {
		return 1.0
	}
	return 0.0
}

// Function to check if a string exists in a slice of strings
func contains(slice []string, str string) bool {
	for _, s := range slice {
		if s == str {
			return true
		}
	}
	return false
}

func ExposeMetrics(ctx context.Context, eventChan <-chan events.SaltEvent, metricsConfig MetricsConfig) {
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

	lastStateHealth := promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "salt_state_health",
			Help: "Last state success state, 0=Failed, 1=Success",
		},
		[]string{"minion", "function", "state"},
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
				newJobCounter.WithLabelValues(event.Data.Fun, state, "true").Inc()
				expectedResponsesNumber.WithLabelValues(event.Data.Fun, state).Add(float64(event.TargetNumber))
			case "ret":
				state := event.ExtractState()
				if event.IsScheduleJob {
					// for scheduled job, when the states in the job actually failed
					// - the global "success" value is always true
					// - the substate success is false, and the global retcode is > 0
					// using retcode could be enough, but in case there are other corner cases, we combine both values
					success := event.Data.Success && (event.Data.Retcode == 0)
					scheduledJobReturnCounter.WithLabelValues(
						event.Data.Fun,
						state,
						strconv.FormatBool(success),
					).Inc()
				} else {
					success := strconv.FormatBool(event.Data.Success)

					responsesCounter.WithLabelValues(
						event.Data.Id,
						success,
					).Inc()
					functionResponsesCounter.WithLabelValues(
						event.Data.Fun,
						state,
						success,
					).Inc()

					if metricsConfig.HealthMinions {
						if contains(metricsConfig.HealthFunctionsFilters, event.Data.Fun) &&
							contains(metricsConfig.HealthStatesFilters, state) {
							lastStateHealth.WithLabelValues(
								event.Data.Id,
								event.Data.Fun,
								state).Set(boolToFloat64(event.Data.Success))
						}
					}
				}
			}

			elapsed := time.Since(start)
			log.Debug().Str("metric conversion took", elapsed.String()).Send()
		}
	}
}
