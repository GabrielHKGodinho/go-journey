package containsduplicate

import "slices"

// containsDuplicate reports whether any value appears at least twice in nums.
// It remembers every value seen so far in a set and stops at the first repeat.
// struct{} is used to show the map works as a set; since Go 1.24 it takes the
// same memory as bool.
func containsDuplicate(nums []int) bool {
	seen := make(map[int]struct{}, len(nums))
	for _, n := range nums {
		if _, found := seen[n]; found {
			return true
		}
		seen[n] = struct{}{}
	}
	return false
}

// containsDuplicateSort sorts a copy of nums so equal values become neighbors.
// Without the copy, the caller's slice would be reordered.
func containsDuplicateSort(nums []int) bool {
	sorted := slices.Clone(nums)
	slices.Sort(sorted)
	for i := 1; i < len(sorted); i++ {
		if sorted[i] == sorted[i-1] {
			return true
		}
	}
	return false
}
