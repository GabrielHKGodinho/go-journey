package containsduplicate

import (
	"slices"
	"testing"
)

func TestContainsDuplicate(t *testing.T) {
	tests := []struct {
		name string
		nums []int
		want bool
	}{
		{"duplicate at the ends", []int{1, 2, 3, 1}, true},
		{"all distinct", []int{1, 2, 3, 4}, false},
		{"many duplicates", []int{1, 1, 1, 3, 3, 4, 3, 2, 4, 2}, true},
		{"single element", []int{7}, false},
		{"repeated zero value", []int{0, 0}, true},
		{"extreme values", []int{-1_000_000_000, 1_000_000_000}, false},
	}

	solutions := map[string]func([]int) bool{
		"hash": containsDuplicate,
		"sort": containsDuplicateSort,
	}

	for solutionName, solve := range solutions {
		for _, tt := range tests {
			t.Run(solutionName+"/"+tt.name, func(t *testing.T) {
				original := slices.Clone(tt.nums)
				if got := solve(tt.nums); got != tt.want {
					t.Errorf("got %v, want %v", got, tt.want)
				}
				if !slices.Equal(original, tt.nums) {
					t.Errorf("input was modified: %v", tt.nums)
				}
			})
		}
	}
}
