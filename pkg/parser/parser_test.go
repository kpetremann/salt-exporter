package parser_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/kpetremann/salt-exporter/pkg/event"
	"github.com/kpetremann/salt-exporter/pkg/parser"
)

func TestParseEvent(t *testing.T) {
	tests := []struct {
		name string
		args map[string]interface{}
		want event.SaltEvent
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

	p := parser.NewEventParser(false)
	for _, test := range tests {
		parsed, err := p.Parse(test.args)
		if err != nil {
			t.Errorf("Unexpected error %s", err.Error())
		}

		if diff := cmp.Diff(parsed, test.want); diff != "" {
			t.Errorf("Mismatch for '%s' test:\n%s", test.name, diff)
		}
	}
}
