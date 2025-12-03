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

func eventToMetrics(e event.SaltEvent, r *Registry) {
	if e.Module == event.BeaconModule {
		if e.Type != "status" {
			return
		}
		r.UpdateLastHeartbeat(e.Data.ID)
		return
	}

	switch e.Type {
	case "new":
		state := e.ExtractState()
		r.IncreaseNewJobTotal(e.Data.Fun, state)
		r.IncreaseExpectedResponsesTotal(e.Data.Fun, state, float64(e.TargetNumber))

	case "ret":
		state := e.ExtractState()
		success := e.Data.Success

		if e.IsScheduleJob {
			// for scheduled job, when the states in the job actually failed
			// - the global "success" value is always true
			// - the state module success is false, but the global retcode is > 0
			// - if defined, the "result" of a state module in event.Return covers
			//   the corner case when retccode is not properly computed by Salt.
			//
			// using retcode and state module success could be enough, but we combine all values
			// in case there are other corner cases.
			success = e.Data.Success && (e.Data.Retcode == 0)
			if e.StateModuleSuccess != nil {
				success = success && *e.StateModuleSuccess
			}
			r.IncreaseScheduledJobReturnTotal(e.Data.Fun, state, e.Data.ID, success)
		} else {
			r.IncreaseFunctionResponsesTotal(e.Data.Fun, state, e.Data.ID, success)
		}

		r.IncreaseResponseTotal(e.Data.ID, success)
		r.SetFunctionStatus(e.Data.ID, e.Data.Fun, state, success)
	}
}

func ExposeMetrics(ctx context.Context, eventChan <-chan event.SaltEvent, watchChan <-chan event.WatchEvent, config Config) {
	registry := NewRegistry(config)

	for {
		select {
		case <-ctx.Done():
			log.Info().Msg("stopping event listener")
			return
		case e := <-watchChan:
			if e.Op == event.Accepted {
				registry.AddObservableMinion(e.MinionName)
			}
			if e.Op == event.Removed {
				registry.DeleteObservableMinion(e.MinionName)
			}
		case e := <-eventChan:
			if config.Global.Filters.IgnoreTest && e.IsTest {
				continue
			}
			if config.Global.Filters.IgnoreMock && e.IsMock {
				continue
			}

			eventToMetrics(e, &registry)
		}
	}
}
