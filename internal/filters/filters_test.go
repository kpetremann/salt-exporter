package filters_test

import (
	"testing"

	"github.com/kpetremann/salt-exporter/internal/filters"
)

func TestMatchesFilters(t *testing.T) {
	var tests = []struct {
		value   string
		filters []string
		want    bool
	}{
		{"foo", []string{"foo"}, true},
		{"foo", []string{"bar"}, false},
		{"foo", []string{"*"}, true},
		{"foo", []string{"*o"}, true},
		{"foo", []string{"f*"}, true},
		{"foo", []string{"*o*"}, true},
		{"foo", []string{"foo", "bar"}, true},
		{"foo", []string{"bar", "baz"}, false},
		{"foo", []string{"*o", "bar"}, true},
		{"foo", []string{"bar", "*o"}, true},
		{"test.ping", []string{"test"}, false},
		{"test.ping", []string{"test.*"}, true},
		{"test.ping", []string{"test.ping*"}, true},
		{"state.sls", []string{"state.*"}, true},
		{"state.sls", []string{"state.*", "test.ping"}, true},
		{"state.sls", []string{"state", "test.ping"}, false},
	}

	for _, tt := range tests {
		got := filters.Match(tt.value, tt.filters)
		if got != tt.want {
			t.Errorf("'MatchFilters(%q, %q)' wants '%v' got '%v'", tt.value, tt.filters, tt.want, got)
		}
	}
}
