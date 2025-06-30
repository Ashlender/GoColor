package pkg

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type ConfigPalette struct {
	Name   string   `json:"name"`
	Colors [][3]int `json:"colors"`
}

func LoadPaletteFromJSON(filename, targetName string) ([]ColorInfo, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("не удалось прочитать %s: %v", filename, err)
	}

	var palettes []ConfigPalette
	if err := json.Unmarshal(data, &palettes); err != nil {
		return nil, fmt.Errorf("ошибка парсинга JSON: %v", err)
	}

	for _, p := range palettes {
		if strings.EqualFold(p.Name, targetName) {
			var result []ColorInfo
			for _, c := range p.Colors {
				if len(c) != 3 {
					continue
				}
				result = append(result, ColorInfo{R: c[0], G: c[1], B: c[2]})
			}
			if len(result) == 0 {
				return nil, fmt.Errorf("палитра '%s' не содержит валидных цветов", targetName)
			}
			return result, nil
		}
	}

	return nil, fmt.Errorf("палитра '%s' не найдена", targetName)
}
