package threesum

import (
	"math/rand/v2"
	"slices"
	"testing"
)

func TestThreeSum(t *testing.T) {
	tests := []struct {
		name string
		nums []int
		want [][]int
	}{
		{"example 1", []int{-1, 0, 1, 2, -1, -4}, [][]int{{-1, -1, 2}, {-1, 0, 1}}},
		{"example 2: no triplet", []int{0, 1, 1}, [][]int{}},
		{"example 3: three zeros", []int{0, 0, 0}, [][]int{{0, 0, 0}}},
		{"all positive", []int{1, 2, 3}, [][]int{}},
		// Many zeros, one triplet: the fixed value must be skipped too.
		{"four zeros", []int{0, 0, 0, 0}, [][]int{{0, 0, 0}}},
		// Two -1s as the fixed value would add [-1,0,1] twice without the
		// skip on i.
		{"repeated fixed value", []int{-1, -1, 0, 1}, [][]int{{-1, 0, 1}}},
		// The same value twice inside one triplet is allowed (different
		// indices).
		{"same value twice in a triplet", []int{-2, 0, 1, 1, 2}, [][]int{{-2, 0, 2}, {-2, 1, 1}}},
		// After [-2,0,2], the next pair is again 0 and 2: without the skip on
		// left the triplet would be added twice.
		{"repeated left value after a match", []int{-2, 0, 0, 2, 2}, [][]int{{-2, 0, 2}}},
		{"minimum length, no answer", []int{-1, 0, 2}, [][]int{}},
		{"values at the constraint limits", []int{-100000, 50000, 50000}, [][]int{{-100000, 50000, 50000}}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			// threeSum sorts in place: pass a copy so the table stays intact.
			got := threeSum(slices.Clone(tt.nums))
			if !equalTriplets(got, tt.want) {
				t.Errorf("threeSum(%v) = %v, want %v", tt.nums, got, tt.want)
			}
		})
	}
}

// TestThreeSumMatchesBruteForce compares against the O(n³) definition on
// random inputs with many repeated values, which is where dedup bugs hide.
func TestThreeSumMatchesBruteForce(t *testing.T) {
	rng := rand.New(rand.NewPCG(28, 10)) // fixed seed: failures are reproducible

	for round := range 500 {
		n := 3 + rng.IntN(15)
		nums := make([]int, n)
		for i := range nums {
			nums[i] = rng.IntN(9) - 4 // values in [-4, 4]: lots of duplicates
		}

		want := bruteForce(nums)
		got := threeSum(slices.Clone(nums))
		if !equalTriplets(got, want) {
			t.Fatalf("round %d: threeSum(%v) = %v, want %v", round, nums, got, want)
		}
	}
}

// bruteForce tries every index triple, normalizes each match by sorting it,
// and removes duplicates with a map keyed by an array (arrays are
// comparable, slices are not). Output is in ascending order, like threeSum.
func bruteForce(nums []int) [][]int {
	seen := map[[3]int]bool{}
	for i := 0; i < len(nums); i++ {
		for j := i + 1; j < len(nums); j++ {
			for k := j + 1; k < len(nums); k++ {
				if nums[i]+nums[j]+nums[k] != 0 {
					continue
				}
				key := [3]int{nums[i], nums[j], nums[k]}
				slices.Sort(key[:])
				seen[key] = true
			}
		}
	}

	triplets := [][]int{}
	for key := range seen {
		triplets = append(triplets, []int{key[0], key[1], key[2]})
	}
	slices.SortFunc(triplets, slices.Compare[[]int])
	return triplets
}

// equalTriplets compares two lists of triplets in order.
func equalTriplets(a, b [][]int) bool {
	return slices.EqualFunc(a, b, func(x, y []int) bool {
		return slices.Equal(x, y)
	})
}
