package event_test

import (
	"testing"

	"github.com/kpetremann/salt-exporter/pkg/event"
)

func getNewStateEvent() event.SaltEvent {
	return event.SaltEvent{
		Tag:          "salt/job/20220630000f000000000/new",
		Type:         "new",
		TargetNumber: 1,
		Data: event.EventData{
			Timestamp: "2022-06-30T00:00:00.000000",
			Fun:       "state.sls",
			Arg:       []interface{}{"test"},
			Jid:       "20220630000000000000",
			Minions:   []string{"node1"},
			Missing:   []string{},
			Tgt:       "node1",
			TgtType:   "glob",
			User:      "salt_user",
		},
		IsScheduleJob: false,
	}
}

func TestExtractState(t *testing.T) {
	stateSls := getNewStateEvent()

	stateSlsFunArg := getNewStateEvent()
	stateSlsFunArg.Data.Arg = nil
	stateSlsFunArg.Data.FunArgs = []interface{}{"test", map[string]bool{"dry_run": true}}

	stateSlsFunArgMap := getNewStateEvent()
	stateSlsFunArgMap.Data.Arg = nil
	stateSlsFunArgMap.Data.FunArgs = []interface{}{map[string]interface{}{"mods": "test", "dry_run": true}}

	stateApplyArg := getNewStateEvent()
	stateApplyArg.Data.Fun = "state.apply"

	stateApplyHighstate := getNewStateEvent()
	stateApplyHighstate.Data.Fun = "state.apply"
	stateApplyHighstate.Data.Arg = nil

	stateHighstate := getNewStateEvent()
	stateHighstate.Data.Fun = "state.highstate"
	stateHighstate.Data.Arg = nil

	tests := []struct {
		name  string
		event event.SaltEvent
		want  string
	}{
		{
			name:  "state via state.sls",
			event: stateSls,
			want:  "test",
		},
		{
			name:  "state via state.sls args + kwargs",
			event: stateSlsFunArg,
			want:  "test",
		},
		{
			name:  "state via state.sls kwargs only",
			event: stateSlsFunArgMap,
			want:  "test",
		},
		{
			name:  "state via state.apply args only",
			event: stateApplyArg,
			want:  "test",
		},
		{
			name:  "state via state.apply",
			event: stateApplyArg,
			want:  "test",
		},
		{
			name:  "highstate via state.apply",
			event: stateApplyHighstate,
			want:  "highstate",
		},
		{
			name:  "state.highstate",
			event: stateHighstate,
			want:  "highstate",
		},
	}

	for _, test := range tests {
		if res := test.event.ExtractState(); res != test.want {
			t.Errorf("Mismatch for '%s', wants '%s' got '%s' ", test.name, test.want, res)
		}
	}

}
