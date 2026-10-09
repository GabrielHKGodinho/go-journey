package twosumsorted

// twoSum returns the 1-based positions of the two numbers in the sorted
// slice numbers that add up to target. The problem guarantees exactly one
// answer. Indices are 0-based internally and converted only on return.
func twoSum(numbers []int, target int) []int {
	left, right := 0, len(numbers)-1
	for left < right {
		sum := numbers[left] + numbers[right]
		switch {
		case sum == target:
			return []int{left + 1, right + 1}
		case sum > target:
			// numbers[right] is too big even with the smallest value left.
			right--
		default:
			// numbers[left] is too small even with the largest value left.
			left++
		}
	}
	return nil
}
