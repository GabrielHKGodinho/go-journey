package validanagram

// isAnagram reports whether t uses exactly the same letters as s.
// It assumes lowercase English letters, as the problem guarantees.
func isAnagram(s string, t string) bool {
	if len(s) != len(t) {
		return false
	}

	var counts [26]int
	for i := 0; i < len(s); i++ {
		counts[s[i]-'a']++
		counts[t[i]-'a']--
	}

	for _, count := range counts {
		if count != 0 {
			return false
		}
	}
	return true
}
