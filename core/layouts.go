package core

import (
	"encoding/json"
	"errors"
	"os"
)

type LayoutsFile struct {
	Layouts map[string]map[string]string `json:"layouts"`
}

func LoadLayouts(path string) (map[string]map[rune]rune, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var raw LayoutsFile
	if err := json.Unmarshal(b, &raw); err != nil {
		return nil, err
	}

	if len(raw.Layouts) == 0 {
		return nil, errors.New("no layouts found")
	}

	result := make(map[string]map[rune]rune)

	for lang, mapping := range raw.Layouts {
		runeMap := make(map[rune]rune)
		for k, v := range mapping {
			rk := []rune(k)
			rv := []rune(v)
			if len(rk) == 1 && len(rv) == 1 {
				runeMap[rk[0]] = rv[0]
			}
		}
		result[lang] = runeMap
	}

	return result, nil
}
