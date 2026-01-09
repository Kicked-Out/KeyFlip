package core

// UaToEn maps Ukrainian keyboard layout runes to English layout runes
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
