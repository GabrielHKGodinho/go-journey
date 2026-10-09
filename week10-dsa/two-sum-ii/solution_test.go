package twosumsorted

import (
	"slices"
	"testing"
)

func TestTwoSum(t *testing.T) {
	tests := []struct {
		name    string
		numbers []int
		target  int
		want    []int
	}{
		{"example", []int{2, 7, 11, 15}, 9, []int{1, 2}},
		{"answer at both ends", []int{2, 3, 4}, 6, []int{1, 3}},
		{"negative numbers", []int{-1, 0}, -1, []int{1, 2}},
		{"repeated values", []int{1, 3, 3, 8}, 6, []int{2, 3}},
		{"both pointers move", []int{1, 2, 4, 10}, 6, []int{2, 3}},
		{"negative target", []int{-5, -3, 0, 4}, -8, []int{1, 2}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := twoSum(tt.numbers, tt.target); !slices.Equal(got, tt.want) {
				t.Errorf("twoSum(%v, %d) = %v, want %v", tt.numbers, tt.target, got, tt.want)
			}
		})
	}
}
