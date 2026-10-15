package validparentheses

import "testing"

func TestIsValid(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want bool
	}{
		{"one pair", "()", true},
		{"three types in sequence", "()[]{}", true},
		{"mixed nesting", "{[()()]}", true},
		// From my own test cases: looks nested, but '[' is closed by ')'.
		{"wrong type deep inside", "(((([())])))", false},
		{"empty string has nothing unmatched", "", true},
		{"wrong closing type", "(]", false},
		// Counting each type would accept this: the pairs cross.
		{"pairs cross", "([)]", false},
		{"single closer", ")", false},
		// Odd lengths never reach the loop in isValid (early exit), so the
		// empty-stack guard needs even-length inputs. Without the guard,
		// these read stack[-1] and panic.
		{"closer on an empty stack", ")(", false},
		{"closer after the stack empties", "())(", false},
		{"opener left over", "(", false},
		{"even length, all openers", "((", false},
	}

	solutions := []struct {
		name string
		fn   func(string) bool
	}{
		{"map", isValid},
		{"switch", isValidSwitch},
	}

	for _, sol := range solutions {
		for _, tt := range tests {
			t.Run(sol.name+"/"+tt.name, func(t *testing.T) {
				if got := sol.fn(tt.s); got != tt.want {
					t.Errorf("%s(%q) = %v, want %v", sol.name, tt.s, got, tt.want)
				}
			})
		}
	}
}
