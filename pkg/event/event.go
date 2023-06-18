package event

import (
	"encoding/json"
	"errors"

	"github.com/vmihailenco/msgpack/v5"
	"gopkg.in/yaml.v3"
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
	RawBody       []byte
}

// RawToJSON converts raw body to JSON
//
// If indent is true, the JSON will be indented
func (e SaltEvent) RawToJSON(indent bool) ([]byte, error) {
	if e.RawBody == nil {
		return nil, errors.New("raw body not registered")
	}

	var data interface{}
	if err := msgpack.Unmarshal(e.RawBody, &data); err != nil {
		return nil, err
	}
	if indent {
		return json.MarshalIndent(data, "", "  ")
	} else {
		return json.Marshal(data)
	}
}

// RawToYAML converts raw body to YAML
func (e SaltEvent) RawToYAML() ([]byte, error) {
	if e.RawBody == nil {
		return nil, errors.New("raw body not registered")
	}

	var data interface{}
	if err := msgpack.Unmarshal(e.RawBody, &data); err != nil {
		return nil, err
	}

	return yaml.Marshal(data)
}

// extractStateFromArgs extracts embedded state info
func extractStateFromArgs(args interface{}, key string) string {
	// args only
	if v, ok := args.(string); ok {
		return v
	}

	// kwargs
	if v, ok := args.(map[string]interface{}); ok {
		if _, keyExists := v[key]; !keyExists {
			return ""
		}
		if ret, isString := v[key].(string); isString {
			return ret
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
