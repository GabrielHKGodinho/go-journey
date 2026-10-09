package validpalindrome

import "testing"

func TestIsPalindrome(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want bool
	}{
		{"sentence with punctuation", "A man, a plan, a canal: Panama", true},
		{"not a palindrome", "race a car", false},
		{"only a space", " ", true},
		{"only punctuation", ".,", true},
		{"digit and letter differ", "0P", false},
		{"digits count", "1a1", true},
		{"single character", "a", true},
		{"invalid at one end", "a.", true},
		{"case is ignored", "Aa", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isPalindrome(tt.s); got != tt.want {
				t.Errorf("isPalindrome(%q) = %v, want %v", tt.s, got, tt.want)
			}
		})
	}
}
