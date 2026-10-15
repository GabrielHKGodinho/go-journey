package validparentheses

// isValid reports whether every bracket in s is closed by the same type, in
// the correct order. LeetCode constraint: s only contains "()[]{}", so any
// byte that is not a closer is an opener.
//
// Every closer must match the most recent unmatched opener, and "most
// recent" is exactly what a stack returns (LIFO).
//
// Time: O(n), n = len(s). One pass; each byte is one map lookup plus a push
// or a pop, all O(1).
// Space: O(n). Not because of the make below: in the worst case ("((((")
// every byte is an opener and stays on the stack.
func isValid(s string) bool {
	// An odd length can never be fully paired. Early exit, same Big-O.
	if len(s)%2 == 1 {
		return false
	}

	openerFor := map[byte]byte{')': '(', ']': '[', '}': '{'}
	stack := make([]byte, 0, len(s))

	for i := 0; i < len(s); i++ {
		c := s[i]

		opener, isCloser := openerFor[c]
		if !isCloser {
			stack = append(stack, c) // opener: wait for its match
			continue
		}

		// Closer: the most recent unmatched opener must be its pair.
		// The length check comes first, so an empty stack never reaches
		// stack[-1] (that would panic).
		if len(stack) == 0 || stack[len(stack)-1] != opener {
			return false
		}
		stack = stack[:len(stack)-1]
	}

	return len(stack) == 0
}

// isValidSwitch is my first accepted version (9:27, 1 submit), cleaned up:
// one case per closer instead of a map. It skips hashing, so it is a bit
// faster; the map version is shorter and a new bracket type is one entry.
func isValidSwitch(s string) bool {
	stack := make([]byte, 0, len(s))

	for i := 0; i < len(s); i++ {
		var opener byte
		switch s[i] {
		case ')':
			opener = '('
		case ']':
			opener = '['
		case '}':
			opener = '{'
		default:
			stack = append(stack, s[i])
			continue
		}

		if len(stack) == 0 || stack[len(stack)-1] != opener {
			return false
		}
		stack = stack[:len(stack)-1]
	}

	return len(stack) == 0
}
