package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

type ColorInfo struct {
	R, G, B int
}

func main() {
	var count int
	var mode string

	flag.IntVar(&count, "count", 3, "Количество цветов в палитре (2, 3, 4 или 5)")
	flag.StringVar(&mode, "mode", "rgb", "Режим генерации: rgb или hsv")
	flag.Parse()

	if count < 2 || count > 5 {
		fmt.Println("Ошибка: поддерживается только count от 2 до 5")
		os.Exit(1)
	}

	mode = strings.ToLower(mode)
	if mode != "rgb" {
		fmt.Println("Ошибка: пока реализован только режим 'rgb'")
		os.Exit(1)
	}

	r, g, b := paletteGenerator()
	colors := paletteRuleRGB(r, g, b, count)

	fmt.Println("Generated palette:")
	for _, c := range colors {
		printColorLine(c)
	}
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

func paletteGenerator() (r, g, b int) {
	rand.Seed(time.Now().UnixNano())
	r = rand.Intn(256)
	g = rand.Intn(256)
	b = rand.Intn(256)
	return
}

func paletteRuleRGB(r, g, b, count int) []ColorInfo {
	base := ColorInfo{r, g, b}
	var colors []ColorInfo
	colors = append(colors, base)

	switch count {
	case 2:
		// комплементарный
		colors = append(colors, ColorInfo{
			255 - r,
			255 - g,
			255 - b,
		})
	case 3:
		// триада: перестановки RGB
		colors = append(colors, ColorInfo{g, b, r})
		colors = append(colors, ColorInfo{b, r, g})
	case 4:
		// тетрада: базовый, комплемент, два акцента
		colors = append(colors, ColorInfo{255 - r, 255 - g, 255 - b})
		colors = append(colors, ColorInfo{g, r, b})
		colors = append(colors, ColorInfo{b, g, r})
	case 5:
		// псевдо-пятиугольник: смешение и противоположные
		colors = append(colors, ColorInfo{255 - r, 255 - g, 255 - b})
		colors = append(colors, ColorInfo{(r + g) / 2, (g + b) / 2, (b + r) / 2})
		colors = append(colors, ColorInfo{g, r, b})
		colors = append(colors, ColorInfo{b, g, r})
	}

	return colors
}

func printColorLine(c ColorInfo) {
	fmt.Printf("\x1b[48;2;%d;%d;%dm", c.R, c.G, c.B)
	fmt.Print("        ")
	fmt.Print("\x1b[0m")
	fmt.Printf("  | RGB - [ %3d, %3d, %3d ] | HEX - [#%02X%02X%02X]\n",
		c.R, c.G, c.B, c.R, c.G, c.B)
}
