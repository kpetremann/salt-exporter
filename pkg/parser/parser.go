package parser

import (
	"errors"
	"strings"

	"github.com/kpetremann/salt-exporter/pkg/event"
	"github.com/rs/zerolog/log"
	"github.com/vmihailenco/msgpack/v5"
)

const testArg = "test"
const mockArg = "mock"

type Event struct {
	KeepRawBody bool
}

func NewEventParser(KeepRawBody bool) Event {
	return Event{KeepRawBody: KeepRawBody}
}

// isDryRun checks if an event is run with test=True
//
// Salt stores can store this info at two locations:
//
// in args:
//
//	"arg": [
//		"somestate",
//		{
//			"__kwarg__": true,
//			"test": true
//		}
//	]
//
// or in fun_args:
//
//	"fun_args": [
//		"somestate",
//		{
//			"test": true
//		}
//	]
func getBoolKwarg(event event.SaltEvent, field string) bool {
	for _, arg := range event.Data.Arg {
		if fields, ok := arg.(map[string]interface{}); ok {
			if val, ok := fields[field].(bool); ok {
				return val
			}
		}
	}

	for _, funArg := range event.Data.FunArgs {
		if fields, ok := funArg.(map[string]interface{}); ok {
			if val, ok := fields[field].(bool); ok {
				return val
			}
		}
	}

	return false
}

func substateResult(event event.SaltEvent) *bool {
	substates, ok := event.Data.Return.(map[string]interface{})
	if !ok {
		return nil
	}

	for _, ret := range substates {
		substate, ok := ret.(map[string]interface{})
		if !ok {
			return nil
		}

		result, ok := substate["result"]
		if !ok {
			return nil
		}

		r, ok := result.(bool)
		if !ok {
			return nil
		}

		if !r {
			return &r
		}
	}

	success := true
	return &success
}

// ParseEvent parses a salt event
func (e Event) Parse(message map[string]interface{}) (event.SaltEvent, error) {
	body := string(message["body"].([]byte))
	lines := strings.SplitN(body, "\n\n", 2)

	tag := lines[0]
	if !(strings.HasPrefix(tag, "salt/job") || strings.HasPrefix(tag, "salt/run")) {
		return event.SaltEvent{}, errors.New("tag not supported")
	}
	log.Debug().Str("tag", tag).Msg("new event")

	// Extract job type from the tag
	job_type := strings.Split(tag, "/")[3]

	// Parse message body
	byteResult := []byte(lines[1])
	ev := event.SaltEvent{Tag: tag, Type: job_type}

	if e.KeepRawBody {
		ev.RawBody = byteResult
	}

	if err := msgpack.Unmarshal(byteResult, &ev.Data); err != nil {
		log.Warn().Str("error", err.Error()).Str("tag", tag).Msg("decoding_failure")
		return event.SaltEvent{}, err
	}

	// Extract other info
	ev.TargetNumber = len(ev.Data.Minions)
	ev.IsScheduleJob = ev.Data.Schedule != ""
	ev.IsTest = getBoolKwarg(ev, testArg)
	ev.IsMock = getBoolKwarg(ev, mockArg)
	ev.StateModuleSuccess = substateResult(ev)

	// A runner are executed on the master but they do not provide their ID in the event
	if strings.HasPrefix(tag, "salt/run") && ev.Data.Id == "" {
		ev.Data.Id = "master"
	}

	return ev, nil
}
