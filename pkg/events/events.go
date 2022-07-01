package events

import (
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/vmihailenco/msgpack/v5"
)

type EventData struct {
	Arg       []interface{} `msgpack:"arg"`
	Cmd       string        `msgpack:"cmd"`
	Fun       string        `msgpack:"fun"`
	FunArgs   []interface{} `msgpack:"fun_args"`
	Id        string        `msgpack:"id"`
	Jid       string        `msgpack:"jid"`
	JidStamp  string        `msgpack:"jid_stamp"`
	Minions   []string      `msgpack:"minions"`
	Missing   []string      `msgpack:"missing"`
	Out       string        `msgpack:"out"`
	Retcode   int           `msgpack:"retcode"`
	Return    interface{}   `msgpack:"return"`
	Tgt       interface{}   `msgpack:"tgt"`
	TgtType   string        `msgpack:"tgt_type"`
	Timestamp string        `msgpack:"_stamp"`
	User      string        `msgpack:"user"`
	Schedule  string        `msgpack:"schedule"`
	Success   bool          `msgpack:"success"`
}

type SaltEvent struct {
	Tag           string
	Type          string
	TargetNumber  int
	Data          EventData
	IsScheduleJob bool
}

func ParseEvent(message map[string]interface{}, eventChan chan SaltEvent) {
	body := string(message["body"].([]byte))
	lines := strings.SplitN(body, "\n\n", 2)

	tag := lines[0]
	if !strings.HasPrefix(tag, "salt/job") {
		return
	}
	log.Debug().Str("tag", tag).Msg("new event")

	// Extract job type from the tag
	job_type := strings.Split(tag, "/")[3]

	// Parse message body
	event := SaltEvent{Tag: tag, Type: job_type}
	byteResult := []byte(lines[1])

	if err := msgpack.Unmarshal(byteResult, &event.Data); err != nil {
		var debug map[string]interface{}
		msgpack.Unmarshal(byteResult, &debug)
		log.Trace().Interface("raw event", debug).Msg("failed to parse event")

		log.Debug().Str("error", err.Error()).Interface("debug", debug).Msg("decoding_failure")
		log.Warn().Str("error", err.Error()).Str("tag", tag).Msg("decoding_failure")

		return
	}

	event.TargetNumber = len(event.Data.Minions)
	event.IsScheduleJob = event.Data.Schedule != ""

	eventChan <- event
}
