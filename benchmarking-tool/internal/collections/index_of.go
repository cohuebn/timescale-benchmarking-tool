package collections

// Get the index of a string in a slice of strings
func IndexOf(vs []string, t string) int {
	for i, v := range vs {
			if v == t {
					return i
			}
	}
	return -1
}