package filters

import (
	"strings"
)

func matchTerm(s string, pattern string) bool {
	switch {
	case s == pattern:
		return true

	case pattern == "*":
		return true

	case strings.HasPrefix(pattern, "*") && strings.HasSuffix(pattern, "*"):
		return strings.Contains(s, pattern[1:len(pattern)-1])

	case strings.HasPrefix(pattern, "*"):
		return strings.HasSuffix(s, pattern[1:])

	case strings.HasSuffix(pattern, "*"):
		return strings.HasPrefix(s, pattern[:len(pattern)-1])

	default:
		return false
	}
}

// Function to check if a string exists in a slice of strings
//
// It supports the following wildcards as a prefix or suffix.
// Wildcard is not supported in the middle of the string.
//
// Examples:
//   - Match("foo", []string{"foo"}) -> true
//   - Match("foo", []string{"bar"}) -> false
//   - Match("foo", []string{"*"}) -> true
//   - Match("foo", []string{"*o"}) -> true
//   - Match("foo", []string{"f*"}) -> true
//   - Match("foo", []string{"*f"}) -> false
//   - Match("foo", []string{"*o*"}) -> true
func Match(value string, filters []string) bool {
	for _, pattern := range filters {

		if matchTerm(value, pattern) {
			return true
		}
	}
	return false
}
