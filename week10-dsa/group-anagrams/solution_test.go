package groupanagrams

import (
	"slices"
	"testing"
)

func TestGroupAnagrams(t *testing.T) {
	tests := []struct {
		name string
		strs []string
		want [][]string
	}{
		{"example", []string{"eat", "tea", "tan", "ate", "nat", "bat"},
			[][]string{{"eat", "tea", "ate"}, {"tan", "nat"}, {"bat"}}},
		{"empty string", []string{""}, [][]string{{""}}},
		{"single letter", []string{"a"}, [][]string{{"a"}}},
		{"empty strings group together", []string{"", "b", ""}, [][]string{{"", ""}, {"b"}}},
		{"same letters, different counts", []string{"aab", "abb"}, [][]string{{"aab"}, {"abb"}}},
	}

	solutions := map[string]func([]string) [][]string{
		"sort":  groupAnagrams,
		"count": groupAnagramsCount,
	}

	for solutionName, solve := range solutions {
		for _, tt := range tests {
			t.Run(solutionName+"/"+tt.name, func(t *testing.T) {
				got := solve(tt.strs)
				if !slices.EqualFunc(got, tt.want, slices.Equal[[]string]) {
					t.Errorf("got %q, want %q", got, tt.want)
				}
			})
		}
	}
}
