package events_test

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/kpetremann/salt-exporter/pkg/events"
)

func TestParseEvent(t *testing.T) {
	eventChan := make(chan events.SaltEvent)
	tests := []struct {
		name string
		args map[string]interface{}
		want events.SaltEvent
	}{
		{
			name: "new job",
			args: fakeEventAsMap(fakeNewJobEvent()),
			want: expectedNewJob,
		},
		{
			name: "return job",
			args: fakeEventAsMap(fakeRetJobEvent()),
			want: expectedReturnJob,
		},
		{
			name: "new schedule job",
			args: fakeEventAsMap(fakeNewScheduleJobEvent()),
			want: expectedNewScheduleJob,
		},
		{
			name: "return ack schedule job",
			args: fakeEventAsMap(fakeAckScheduleJobEvent()),
			want: expectedAckScheduleJob,
		},
		{
			name: "return schedule job",
			args: fakeEventAsMap(fakeScheduleJobReturnEvent()),
			want: expectedScheduleJobReturn,
		},
		{
			name: "new state.sls",
			args: fakeEventAsMap(fakeNewStateSlsJobEvent()),
			want: expectedNewStateSlsJob,
		},
		{
			name: "return state.sls",
			args: fakeEventAsMap(fakeStateSlsReturnEvent()),
			want: expectedStateSlsReturn,
		},
		{
			name: "new state.single",
			args: fakeEventAsMap(fakeNewStateSingleEvent()),
			want: expectedNewStateSingle,
		},
		{
			name: "return state.single",
			args: fakeEventAsMap(fakeStateSingleReturnEvent()),
			want: expectedStateSingleReturn,
		},
	}

	for _, test := range tests {
		var parsed events.SaltEvent
		go events.ParseEvent(test.args, eventChan, false)

		select {
		case parsed = <-eventChan:
		case <-time.After(1 * time.Millisecond):
		}

		if diff := cmp.Diff(parsed, test.want); diff != "" {
			t.Errorf("Mismatch for '%s' test:\n%s", test.name, diff)
		}
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
		event events.SaltEvent
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
