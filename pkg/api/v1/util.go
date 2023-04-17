package api

// sliceCompare takes two slices and returns a map of the slice value as the key
// and an int as the value. the int will be 0 if the key exists in both slices,
// greater than one if it exists in s2 but not s1 and less than one if it exists
// in s1 but not s2.
func sliceCompare(s1, s2 []string) map[string]int {
	m := make(map[string]int)
	for _, p := range s2 {
		m[p]++
	}

	for _, p := range s1 {
		m[p]--
	}

	return m
}
