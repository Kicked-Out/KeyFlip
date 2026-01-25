package core

import (
	"encoding/json"
	"errors"
	"os"
)

// Struct to match the JSON layout file structure
type LayoutsFile struct {
	Layouts map[string]map[string]string `json:"layouts"`
}
// LoadLayouts loads keyboard layouts from the specified JSON file.
func LoadLayouts(path string) (map[string]map[rune]rune, error) {
	// Read the JSON file
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	// Transform JSON data into the desired map structure
	var raw LayoutsFile
	if err := json.Unmarshal(b, &raw); err != nil {
		return nil, err
	}

	if len(raw.Layouts) == 0 {
		return nil, errors.New("no layouts found")
	}
	// Convert string keys/values to rune keys/values
	result := make(map[string]map[rune]rune)
	
	// Iterate through each layout
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
