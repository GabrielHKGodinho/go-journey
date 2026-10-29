package threesum

import "slices"

// threeSum returns every unique triplet of values, taken from three
// different indices of nums, that sums to zero. Two triplets with the same
// values in a different order count as the same triplet.
//
// Idea: sort once, fix nums[i], then run Two Sum II (167) on the rest of the
// slice with target -nums[i]. Sorting is also what removes duplicates:
// equal values end up next to each other, so they can be skipped without a
// map.
//
// The triplets come out in ascending order (by first value, then second),
// because i only moves forward and, for a fixed i, left only moves forward.
//
// It sorts nums in place: a caller that needs the original order must pass
// slices.Clone(nums).
//
// Time: O(n²), n = len(nums). Sorting is O(n log n); then, for each of the n
// fixed values, left and right only move toward each other, so together
// they take at most n steps. n * O(n) dominates the sort.
// Space: O(1) extra, not counting the output (slices.Sort works in place
// and only uses O(log n) stack).
func threeSum(nums []int) [][]int {
	slices.Sort(nums)
	triplets := [][]int{}

	for i := 0; i < len(nums)-2; i++ {
		// The smallest value left is positive: nothing from here on can sum
		// to zero. This is only a shortcut; the result is the same without it.
		if nums[i] > 0 {
			break
		}
		// Same fixed value as the previous round: its triplets were already
		// found. i > 0 keeps nums[i-1] in range.
		if i > 0 && nums[i] == nums[i-1] {
			continue
		}

		// Reset for every fixed value: the pair search starts over.
		left, right := i+1, len(nums)-1
		for left < right {
			sum := nums[i] + nums[left] + nums[right]
			switch {
			case sum < 0:
				// nums[left] is too small even with the largest partner left.
				left++
			case sum > 0:
				// nums[right] is too big even with the smallest partner left.
				right--
			default:
				triplets = append(triplets, []int{nums[i], nums[left], nums[right]})
				// Found → move both → skip equal values on the left.
				// With nums[i] and nums[left] fixed, the third value is forced,
				// so an equal nums[left] would rebuild the same triplet. Equal
				// values on the right need no skip: once left changes value,
				// the old right is too big and leaves through sum > 0.
				left++
				right--
				for left < right && nums[left] == nums[left-1] {
					left++
				}
			}
		}
	}
	return triplets
}
