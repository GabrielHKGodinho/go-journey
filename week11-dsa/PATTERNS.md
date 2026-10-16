# Week 11 — Pattern library

Same rules as `week10-dsa/PATTERNS.md`: match the statement against each
**Signal**, write 2-3 test cases of my own before submitting, and state
Big-O with every variable defined. The re-solve log stays in the week 10
file.

---

## Tree DFS (recursive) — "the answer for a node comes from its subtrees"

- **Signal:** a binary tree where the answer for a node can be built from
  the answers for its left and right subtrees (height, size, sum, "is it
  balanced", "are two trees the same").
- **Core idea:** ask *"if I already knew the answer for both children, how
  would I compute this node's answer?"* Base case on `nil`, recurse on both
  children, combine. Trust the recursive call instead of simulating it.
- **Orders:** the path is always the same; only *when* the node is
  processed changes.
  - Pre-order (node, left, right): the node passes something **down**, or
    copying / serializing a tree.
  - In-order (left, node, right): a BST yields its values sorted.
  - Post-order (left, right, node): the node needs its **children's
    answers** first (height, size, folder size, freeing memory).
- **Go tools:**
  - `TreeNode` holds `*TreeNode`, never `TreeNode` (a struct cannot
    contain itself by value); a missing child is `nil`.
  - Builtin `max` (Go 1.21+).
  - Recursive closure: `var visit func(*TreeNode)` first, then
    `visit = func(...) { ... visit(...) }`. Capture the result slice in the
    closure instead of passing a slice (its header is copied by value).
- **Complexity:** O(n) time (2n+1 calls). O(h) space for the call stack,
  h+1 frames at the peak: O(log n) balanced, O(n) skewed.
- **Traps:**
  - Base case on "is a leaf" instead of `nil`: a node with one child then
    recurses into `nil` and panics.
  - Depth in nodes vs edges differs by 1. LeetCode 104 counts nodes; ask.
  - Very deep input: a goroutine stack grows to 1 GB by default (about 32
    bytes per `maxDepth` frame, measured), and overflowing it is a
    `fatal error` that `recover` cannot catch. Bound the depth or go
    iterative for untrusted input.
- **Problems:** 104 Maximum Depth of Binary Tree.

## Tree BFS (level order) — "process the tree level by level"

- **Signal:** the statement talks about **levels** (values per level, the
  right side view, level averages), or wants the node **closest to the
  root** (minimum depth: stop at the first leaf).
- **Core idea:** a FIFO queue. Take from the front, push non-nil children
  at the back. To know where a level ends, take a **snapshot** of the queue
  length at the start of each round and process exactly that many nodes:
  at that moment the queue holds exactly one full level.
- **Go tools:**
  - The queue is a slice: `queue = append(queue, x)` to enqueue,
    `node := queue[0]` then `queue = queue[1:]` to dequeue (O(1), no copy).
  - `queue[1:]` on a slice of length 1 is valid and empty: slicing allows
    `low == len`, indexing `queue[1]` does not.
  - `for range n` (Go 1.22+) for "repeat n times".
  - Not a channel (it blocks and needs another goroutine) and not
    `container/list` (one allocation per element, values are `any`).
- **Complexity:** O(n) time (each node enters and leaves once). O(w) space,
  w = maximum level width, never above 2w mid-level: O(1) skewed, O(n) full.
- **Traps:**
  - Looping on `len(queue)` instead of the snapshot: the limit moves while
    the loop runs. It does not panic, and for depth it is right on most
    shapes (only 4 of the 132 trees with 6 nodes expose it), so it needs a
    test built to catch it.
  - Enqueuing `nil` children and then reading `node.Val`.
- **Problems:** 104 Maximum Depth of Binary Tree (BFS version).

## Linked list, in-place reversal — "turn the arrows around"

- **Signal:** a singly linked list that must be reversed (whole or in
  part), with O(1) extra space or "in place" in the statement.
- **Core idea:** three pointers. Every iteration: save `next := curr.Next`,
  turn `curr.Next = prev`, then `prev = curr` and `curr = next`. Invariant:
  `prev` heads the reversed part, `curr` heads the rest, nothing is lost.
  When `curr` is nil, `prev` is the new head.
- **Go tools:**
  - `var prev *ListNode` (nil is the zero value; `= nil` is redundant).
  - `for curr := head; curr != nil; curr = curr.Next` to walk a list.
  - A pointer argument is a copy: changing `node.Next` is visible to the
    caller, reassigning `head` is not. Return the new head instead.
- **Complexity:** O(n) time (each node once, constant work). O(1) extra
  space iterative; O(n) recursive (one stack frame per node, since a list
  is a skewed tree).
- **Traps:**
  - Turning the arrow before saving `next`: the rest of the list is lost
    silently and the function returns only the first node.
  - Walking with `head` itself and losing the start of the list.
  - Corner cases: empty list, one node, two nodes.
- **Problems:** 206 Reverse Linked List.

## DFS vs BFS — choosing

| Tree shape | DFS recursive: O(h) | DFS iterative: O(h) | BFS: O(w) |
|---|---|---|---|
| skewed (a linked list) | O(n), worst | O(1), best | O(1), best |
| balanced / full | O(log n), best | O(log n) | O(n), worst |
| comb (a leaf on the left of every level) | O(n) | O(n), worst | O(1) |

Same algorithm, different container: a stack (LIFO) gives DFS, a queue
(FIFO) gives BFS.

## Stack — "match with the most recent unmatched thing"

- **Signal:** nested or paired structure (brackets, tags, `3[a2[c]]`), or
  "undo the last step": every new item must be checked against the **most
  recent one still open**.
- **Core idea:** push what is still open; when a closing item arrives,
  compare it with the top and pop. Whatever is left at the end was never
  closed.
- **Go tools:**
  - The stack is a slice with the top at the **end**:
    `stack = append(stack, x)` to push, `stack[len(stack)-1]` to peek,
    `stack = stack[:len(stack)-1]` to pop. Pop copies nothing, O(1); push
    is O(1) amortized. A top at the front would copy everything per push.
  - A popped element still sits in the backing array. With pointers in a
    long-lived stack, set `stack[len(stack)-1] = nil` before popping so the
    GC can free it (go.dev/wiki/SliceTricks).
  - `map[closer]opener` turns one case per type into one lookup.
- **Complexity:** O(n) time, one push or pop per item. O(n) space in the
  worst case, when everything stays open.
- **Traps:**
  - Peeking or popping an empty stack: `stack[-1]` panics in Go (PHP and
    JS just return null / undefined). Check `len(stack) == 0` first.
  - Forgetting the final `len(stack) == 0`: "(" is not valid.
  - An early exit (odd length) can hide a guard from the tests: test the
    empty-stack case with an even-length input like ")(".
- **Problems:** 20 Valid Parentheses. Next: 394 Decode String, 232
  Implement Queue using Stacks.

## Tree DFS (iterative) — "DFS without recursion"

- **Signal:** same as recursive DFS, plus "without recursion", very deep
  input, or the interviewer's follow-up "can you do it iteratively?".
- **Core idea:** replace the call stack with your own slice. A loop has
  no "way back" for return values, so the computation turns **top-down**:
  each stack item carries what the frame used to carry (the node and its
  depth), children get `depth+1` when pushed, and a variable keeps the
  best answer so far.
- **Go tools:**
  - A small struct for the item, declared inside the function when nothing
    else needs it: `type frame struct { node *TreeNode; depth int }`.
  - `[]frame{{node: root, depth: 1}}`: the inner type name can be omitted.
  - Builtin `max` (Go 1.21+).
- **Complexity:** O(n) time. O(h) space: the current path plus at most one
  waiting sibling per level. Measured: full tree (10 levels) peaks at 10,
  skewed at 1, comb at about n/2.
- **Traps:**
  - Pushing `nil` children and then reading `top.node.Left`: filter before
    pushing, or skip `nil` after popping, but pick one.
  - Order matters when the output is a list: for pre-order (node, left,
    right), push **right first**, so left is popped first (LIFO).
  - Explaining its space with the BFS argument ("the last level holds n/2
    nodes"): the DFS stack never holds a whole level.
- **Problems:** 104 Maximum Depth of Binary Tree (iterative version).

## Backtracking — "list every subset / combination / permutation"

- **Signal:** the output is **all** the ways to choose or arrange items
  (subsets, combinations, permutations, every valid string). The answer
  itself is exponential, so nothing faster than enumerating exists.
- **Core idea:** a DFS over a tree of decisions. Before coding, answer on
  paper, with a small input:
  1. What does one **level** of the tree decide?
  2. What are the **choices** at that level?
  3. When is a node a **leaf** (the answer is complete)?
  Then: choose, explore, undo. One `path` is reused by every branch; the
  undo runs on the way back up and restores the state the node arrived
  with.
- **Go tools:**
  - Recursive closure: `var backtrack func(i int)` first, then assign.
  - `path` is a stack: `append` to choose, `path[:len(path)-1]` to undo.
  - Store `slices.Clone(path)` at the leaf, never `path` itself.
  - `1<<n` is 2^n: pre-size the result when the count is known.
  - Debug a trace by indenting with depth:
    `fmt.Printf("%s%v\n", strings.Repeat("  ", i), path)`.
- **Complexity:** subsets: O(n·2^n) time (2^n leaves, an O(n) clone each),
  O(n) extra space beyond the output (path plus n+1 frames).
- **Traps:**
  - Forgetting the undo. If the leaf test is `len(path) == n`, the path
    grows past n and the recursion never stops (stack overflow, a fatal
    error). If the level is a parameter `i`, it stops but returns wrong
    subsets.
  - Counting calls without the base cases: `walk(n)` with two calls makes
    2^(n+1)-1 calls, and its peak stack is n+1 frames, not n.
  - Hard to trace by hand: after a call returns, you must remember **where
    in the parent** you resume (after the first call or the second) and
    what `path` held then. Draw the tree; do not keep it in your head.
- **Problems:** 78 Subsets (backtracking and iterative). Next: 77
  Combinations, 22 Generate Parentheses, 46 Permutations.

## Slice aliasing — "is this array only mine?"

Not a problem pattern, a Go trap that broke 78 (and 49 on 09/10).

- A slice is a header `{ptr, len, cap}`. Copying it (assignment, passing
  it, storing it in another slice) copies the header, not the array.
- `append(s, x)` writes **into s's array** when `len < cap`, and does not
  change `s` itself (its len stays the same). Two appends on the same base
  with spare capacity write the **same slot**: `b := append(a, 1)` then
  `c := append(a, 2)` leaves both as `[... 2]`.
- `slices.Clone` copies one level. On a `[][]int` it copies the inner
  headers, so the inner arrays are still shared (like `[...grid]` in JS).
  Clone the slice that is about to grow or be stored.
- The bug is intermittent: it only shows when some array has spare
  capacity and is appended to twice. Test sizes beyond the examples, or
  test the cause directly (no two results share `&s[0]`).
- Rule: `append(a, x)` without `a = ...` is a warning sign. Either you
  want a new slice derived from `a` (clone first), or you forgot to
  reassign.
