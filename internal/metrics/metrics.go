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
	newJobsCounter := promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "salt_new_job_total",
			Help: "Total number of new job processed",
		},
		[]string{"function", "success"},
	)
	responsesCounter := promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "salt_responses_total",
			Help: "Total number of response job processed",
		},
		[]string{"function", "minion", "success"},
	)
	scheduledJobReturnCounter := promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "salt_scheduled_job_return_total",
			Help: "Total number of scheduled job response",
		},
		[]string{"function", "minion", "success"},
	)
	expectedResponsesNumber := promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "salt_expected_responses_total",
			Help: "Total number of expected minions responses",
		},
		[]string{"function"},
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
				newJobsCounter.WithLabelValues(event.Data.Fun, strconv.FormatBool(event.Data.Success)).Inc()
				expectedResponsesNumber.WithLabelValues(event.Data.Fun).Add(float64(event.TargetNumber))
			case "ret":
				if event.IsScheduleJob {
					scheduledJobReturnCounter.WithLabelValues(
						event.Data.Fun,
						event.Data.Id,
						strconv.FormatBool(event.Data.Success),
					).Inc()
				} else {
					responsesCounter.WithLabelValues(
						event.Data.Fun,
						event.Data.Id,
						strconv.FormatBool(event.Data.Success),
					).Inc()
				}
			}
			elapsed := time.Since(start)
			log.Debug().Str("metric conversion took", elapsed.String()).Send()
		}
	}
}
