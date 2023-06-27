package metrics

import (
	"context"

	"github.com/kpetremann/salt-exporter/pkg/event"
	"github.com/rs/zerolog/log"
)

func boolToFloat64(b bool) float64 {
	if b {
		return 1.0
	}
	return 0.0
}

func eventToMetrics(event event.SaltEvent, r Registry) {
	switch event.Type {
	case "new":
		state := event.ExtractState()
		r.IncreaseNewJobTotal(event.Data.Fun, state, true)
		r.IncreaseExpectedResponsesTotal(event.Data.Fun, state, float64(event.TargetNumber))

	case "ret":
		state := event.ExtractState()
		success := event.Data.Success

		if event.IsScheduleJob {
			// for scheduled job, when the states in the job actually failed
			// - the global "success" value is always true
			// - the substate success is false, and the global retcode is > 0
			// using retcode could be enough, but in case there are other corner cases, we combine both values
			success = event.Data.Success && (event.Data.Retcode == 0)
			r.IncreaseScheduledJobReturnTotal(event.Data.Fun, state, event.Data.Id, success)
		} else {
			r.IncreaseResponseTotal(event.Data.Id, success) // TODO: move outside the if?
			r.IncreaseFunctionResponsesTotal(event.Data.Fun, state, event.Data.Id, success)
		}

		r.SetFunctionStatus(event.Data.Id, event.Data.Fun, state, success)
	}
}

func ExposeMetrics(ctx context.Context, eventChan <-chan event.SaltEvent, config Config) {
	registry := NewRegistry(config)

	for {
		select {
		case <-ctx.Done():
			log.Info().Msg("stopping event listener")
			return
		case event := <-eventChan:
			if config.Global.Filters.IgnoreTest && event.IsTest {
				return
			}
			if config.Global.Filters.IgnoreMock && event.IsMock {
				return
			}

			eventToMetrics(event, registry)
		}
	}
}
