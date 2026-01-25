package core

// Transform applies a character mapping to the input string.
// It replaces each rune in the input string according to the provided mapping.
// If a rune is not present in the mapping, it remains unchanged.
func Transform(input string, mapping map[rune]rune) string {
	out := make([]rune, 0, len(input))
	// Iterate over each rune in the input string
	for _, r := range input {
		// Map rune using provided mapping, or keep it unchanged
		if v, ok := mapping[r]; ok {
			out = append(out, v)
		} else {
			out = append(out, r)
		}
	}
	return string(out)
}
