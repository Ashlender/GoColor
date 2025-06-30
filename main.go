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
	var baseColor string
	var htmlOut string

	flag.IntVar(&count, "count", 3, "Количество цветов в палитре (2–5)")
	flag.StringVar(&mode, "mode", "rgb", "Режим генерации: rgb или hsv")
	flag.StringVar(&conf, "conf", "", "Имя палитры из palettes.json (заменяет mode/count)")
	flag.StringVar(&baseColor, "base", "", "Базовый цвет в формате R,G,B для RGB-режима")
	flag.StringVar(&htmlOut, "html", "", "Имя HTML-файла для вывода палитры")
	flag.Parse()

	if conf != "" {
		colors, err := loadPaletteFromJSON("palettes.json", conf)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Ошибка:", err)
			os.Exit(1)
		}
		for _, c := range colors {
			printColorLine(c)
		}
		if htmlOut != "" {
			if err := savePaletteToHTML(htmlOut, colors); err != nil {
				fmt.Fprintln(os.Stderr, "Ошибка сохранения HTML:", err)
				os.Exit(1)
			}
			fmt.Println("HTML-файл сохранён:", htmlOut)
		}
		return
	}

	if count < 2 || count > 5 {
		fmt.Fprintln(os.Stderr, "Ошибка: поддерживаются только значения count от 2 до 5")
		os.Exit(1)
	}

	mode = strings.ToLower(mode)
	var colors []ColorInfo

	baseProvided := false
	var baseR, baseG, baseB int

	if baseColor != "" {
		parts := strings.Split(baseColor, ",")
		if len(parts) != 3 {
			fmt.Fprintln(os.Stderr, "Ошибка: base должен быть в формате R,G,B")
			os.Exit(1)
		}
		_, err := fmt.Sscanf(baseColor, "%d,%d,%d", &baseR, &baseG, &baseB)
		if err != nil || baseR < 0 || baseR > 255 || baseG < 0 || baseG > 255 || baseB < 0 || baseB > 255 {
			fmt.Fprintln(os.Stderr, "Ошибка: base содержит некорректные значения RGB")
			os.Exit(1)
		}
		baseProvided = true
	}

	switch mode {
	case "rgb":
		if baseProvided {
			colors = paletteRuleRGB(baseR, baseG, baseB, count)
		} else {
			r, g, b := paletteGenerator()
			colors = paletteRuleRGB(r, g, b, count)
		}
	case "hsv":
		colors = paletteRuleHSV(count)
	default:
		fmt.Fprintln(os.Stderr, "Ошибка: режим должен быть 'rgb' или 'hsv'")
		os.Exit(1)
	}

	for _, c := range colors {
		printColorLine(c)
	}

	if htmlOut != "" {
		if err := savePaletteToHTML(htmlOut, colors); err != nil {
			fmt.Fprintln(os.Stderr, "Ошибка сохранения HTML:", err)
			os.Exit(1)
		}
		fmt.Println("HTML-файл сохранён:", htmlOut)
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
			if len(result) == 0 {
				return nil, fmt.Errorf("палитра '%s' не содержит валидных цветов", targetName)
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
	rand.Seed(time.Now().UnixNano())

	for i := 0; i < count; i++ {
		h := float64(i) * (360.0 / float64(count))
		s := 0.6 + rand.Float64()*0.4 // насыщенность от 0.6 до 1.0
		v := 0.7 + rand.Float64()*0.3 // яркость от 0.7 до 1.0
		r, g, b := HSVtoRGB(h, s, v)
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
	fmt.Print("                ")
	fmt.Print("\x1b[0m")
	fmt.Printf("  | RGB - [ %3d, %3d, %3d ] | HEX - [#%02X%02X%02X]\n",
		c.R, c.G, c.B, c.R, c.G, c.B)
}

func savePaletteToHTML(filename string, colors []ColorInfo) error {
	html := "<!DOCTYPE html><html><head><meta charset='utf-8'><title>Palette</title></head><body style='font-family:sans-serif'>"
	html += "<h2>Сгенерированная палитра:</h2><div style='display:flex;'>"

	for _, c := range colors {
		hex := fmt.Sprintf("#%02X%02X%02X", c.R, c.G, c.B)
		block := fmt.Sprintf(
			"<div style='width:120px;height:120px;background:%s;display:flex;align-items:center;justify-content:center;margin:4px;border:1px solid #000;'>"+
				"<span style='color:#000;background:#fff;padding:2px;font-size:12px;'>%s</span></div>", hex, hex)
		html += block
	}

	html += "</div></body></html>"
	return os.WriteFile(filename, []byte(html), 0644)
}
