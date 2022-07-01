package events

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestParseEvent(t *testing.T) {
	eventChan := make(chan SaltEvent)
	tests := []struct {
		name string
		args map[string]interface{}
		want SaltEvent
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
	}

	for _, test := range tests {
		var parsed SaltEvent
		go ParseEvent(test.args, eventChan)

		select {
		case parsed = <-eventChan:
		case <-time.After(1 * time.Millisecond):
		}

		if diff := cmp.Diff(parsed, test.want); diff != "" {
			t.Errorf("Mismatch for '%s' test:\n%s", test.name, diff)
		}
	}
}
