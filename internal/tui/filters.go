package tui

import (
	"strings"

	"github.com/kpetremann/salt-exporter/internal/tui/list"
)

const negativePrefix = "!"

// Filter via multiple fields and keywords.
func WordsFilter(term string, targets []string) []list.Rank {
	var ranks []int

	term = strings.ToLower(term)
	splitedTerm := strings.Split(term, " ")

	for i, target := range targets {
		match := true

		for _, word := range splitedTerm {
			// check if excluding a substring
			negated := strings.HasPrefix(word, negativePrefix)
			if negated {
				word = strings.TrimPrefix(word, negativePrefix)
			}

			// we ignore empty terms
			if word == "" {
				continue
			}

			// look for the substring
			matching := strings.Contains(strings.ToLower(target), word)

			// if excluding the substring, invert the result
			if negated {
				matching = !matching
			}

			// if one of the words is not matching, the whole target is not matching
			if !matching {
				match = false
				break
			}
		}

		if match {
			ranks = append(ranks, i)
		}
	}

	result := make([]list.Rank, len(ranks))
	for i, indexMatching := range ranks {
		result[i] = list.Rank{
			Index:          indexMatching,
			MatchedIndexes: []int{},
		}
	}
	return result
}
