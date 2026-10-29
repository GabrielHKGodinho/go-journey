# Week 13 — Pattern library

Same rules as `week10-dsa/PATTERNS.md`: read the **Constraints** first,
match the statement against each **Signal**, write 2-3 test cases of my
own (with expected outputs) before submitting, and state Big-O with every
variable defined.

---

## Two pointers inside a loop (k-sum) — "fix one, two-pointer the rest"

- **Signal:** "all **unique** triplets (or quadruplets) whose sum is a
  target", values may repeat, and `n` up to a few thousand. With
  n = 3·10³, O(n³) ≈ 2.7·10¹⁰ is too slow and O(n²) ≈ 9·10⁶ fits.
- **Core idea:** sort once, fix `nums[i]`, then run Two Sum II (167) on
  `nums[i+1:]` with target `-nums[i]`. 3Sum is 167 inside a loop.
- **Why moving a pointer is safe (the invariant):** the array is sorted,
  so `nums[right]` is the largest partner left for `nums[left]`. If even
  that pair is too small, every pair with `nums[left]` is too small:
  discard it (`left++`). Mirror argument for `right--` when the sum is
  too big. Without sorting, the argument falls apart.
- **Dedup without a map:** sorting puts equal values side by side.
  - Fixed value: skip when `i > 0 && nums[i] == nums[i-1]`.
  - After a match: **found → move both → skip equal values on the left.**
    With two values fixed the third is forced, so an equal `nums[left]`
    would rebuild the same triplet. No skip needed on the right.
- **Go tools:**
  - `slices.Sort` sorts **in place** (pdqsort, O(log n) stack). It changes
    the caller's slice: say so, and `slices.Clone` first if the original
    order matters.
  - If you ever dedup with a map, the key must be `[3]int`: arrays are
    comparable, slices are not (`map[[]int]bool` does not compile).
- **Complexity:** O(n²) time (sort O(n log n) + n fixed values × O(n)
  pointer steps). O(1) extra space, not counting the output.
- **Traps (all from my first attempt, 28/10):**
  - Two pointers on an **unsorted** slice.
  - Direction inverted: sum too small → `left++`, not `right--`.
  - `left, right` declared **outside** the loop over `i`: they must reset
    for every fixed value.
  - `right := len(nums)` is out of range; the last index is `len - 1`.
  - Not moving any pointer after a match → infinite loop.
  - `slice.Clone()` does not exist (slices have no methods); and cloning a
    slice built fresh with `[]int{...}` is pointless — clone only when an
    underlying array is shared (like `path` in 78).
  - Splitting into negatives and positives: duplicated code, and
    `[0,0,0]` (three non-negatives) was never searched.
- **Problems:** 15 3Sum. Next: 16 3Sum Closest, 18 4Sum.

## Tree BFS, grouping by level — update to week 11

The week 11 entry has the queue mechanics. What 102 added:

- **Answer to "where does a level end?":** *"At the start of each round,
  the queue holds exactly the nodes of one level. I capture the queue
  length before the inner loop and process only that many nodes; the
  children I enqueue go behind them and form the next level."*
- **Group by depth, not by parent.** In `[1,2,3,null,4,5,6]`, node 4
  (child of 2) and nodes 5, 6 (children of 3) are on the same level:
  `[[1],[2,3],[4,5,6]]`.
- **DFS alternative:** pass `depth` down; when `depth == len(levels)` this
  is the first node at that depth, so append a new empty level. Depths are
  discovered in order because reaching depth d means passing d-1. Visit
  left before right. O(h) space instead of O(w). Reassign
  `levels[depth] = append(levels[depth], v)` (the 09/10 slice-header trap).
- **Sentinel alternative** (a `nil` marker at the end of each level):
  works, but needs an extra rule (only re-push the marker if the queue is
  not empty, or the loop never ends). The snapshot is simpler.
- **Space:** O(w), w = widest level. Worst case is the **full** tree
  (last level ≈ n/2 → O(n)); best case is the skewed tree (O(1)).
- **Testing:** trees with one node per level (skewed, zigzag) do not test
  the grouping at all. Include a tree with several nodes per level coming
  from different parents, and the empty tree (constraint `[0, 2000]`).
- **Make/append detail:** the level size is known before the inner loop,
  so `make([]int, 0, levelSize)` avoids reallocations.
- **Problems:** 102 Binary Tree Level Order Traversal (BFS and DFS). Next:
  103 Zigzag, 199 Right Side View; 994 Rotting Oranges uses the same
  snapshot on a grid.

---

## Attempt log

| Date | Problem | Kind | Result | Notes |
|---|---|---|---|---|
| 28/10 | 15 3Sum | Spaced review (two pointers, last seen 07/10) | Not solved in 25 min; after the walkthrough, retyped from memory in 8:13, 1 submit | Recalled the pattern (fix one + two pointers) but not why it needs sorting. Line that needed thought on the retype: the skip on `left` after a match. |
| 28/10 | 102 Level Order | Weakest topic (BFS per level, missed on 12/10) | 10:00, 1 submit, no hints | Snapshot and non-nil filter both present. Big-O right with O(w) and the full-tree worst case. Reading the statement, grouped `[1,2,3,null,4,5,6]` by parent. |
