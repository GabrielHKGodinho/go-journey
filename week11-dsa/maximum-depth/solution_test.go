package maxdepth

import (
	"math"
	"testing"
)

// null marks a missing child in LeetCode's level-order notation.
// Node values are in [-100, 100], so math.MinInt can never be a real value.
const null = math.MinInt

// fromLevelOrder builds a tree from LeetCode's array notation, for example
// [3,9,20,null,null,15,7]. It is itself a BFS: every real node, in queue
// order, takes the next two entries as its left and right children.
// Children of null are not listed, and trailing nulls may be omitted.
func fromLevelOrder(values []int) *TreeNode {
	if len(values) == 0 || values[0] == null {
		return nil
	}

	root := &TreeNode{Val: values[0]}
	queue := []*TreeNode{root}
	next := 1

	for len(queue) > 0 && next < len(values) {
		node := queue[0]
		queue = queue[1:]

		if next < len(values) && values[next] != null {
			node.Left = &TreeNode{Val: values[next]}
			queue = append(queue, node.Left)
		}
		next++

		if next < len(values) && values[next] != null {
			node.Right = &TreeNode{Val: values[next]}
			queue = append(queue, node.Right)
		}
		next++
	}
	return root
}

// skewed builds a tree where every node has only a right child: a linked
// list, the worst case for the recursion depth.
func skewed(n int) *TreeNode {
	var root *TreeNode
	for i := range n {
		root = &TreeNode{Val: i % 100, Right: root}
	}
	return root
}

func TestFromLevelOrder(t *testing.T) {
	//       1
	//      / \
	//     2   3
	//    /
	//   4
	//    \
	//     5
	root := fromLevelOrder([]int{1, 2, 3, 4, null, null, null, null, 5})

	if got := root.Left.Left.Right.Val; got != 5 {
		t.Fatalf("root.Left.Left.Right.Val = %d, want 5", got)
	}
	if root.Right.Left != nil || root.Right.Right != nil {
		t.Errorf("node 3 should be a leaf")
	}
	if root.Left.Right != nil {
		t.Errorf("node 2 should have no right child")
	}
}

func TestMaxDepth(t *testing.T) {
	tests := []struct {
		name string
		root *TreeNode
		want int
	}{
		{"empty tree", fromLevelOrder(nil), 0},
		{"single node counts as one", fromLevelOrder([]int{1}), 1},
		{"example 1", fromLevelOrder([]int{3, 9, 20, null, null, 15, 7}), 3},
		{"only a right child", fromLevelOrder([]int{1, null, 2}), 2},
		{"only a left child", fromLevelOrder([]int{1, 2}), 2},
		{"deepest path on the right", fromLevelOrder([]int{1, 2, 3, 4, null, null, 5, null, null, null, 6}), 4},
		{"deepest path on the left", fromLevelOrder([]int{1, 2, 3, 4, null, null, null, null, 5}), 4},
		{"full tree", fromLevelOrder([]int{1, 2, 3, 4, 5, 6, 7}), 3},
		// A BFS that loops on len(queue) instead of a snapshot still returns
		// the right depth for most trees (only 4 of the 132 shapes with 6
		// nodes expose it). This shape makes it return 4.
		{"uneven levels break a BFS without snapshot", fromLevelOrder([]int{1, 2, 3, null, 4, 5, 6}), 3},
		{"skewed, LeetCode maximum size", skewed(10_000), 10_000},
	}

	solutions := []struct {
		name string
		fn   func(*TreeNode) int
	}{
		{"DFS", maxDepth},
		{"BFS", maxDepthBFS},
	}

	for _, s := range solutions {
		for _, tt := range tests {
			t.Run(s.name+"/"+tt.name, func(t *testing.T) {
				if got := s.fn(tt.root); got != tt.want {
					t.Errorf("%s = %d, want %d", s.name, got, tt.want)
				}
			})
		}
	}
}
