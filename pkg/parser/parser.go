package parser

import (
	"errors"
	"strings"

	"github.com/kpetremann/salt-exporter/pkg/event"
	"github.com/rs/zerolog/log"
	"github.com/vmihailenco/msgpack/v5"
)

type Event struct {
	KeepRawBody bool
}

func NewEventParser(KeepRawBody bool) Event {
	return Event{KeepRawBody: KeepRawBody}
}

// ParseEvent parses a salt event
//
// KeepRawBody is used to keep the raw body of the event.
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

	ev.TargetNumber = len(ev.Data.Minions)
	ev.IsScheduleJob = ev.Data.Schedule != ""

	// A runner are executed on the master but they do not provide their ID in the event
	if strings.HasPrefix(tag, "salt/run") && ev.Data.Id == "" {
		ev.Data.Id = "master"
	}

	return ev, nil
}
