package reverselist

// ListNode is one node of a singly linked list, as defined by LeetCode.
type ListNode struct {
	Val  int
	Next *ListNode
}

// reverseList reverses the list in place and returns the new head.
//
// Invariant after every iteration: prev is the head of the part that is
// already reversed, curr is the head of the part still to reverse, and no
// node is lost between them. An empty list skips the loop and returns nil.
//
// Time: O(n), n = number of nodes. Each node is visited once, and each
// iteration does four assignments.
// Space: O(1) extra. Three pointers, no new nodes, no recursion.
func reverseList(head *ListNode) *ListNode {
	var prev *ListNode
	curr := head
	for curr != nil {
		next := curr.Next // save the rest before cutting the link
		curr.Next = prev  // turn the current node's arrow around
		prev = curr
		curr = next
	}
	return prev
}
