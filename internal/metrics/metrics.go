package metrics

import (
	"context"
	"strconv"

	"github.com/kpetremann/salt-exporter/internal/filters"
	"github.com/kpetremann/salt-exporter/pkg/event"
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

func ExposeMetrics(ctx context.Context, eventChan <-chan event.SaltEvent, metricsConfig MetricsConfig) {
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
	lastFunctionStatus := promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "salt_function_status",
			Help: "Last function/state success, 0=Failed, 1=Success",
		},
		[]string{"minion", "function", "state"},
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
			if metricsConfig.IgnoreTest && event.IsTest || metricsConfig.IgnoreMock && event.IsMock {
				return
			}

			switch event.Type {
			case "new":
				state := event.ExtractState()
				newJobCounter.WithLabelValues(event.Data.Fun, state, "true").Inc()
				expectedResponsesNumber.WithLabelValues(event.Data.Fun, state).Add(float64(event.TargetNumber))
			case "ret":
				var success bool
				state := event.ExtractState()
				if event.IsScheduleJob {
					// for scheduled job, when the states in the job actually failed
					// - the global "success" value is always true
					// - the substate success is false, and the global retcode is > 0
					// using retcode could be enough, but in case there are other corner cases, we combine both values
					success = event.Data.Success && (event.Data.Retcode == 0)
					scheduledJobReturnCounter.WithLabelValues(
						event.Data.Fun,
						state,
						strconv.FormatBool(success),
					).Inc()
				} else {
					success = event.Data.Success

					if metricsConfig.HealthMinions {
						responsesCounter.WithLabelValues(
							event.Data.Id,
							strconv.FormatBool(success),
						).Inc()
					}

					functionResponsesCounter.WithLabelValues(
						event.Data.Fun,
						state,
						strconv.FormatBool(success),
					).Inc()
				}

				// Expose state/func status if feature enabled and matching filters
				if !metricsConfig.HealthMinions {
					continue
				}
				if !filters.Match(event.Data.Fun, metricsConfig.HealthFunctionsFilters) {
					continue
				}
				log.Debug().Msg("function matches")
				if !filters.Match(state, metricsConfig.HealthStatesFilters) {
					continue
				}
				lastFunctionStatus.WithLabelValues(event.Data.Id, event.Data.Fun, state).Set(boolToFloat64(success))
			}
		}
	}
}
