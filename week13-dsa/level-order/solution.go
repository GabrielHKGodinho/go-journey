package levelorder

// TreeNode is one node of a binary tree, as defined by LeetCode.
type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

// levelOrder returns the node values grouped by depth, top to bottom, each
// level from left to right (LeetCode 102).
//
// Invariant: at the start of every round of the outer loop, the queue holds
// exactly the nodes of one level. Taking len(queue) before the inner loop
// (the snapshot) tells how many nodes belong to this level; the children
// pushed during the round go behind them and form the next level.
//
// Time: O(n), n = number of nodes. Every node is enqueued and dequeued
// once, both O(1).
// Space: O(w) extra, w = the largest number of nodes on one level, not
// counting the output. Best case a skewed tree (w = 1); worst case a full
// tree, whose last level has about n/2 nodes, so O(n).
func levelOrder(root *TreeNode) [][]int {
	if root == nil {
		return [][]int{}
	}

	levels := [][]int{}
	queue := []*TreeNode{root}

	for len(queue) > 0 {
		// Snapshot: at this point the queue holds exactly one level.
		levelSize := len(queue)
		levelValues := make([]int, 0, levelSize)

		for range levelSize {
			node := queue[0]
			queue = queue[1:]
			levelValues = append(levelValues, node.Val)

			// Children go behind the current level and form the next one.
			if node.Left != nil {
				queue = append(queue, node.Left)
			}
			if node.Right != nil {
				queue = append(queue, node.Right)
			}
		}
		levels = append(levels, levelValues)
	}
	return levels
}

// levelOrderDFS gives the same result with a recursive pre-order DFS that
// carries the depth of each node.
//
// The DFS reaches depth d only after passing through depth d-1, so depths
// are discovered in order: the first node to arrive at a new depth always
// finds depth == len(levels), and that is when the new level is created.
// Visiting left before right keeps each level in left-to-right order.
//
// Time: O(n), 2n+1 calls (one per node plus one per nil child).
// Space: O(h) extra for the call stack, h = height: O(log n) balanced,
// O(n) skewed. Better than BFS on wide trees, worse on deep ones.
func levelOrderDFS(root *TreeNode) [][]int {
	levels := [][]int{}

	var visit func(node *TreeNode, depth int)
	visit = func(node *TreeNode, depth int) {
		if node == nil {
			return
		}
		if depth == len(levels) { // first node seen at this depth
			levels = append(levels, []int{})
		}
		// Reassign: append returns a new slice header, and the copy inside
		// levels would not see it otherwise.
		levels[depth] = append(levels[depth], node.Val)
		visit(node.Left, depth+1)
		visit(node.Right, depth+1)
	}

	visit(root, 0)
	return levels
}
