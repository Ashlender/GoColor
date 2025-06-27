package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"math"
	"math/rand"
	"os"
	"strings"
	"time"
)

type ColorInfo struct {
	R, G, B int
}

type ConfigPalette struct {
	Name   string   `json:"name"`
	Colors [][3]int `json:"colors"`
}

func main() {
	var count int
	var mode string
	var conf string

	flag.IntVar(&count, "count", 3, "Количество цветов в палитре (2–5)")
	flag.StringVar(&mode, "mode", "rgb", "Режим генерации: rgb или hsv")
	flag.StringVar(&conf, "conf", "", "Имя палитры из palettes.json (заменяет mode/count)")
	flag.Parse()

	if conf != "" {
		colors, err := loadPaletteFromJSON("palettes.json", conf)
		if err != nil {
			fmt.Println("Ошибка:", err)
			os.Exit(1)
		}
		fmt.Printf("Loaded palette: %s\n", conf)
		for _, c := range colors {
			printColorLine(c)
		}
		return
	}

	if count < 2 || count > 5 {
		fmt.Println("Ошибка: поддерживаются только значения count от 2 до 5")
		os.Exit(1)
	}

	mode = strings.ToLower(mode)

	var colors []ColorInfo

	switch mode {
	case "rgb":
		r, g, b := paletteGenerator()
		colors = paletteRuleRGB(r, g, b, count)
	case "hsv":
		colors = paletteRuleHSV(count)
	default:
		fmt.Println("Ошибка: режим должен быть 'rgb' или 'hsv'")
		os.Exit(1)
	}

	fmt.Println("Generated palette:")
	for _, c := range colors {
		printColorLine(c)
	}
}

func loadPaletteFromJSON(filename, targetName string) ([]ColorInfo, error) {
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
				result = append(result, ColorInfo{c[0], c[1], c[2]})
			}
			return result, nil
		}
	}

	return nil, fmt.Errorf("палитра '%s' не найдена", targetName)
}

func paletteGenerator() (r, g, b int) {
	rand.Seed(time.Now().UnixNano())
	r = rand.Intn(256)
	g = rand.Intn(256)
	b = rand.Intn(256)
	return
}

func clamp(x int) int {
	if x < 0 {
		return 0
	}
	if x > 255 {
		return 255
	}
	return x
}

func paletteRuleRGB(r, g, b, count int) []ColorInfo {
	base := ColorInfo{r, g, b}
	colors := []ColorInfo{base}

	switch count {
	case 2:
		colors = append(colors, ColorInfo{255 - r, 255 - g, 255 - b})
	case 3:
		colors = append(colors, ColorInfo{g, b, r})
		colors = append(colors, ColorInfo{b, r, g})
	case 4:
		colors = append(colors, ColorInfo{255 - r, 255 - g, 255 - b})
		colors = append(colors, ColorInfo{g, r, b})
		colors = append(colors, ColorInfo{b, g, r})
	case 5:
		colors = append(colors, ColorInfo{255 - r, 255 - g, 255 - b})
		colors = append(colors, ColorInfo{(r + g) / 2, (g + b) / 2, (b + r) / 2})
		colors = append(colors, ColorInfo{g, r, b})
		colors = append(colors, ColorInfo{b, g, r})
	}

	return colors
}

func paletteRuleHSV(count int) []ColorInfo {
	var colors []ColorInfo

	for i := 0; i < count; i++ {
		h := float64(i) * (360.0 / float64(count))
		r, g, b := HSVtoRGB(h, 1.0, 1.0)
		colors = append(colors, ColorInfo{r, g, b})
	}

	return colors
}

func HSVtoRGB(h, s, v float64) (int, int, int) {
	c := v * s
	x := c * (1 - math.Abs(math.Mod(h/60.0, 2)-1))
	m := v - c

	var r1, g1, b1 float64

	switch {
	case h >= 0 && h < 60:
		r1, g1, b1 = c, x, 0
	case h >= 60 && h < 120:
		r1, g1, b1 = x, c, 0
	case h >= 120 && h < 180:
		r1, g1, b1 = 0, c, x
	case h >= 180 && h < 240:
		r1, g1, b1 = 0, x, c
	case h >= 240 && h < 300:
		r1, g1, b1 = x, 0, c
	case h >= 300 && h < 360:
		r1, g1, b1 = c, 0, x
	default:
		r1, g1, b1 = 0, 0, 0
	}

	r := int((r1 + m) * 255)
	g := int((g1 + m) * 255)
	b := int((b1 + m) * 255)
	return clamp(r), clamp(g), clamp(b)
}

func printColorLine(c ColorInfo) {
	fmt.Printf("\x1b[48;2;%d;%d;%dm", c.R, c.G, c.B)
	fmt.Print("        ")
	fmt.Print("\x1b[0m")
	fmt.Printf("  | RGB - [ %3d, %3d, %3d ] | HEX - [#%02X%02X%02X]\n",
		c.R, c.G, c.B, c.R, c.G, c.B)
}
