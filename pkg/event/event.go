package event

import (
	"encoding/json"
	"errors"
	"strings"

	"github.com/vmihailenco/msgpack/v5"
	"gopkg.in/yaml.v3"
)

type EventModule int

type WatchOp uint32

const (
	UnknownModule EventModule = iota
	RunnerModule
	JobModule
	BeaconModule
)

const (
	Accepted WatchOp = iota
	Removed
)

type WatchEvent struct {
	MinionName string
	Op         WatchOp
}

type EventData struct {
	Arg       []interface{} `msgpack:"arg"`
	Cmd       string        `msgpack:"cmd"`
	Fun       string        `msgpack:"fun"`
	FunArgs   []interface{} `msgpack:"fun_args"`
	ID        string        `msgpack:"id"`
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
	Tag                string
	Type               string
	Module             EventModule
	TargetNumber       int
	Data               EventData
	IsScheduleJob      bool
	RawBody            []byte
	IsTest             bool
	IsMock             bool
	StateModuleSuccess *bool
}

// RawToJSON converts raw body to JSON
//
// If indent is true, the JSON will be indented.
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

// RawToYAML converts raw body to YAML.
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

func GetEventModule(tag string) EventModule {
	tagParts := strings.Split(tag, "/")
	if len(tagParts) < 2 {
		return UnknownModule
	}
	switch tagParts[1] {
	case "run":
		return RunnerModule
	case "job":
		return JobModule
	case "beacon":
		return BeaconModule
	default:
		return UnknownModule
	}
}

// extractStateFromArgs extracts embedded state info.
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

// Extract state info from event.
func (e *SaltEvent) ExtractState() string {
	switch e.Data.Fun {
	case "state.sls", "state.apply":
		switch {
		case len(e.Data.Arg) > 0:
			return extractStateFromArgs(e.Data.Arg[0], "mods")
		case len(e.Data.FunArgs) > 0:
			return extractStateFromArgs(e.Data.FunArgs[0], "mods")
		case e.Data.Fun == "state.apply":
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
