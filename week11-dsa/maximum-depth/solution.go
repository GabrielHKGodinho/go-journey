package maxdepth

// TreeNode is one node of a binary tree, as defined by LeetCode.
type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

// maxDepth returns the number of nodes on the longest path from root down
// to a leaf (LeetCode 104 counts nodes, not edges). It is a post-order DFS:
// a node's depth needs both children's depths first.
//
// Time: O(n), n = number of nodes. There are 2n+1 calls (one per node plus
// one per nil child), each doing constant work.
// Space: O(h), h = height in nodes. Only the root-to-current path is on the
// call stack (h+1 frames at the peak): O(log n) balanced, O(n) skewed.
func maxDepth(root *TreeNode) int {
	if root == nil {
		return 0
	}
	return 1 + max(maxDepth(root.Left), maxDepth(root.Right))
}

// maxDepthBFS returns the same answer by counting levels with a queue.
//
// Time: O(n), every node is enqueued and dequeued once, and both are O(1).
// Space: O(w), w = maximum width of a level. Mid-level the queue mixes the
// rest of one level with the start of the next, but never exceeds 2w:
// O(1) skewed, O(n) full (the last level holds about n/2 nodes).
func maxDepthBFS(root *TreeNode) int {
	if root == nil {
		return 0
	}

	levels := 0
	queue := []*TreeNode{root}

	for len(queue) > 0 {
		nodesInLevel := len(queue) // snapshot: the queue holds exactly one level now
		levels++

		for range nodesInLevel {
			node := queue[0]
			queue = queue[1:] // O(1): moves the slice header, copies nothing

			if node.Left != nil {
				queue = append(queue, node.Left)
			}
			if node.Right != nil {
				queue = append(queue, node.Right)
			}
		}
	}
	return levels
}
