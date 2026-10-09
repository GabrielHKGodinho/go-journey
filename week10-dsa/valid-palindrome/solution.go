package validpalindrome

import "unicode"

// isPalindrome reports whether s reads the same in both directions after
// lowercasing it and ignoring every character that is not a letter or digit.
// It compares from both ends inward, so it needs no extra memory.
func isPalindrome(s string) bool {
	left, right := 0, len(s)-1
	for left < right {
		for left < right && !isAlphanumeric(s[left]) {
			left++
		}
		for left < right && !isAlphanumeric(s[right]) {
			right--
		}
		if unicode.ToLower(rune(s[left])) != unicode.ToLower(rune(s[right])) {
			return false
		}
		left++
		right--
	}
	return true
}

// isAlphanumeric reports whether the ASCII byte b is a letter or a digit.
func isAlphanumeric(b byte) bool {
	r := rune(b)
	return unicode.IsLetter(r) || unicode.IsDigit(r)
}
