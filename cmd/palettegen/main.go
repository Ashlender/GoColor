package main

import (
	"GoWall/pkg"
	"flag"
	"fmt"
	"os"
	"strings"
)

func main() {
	var count int
	var mode, conf, baseColor, htmlOut string

	flag.IntVar(&count, "count", 3, "Количество цветов в палитре (2–5)")
	flag.StringVar(&mode, "mode", "rgb", "Режим генерации: rgb или hsv")
	flag.StringVar(&conf, "conf", "", "Имя палитры из palettes.json (заменяет mode/count)")
	flag.StringVar(&baseColor, "base", "", "Базовый цвет в формате R,G,B для RGB-режима")
	flag.StringVar(&htmlOut, "html", "", "Имя HTML-файла для вывода палитры")
	flag.Parse()

	if conf != "" {
		colors, err := pkg.LoadPaletteFromJSON("palettes.json", conf)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Ошибка:", err)
			os.Exit(1)
		}
		printColors(colors)
		saveHTML(htmlOut, colors)
		return
	}

	if count < 2 || count > 5 {
		fmt.Fprintln(os.Stderr, "Ошибка: поддерживаются только значения count от 2 до 5")
		os.Exit(1)
	}

	mode = strings.ToLower(mode)
	var colors []pkg.ColorInfo
	var baseProvided bool
	var baseR, baseG, baseB int

	if baseColor != "" {
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
			colors = pkg.PaletteRuleRGB(baseR, baseG, baseB, count)
		} else {
			r, g, b := pkg.PaletteGenerator()
			colors = pkg.PaletteRuleRGB(r, g, b, count)
		}
	case "hsv":
		colors = pkg.PaletteRuleHSV(count)
	default:
		fmt.Fprintln(os.Stderr, "Ошибка: режим должен быть 'rgb' или 'hsv'")
		os.Exit(1)
	}

	printColors(colors)
	saveHTML(htmlOut, colors)
}

func printColors(colors []pkg.ColorInfo) {
	for _, c := range colors {
		fmt.Printf("\x1b[48;2;%d;%d;%dm", c.R, c.G, c.B)
		fmt.Print("                ")
		fmt.Print("\x1b[0m")
		fmt.Printf("  | RGB - [ %3d, %3d, %3d ] | HEX - [#%02X%02X%02X]\n",
			c.R, c.G, c.B, c.R, c.G, c.B)
	}
}

func saveHTML(htmlOut string, colors []pkg.ColorInfo) {
	if htmlOut != "" {
		if err := pkg.SavePaletteToHTML(htmlOut, colors); err != nil {
			fmt.Fprintln(os.Stderr, "Ошибка сохранения HTML:", err)
			os.Exit(1)
		}
		fmt.Println("HTML-файл сохранён:", htmlOut)
	}
}
