package core

func Transform(input string, mapping map[rune]rune) string {
	out := make([]rune, 0, len(input))
	for _, r := range input {
		if v, ok := mapping[r]; ok {
			out = append(out, v)
		} else {
			out = append(out, r)
		}
	}
	return string(out)
}
