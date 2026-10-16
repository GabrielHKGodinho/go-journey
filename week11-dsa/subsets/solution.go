package subsets

import "slices"

// Both solutions walk the same decision tree: one level per number, two
// branches per level (with it / without it), and every leaf is a subset.
// LeetCode 78 allows any order, both of the subsets and inside each one.

// subsetsIterative is my own idea (17 min, 2 submits): start from the empty
// set, and for every number copy each existing subset and add the number
// to the copy, so the count doubles every round. It builds the tree one
// level at a time: after k numbers, result holds the 2^k nodes of level k.
//
// Time: O(n·2^n), n = len(nums). 2^n subsets, each built by copying up to n
// numbers.
// Space: O(n·2^n) total for the output. Extra beyond it: withNum, up to
// 2^(n-1) slice headers in the last round.
func subsetsIterative(nums []int) [][]int {
	result := [][]int{{}}

	for _, num := range nums {
		withNum := make([][]int, len(result))
		for i, subset := range result {
			// Clone the inner slice, the one that is about to grow. Cloning
			// result itself would be shallow: new headers pointing to the
			// same arrays, so this append could write into an array that a
			// subset already in result still uses (broke on [1,2,3,4,5]).
			withNum[i] = append(slices.Clone(subset), num)
		}
		result = slices.Concat(result, withNum)
	}

	return result
}

// subsetsBacktracking walks the tree depth first: choose, explore, undo.
// Including comes first, so the first subset stored is the whole input
// and the last one is the empty set.
//
// Time: O(n·2^n). 2^n leaves, each cloned in O(n); the 2^n-1 inner calls
// do O(1) work.
// Space: O(n) extra beyond the output: path plus at most n+1 frames.
func subsetsBacktracking(nums []int) [][]int {
	result := make([][]int, 0, 1<<len(nums)) // 1<<n = 2^n, the exact count
	path := make([]int, 0, len(nums))        // one path reused by every branch

	// i is the tree level: the index of the number this call decides.
	var backtrack func(i int)
	backtrack = func(i int) {
		if i == len(nums) { // leaf: every number has been decided
			result = append(result, slices.Clone(path)) // store a copy, never path
			return
		}

		// Branch 1: include nums[i].
		path = append(path, nums[i]) // choose
		backtrack(i + 1)             // explore
		path = path[:len(path)-1]    // undo

		// Branch 2: exclude nums[i]. Nothing to choose or undo.
		backtrack(i + 1)
	}

	backtrack(0)
	return result
}
