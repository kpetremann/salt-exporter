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

func extractStateFromArgs(args interface{}, key string) string {
	// args only
	if v, ok := args.(string); ok {
		return v
	}
	// kwargs
	if v, ok := args.(map[string]interface{}); ok {
		if _, ok := v[key]; ok {
			return v[key].(string)
		}
	}

	return ""
}

// Extract state info from event
func (e *SaltEvent) ExtractState() string {
	switch e.Data.Fun {
	case "state.sls", "state.apply":
		if len(e.Data.Arg) > 0 {
			return extractStateFromArgs(e.Data.Arg[0], "mods")
		} else if len(e.Data.FunArgs) > 0 {
			return extractStateFromArgs(e.Data.FunArgs[0], "mods")
		} else if e.Data.Fun == "state.apply" {
			return "highstate"
		}
	case "state.single":
		if len(e.Data.Arg) > 0 {
			return extractStateFromArgs(e.Data.Arg[0], "fun")
		} else if len(e.Data.FunArgs) > 0 {
			return extractStateFromArgs(e.Data.FunArgs[0], "fun")
		}
	case "state.highstate":
		return "highstate"
	}
	return ""
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
		log.Warn().Str("error", err.Error()).Str("tag", tag).Msg("decoding_failure")
		return
	}

	event.TargetNumber = len(event.Data.Minions)
	event.IsScheduleJob = event.Data.Schedule != ""

	eventChan <- event
}
