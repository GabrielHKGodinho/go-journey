// Package combinations solves LeetCode 77 (Combinations): all ways to choose
// k numbers out of 1..n, each combination in increasing order.
package combinations

import "slices"

// Combine builds each combination one position at a time.
//
// Template (backtracking):
//   - each level decides: the next number of the current combination
//   - choices: any number greater than the last one chosen, up to n
//   - leaf: the combination has k numbers
//
// Pruning: a level only tries numbers that still leave enough numbers after
// them to complete the combination, so no branch is explored in vain.
//
// Time: O(k · C(n,k)), the C(n,k) leaves, each copied in O(k).
// Space: O(k) auxiliary (path and recursion depth), O(k · C(n,k)) for the output.
func Combine(n, k int) [][]int {
	combinations := make([][]int, 0, binomial(n, k))
	path := make([]int, 0, k)

	var backtrack func(start int)
	backtrack = func(start int) {
		if len(path) == k {
			combinations = append(combinations, slices.Clone(path))
			return
		}

		// With `need` numbers still missing, the largest number worth trying
		// is n-need+1: picking anything after it leaves fewer than need-1 numbers.
		need := k - len(path)
		for next := start; next <= n-need+1; next++ {
			path = append(path, next) // choose
			backtrack(next + 1)       // explore
			path = path[:len(path)-1] // undo
		}
	}

	backtrack(1)
	return combinations
}

// CombineTakeSkip is the include/exclude version (same tree shape as Subsets):
// each level is one number i, and the two choices are "take i" or "skip i".
// Its recursion is n levels deep instead of k, so auxiliary space is O(n).
func CombineTakeSkip(n, k int) [][]int {
	combinations := make([][]int, 0, binomial(n, k))
	path := make([]int, 0, k)

	var backtrack func(i int)
	backtrack = func(i int) {
		if len(path) == k {
			combinations = append(combinations, slices.Clone(path))
			return
		}
		// Prune: even taking every remaining number (i..n) can't reach k.
		// This check also covers i > n, since then no numbers remain.
		if len(path)+(n-i+1) < k {
			return
		}

		path = append(path, i) // take i
		backtrack(i + 1)
		path = path[:len(path)-1] // undo

		backtrack(i + 1) // skip i
	}

	backtrack(1)
	return combinations
}

// binomial returns C(n,k), used only to preallocate the result.
func binomial(n, k int) int {
	result := 1
	for i := 1; i <= k; i++ {
		result = result * (n - k + i) / i // stays an integer at every step
	}
	return result
}
