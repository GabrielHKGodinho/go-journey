package levelorder

import (
	"math"
	"math/rand/v2"
	"slices"
	"testing"
)

// null marks a missing child in LeetCode's level-order notation.
// Node values are in [-1000, 1000], so math.MinInt can never be a real value.
const null = math.MinInt

// fromLevelOrder builds a tree from LeetCode's array notation, for example
// [3,9,20,null,null,15,7]. Same helper as week11-dsa/maximum-depth: every
// real node, in queue order, takes the next two entries as its children.
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

var implementations = []struct {
	name string
	fn   func(*TreeNode) [][]int
}{
	{"BFS", levelOrder},
	{"DFS", levelOrderDFS},
}

func TestLevelOrder(t *testing.T) {
	tests := []struct {
		name   string
		values []int
		want   [][]int
	}{
		{"empty tree", []int{}, [][]int{}},
		{"single node", []int{1}, [][]int{{1}}},
		{"example 1", []int{3, 9, 20, null, null, 15, 7}, [][]int{{3}, {9, 20}, {15, 7}}},
		// 4 is a child of 2, while 5 and 6 are children of 3, but all three
		// are at depth 2. Catches grouping by parent instead of by depth, and
		// a BFS without the per-level snapshot.
		{"one level, different parents", []int{1, 2, 3, null, 4, 5, 6}, [][]int{{1}, {2, 3}, {4, 5, 6}}},
		{"skewed right", []int{1, null, 2, null, 3, null, 4, null, 5}, [][]int{{1}, {2}, {3}, {4}, {5}}},
		{"zigzag", []int{1, null, 2, 3, null, null, 4}, [][]int{{1}, {2}, {3}, {4}}},
		{"full tree", []int{1, 2, 3, 4, 5, 6, 7}, [][]int{{1}, {2, 3}, {4, 5, 6, 7}}},
		{"negative and repeated values", []int{-1000, 1000, 1000, -5}, [][]int{{-1000}, {1000, 1000}, {-5}}},
	}

	for _, impl := range implementations {
		for _, tt := range tests {
			t.Run(impl.name+"/"+tt.name, func(t *testing.T) {
				t.Parallel()
				got := impl.fn(fromLevelOrder(tt.values))
				if !equalLevels(got, tt.want) {
					t.Errorf("%s(%v) = %v, want %v", impl.name, tt.values, got, tt.want)
				}
			})
		}
	}
}

// TestBFSMatchesDFS builds random trees and checks that the two
// implementations, which find levels in completely different ways (queue
// snapshot vs. depth parameter), always agree.
func TestBFSMatchesDFS(t *testing.T) {
	rng := rand.New(rand.NewPCG(28, 10)) // fixed seed: failures are reproducible

	for round := range 300 {
		values := make([]int, 1+rng.IntN(40))
		for i := range values {
			if i > 0 && rng.IntN(4) == 0 {
				values[i] = null
			} else {
				values[i] = rng.IntN(2001) - 1000
			}
		}

		root := fromLevelOrder(values)
		bfs, dfs := levelOrder(root), levelOrderDFS(root)
		if !equalLevels(bfs, dfs) {
			t.Fatalf("round %d, tree %v: BFS = %v, DFS = %v", round, values, bfs, dfs)
		}
	}
}

// equalLevels compares two level lists in order.
func equalLevels(a, b [][]int) bool {
	return slices.EqualFunc(a, b, func(x, y []int) bool {
		return slices.Equal(x, y)
	})
}
