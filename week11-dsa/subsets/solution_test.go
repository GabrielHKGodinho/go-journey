package subsets

import (
	"cmp"
	"slices"
	"testing"
)

// normalize sorts each subset and then the list of subsets, so results in
// any order can be compared (LeetCode accepts any order).
func normalize(sets [][]int) [][]int {
	out := make([][]int, len(sets))
	for i, s := range sets {
		out[i] = slices.Sorted(slices.Values(s))
	}
	slices.SortFunc(out, func(a, b []int) int {
		if c := cmp.Compare(len(a), len(b)); c != 0 {
			return c
		}
		return slices.Compare(a, b)
	})
	return out
}

func rangeOf(n int) []int {
	nums := make([]int, n)
	for i := range nums {
		nums[i] = i + 1
	}
	return nums
}

var solutions = []struct {
	name string
	fn   func([]int) [][]int
}{
	{"iterative", subsetsIterative},
	{"backtracking", subsetsBacktracking},
}

func TestSubsets(t *testing.T) {
	tests := []struct {
		name string
		nums []int
		want [][]int
	}{
		{"single number", []int{0}, [][]int{{}, {0}}},
		{"two numbers, negative", []int{-5, 5}, [][]int{{}, {-5}, {5}, {-5, 5}}},
		{"example 1", []int{1, 2, 3}, [][]int{{}, {1}, {2}, {3}, {1, 2}, {1, 3}, {2, 3}, {1, 2, 3}}},
		// My own case: an unsorted input broke my first version, where
		// sorting a shared array changed a subset already in the result.
		{"unsorted input", []int{3, 2, 4, 1}, [][]int{
			{}, {1}, {2}, {3}, {4},
			{1, 2}, {1, 3}, {1, 4}, {2, 3}, {2, 4}, {3, 4},
			{1, 2, 3}, {1, 2, 4}, {1, 3, 4}, {2, 3, 4},
			{1, 2, 3, 4},
		}},
	}

	for _, sol := range solutions {
		for _, tt := range tests {
			t.Run(sol.name+"/"+tt.name, func(t *testing.T) {
				got := normalize(sol.fn(tt.nums))
				want := normalize(tt.want)
				if !slices.EqualFunc(got, want, slices.Equal[[]int]) {
					t.Errorf("%s(%v)\n got  %v\n want %v", sol.name, tt.nums, got, want)
				}
			})
		}
	}
}

// TestSubsetsCountAndDistinct covers every size up to LeetCode's maximum.
// Shared arrays first corrupt a subset at n = 5 (a duplicate replaces a
// missing subset), so small examples alone cannot catch that bug.
func TestSubsetsCountAndDistinct(t *testing.T) {
	for _, sol := range solutions {
		for n := 1; n <= 10; n++ {
			got := normalize(sol.fn(rangeOf(n)))
			if len(got) != 1<<n {
				t.Errorf("%s: n=%d returned %d subsets, want %d", sol.name, n, len(got), 1<<n)
				continue
			}
			for i := 1; i < len(got); i++ {
				if slices.Equal(got[i-1], got[i]) {
					t.Errorf("%s: n=%d has duplicate subset %v", sol.name, n, got[i])
					break
				}
			}
		}
	}
}

// TestSubsetsOwnArrays checks the cause, not only the symptom: no two
// non-empty subsets may share a backing array, even when the values still
// happen to be right.
func TestSubsetsOwnArrays(t *testing.T) {
	for _, sol := range solutions {
		owner := map[*int][]int{}
		for _, s := range sol.fn(rangeOf(10)) {
			if len(s) == 0 {
				continue
			}
			if other, ok := owner[&s[0]]; ok {
				t.Errorf("%s: %v and %v share a backing array", sol.name, other, s)
				break
			}
			owner[&s[0]] = s
		}
	}
}
