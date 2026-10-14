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

| Tree shape | DFS: O(h) | BFS: O(w) |
|---|---|---|
| skewed (a linked list) | O(n), worst | O(1), best |
| balanced / full | O(log n), best | O(n), worst |

Same algorithm, different container: a stack (LIFO) gives DFS, a queue
(FIFO) gives BFS.
