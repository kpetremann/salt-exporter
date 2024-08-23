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
	KeepRewBody bool
}

func NewEventParser(keepRawBody bool) Event {
	return Event{KeepRewBody: keepRawBody}
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

func statemoduleResult(event event.SaltEvent) *bool {
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

// ParseEvent parses a salt event.
func (e Event) Parse(message map[string]interface{}) (event.SaltEvent, error) {
	var body string

	if raw, ok := message["body"].([]byte); ok {
		body = string(raw)
	} else {
		body = message["body"].(string)
	}
	lines := strings.SplitN(body, "\n\n", 2)

	tag := lines[0]
	if !(strings.HasPrefix(tag, "salt/")) {
		return event.SaltEvent{}, errors.New("tag not supported")
	}
	log.Debug().Str("tag", tag).Msg("new event")

	parts := strings.Split(tag, "/")

	if len(parts) < 3 {
		return event.SaltEvent{}, errors.New("tag not supported")
	}

	eventModule := event.GetEventModule(tag)

	if eventModule == event.UnknownModule {
		return event.SaltEvent{}, errors.New("tag not supported. Module unknown")
	}

	// Extract job type from the tag
	jobType := strings.Split(tag, "/")[3]

	// Parse message body
	byteResult := []byte(lines[1])
	ev := event.SaltEvent{Tag: tag, Type: jobType, Module: eventModule}

	if e.KeepRewBody {
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
	ev.StateModuleSuccess = statemoduleResult(ev)

	// A runner are executed on the master but they do not provide their ID in the event
	if strings.HasPrefix(tag, "salt/run") && ev.Data.ID == "" {
		ev.Data.ID = "master"
	}

	return ev, nil
}
