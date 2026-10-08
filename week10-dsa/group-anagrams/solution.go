package groupanagrams

import "slices"

// groupAnagrams groups the strings that are anagrams of each other.
// Anagrams share the same sorted form, which is used as the map key.
// Groups are returned in the order their first word appears in strs.
func groupAnagrams(strs []string) [][]string {
	groupIndexByKey := make(map[string]int, len(strs))
	groups := make([][]string, 0, len(strs))

	for _, word := range strs {
		key := sortedASCII(word)
		if i, found := groupIndexByKey[key]; found {
			groups[i] = append(groups[i], word)
			continue
		}
		groupIndexByKey[key] = len(groups)
		groups = append(groups, []string{word})
	}
	return groups
}

// sortedASCII returns s with its bytes in ascending order.
// It assumes s is ASCII: multi-byte characters would be split apart.
func sortedASCII(s string) string {
	b := []byte(s)
	slices.Sort(b)
	return string(b)
}

// groupAnagramsCount uses the letter counts as the key instead of sorting.
// Arrays are comparable in Go, so [26]int can be a map key; slices cannot.
func groupAnagramsCount(strs []string) [][]string {
	groupIndexByKey := make(map[[26]int]int, len(strs))
	groups := make([][]string, 0, len(strs))

	for _, word := range strs {
		var key [26]int
		for i := 0; i < len(word); i++ {
			key[word[i]-'a']++
		}
		if i, found := groupIndexByKey[key]; found {
			groups[i] = append(groups[i], word)
			continue
		}
		groupIndexByKey[key] = len(groups)
		groups = append(groups, []string{word})
	}
	return groups
}
