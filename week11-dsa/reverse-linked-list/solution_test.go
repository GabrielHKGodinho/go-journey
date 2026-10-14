package reverselist

import (
	"slices"
	"testing"
)

// fromSlice builds a list in LeetCode's notation: [1,2,3] is 1 -> 2 -> 3.
func fromSlice(values []int) *ListNode {
	var head *ListNode
	for i := len(values) - 1; i >= 0; i-- {
		head = &ListNode{Val: values[i], Next: head}
	}
	return head
}

// toSlice walks the list and fails the test if it is longer than limit,
// so a cycle created by a broken reversal fails instead of hanging.
func toSlice(t *testing.T, head *ListNode, limit int) []int {
	t.Helper()
	var values []int
	for curr := head; curr != nil; curr = curr.Next {
		if len(values) == limit {
			t.Fatalf("list is longer than %d nodes: probably a cycle", limit)
		}
		values = append(values, curr.Val)
	}
	return values
}

func TestReverseList(t *testing.T) {
	tests := []struct {
		name  string
		input []int
		want  []int
	}{
		{"empty list", nil, nil},
		{"single node", []int{1}, []int{1}},
		{"two nodes", []int{1, 2}, []int{2, 1}},
		{"example", []int{1, 2, 3, 4, 5}, []int{5, 4, 3, 2, 1}},
		{"repeated values", []int{7, 7, 3}, []int{3, 7, 7}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := toSlice(t, reverseList(fromSlice(tt.input)), len(tt.input))
			if !slices.Equal(got, tt.want) {
				t.Errorf("reverseList(%v) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

// TestReverseListReusesNodes checks the "in place" claim: the reversed list
// is made of the same nodes, only with their Next pointers turned around.
func TestReverseListReusesNodes(t *testing.T) {
	head := fromSlice([]int{1, 2, 3})
	first, last := head, head.Next.Next

	reversed := reverseList(head)

	if reversed != last {
		t.Errorf("new head is not the old last node")
	}
	if first.Next != nil {
		t.Errorf("old head should now be the tail and point to nil")
	}
}
