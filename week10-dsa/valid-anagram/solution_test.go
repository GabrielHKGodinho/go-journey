package validanagram

import "testing"

func TestIsAnagram(t *testing.T) {
	tests := []struct {
		name string
		s, t string
		want bool
	}{
		{"anagram", "anagram", "nagaram", true},
		{"different letters", "rat", "car", false},
		{"different lengths", "a", "ab", false},
		{"same letters, different counts", "aab", "abb", false},
		{"letter only in t", "ab", "cc", false},
		{"single letter", "a", "a", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isAnagram(tt.s, tt.t); got != tt.want {
				t.Errorf("isAnagram(%q, %q) = %v, want %v", tt.s, tt.t, got, tt.want)
			}
		})
	}
}
