package main

import (
	"fmt"
	"math/rand"
	"time"
)

type ColorInfo struct {
	R, G, B int
}

func main() {
	r, g, b := paletteGenerator()
	paletteRule(r, g, b)
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

	//rand based on unix-time
	rand.Seed(time.Now().UnixNano())
	r = rand.Intn(255) + 1
	g = rand.Intn(255) + 1
	b = rand.Intn(255) + 1

	return
}

func paletteRule(r, g, b int) {

	colors := []ColorInfo{{r, g, b}}

	for i := 1; i < 5; i++ {
		offset := i * 51
		colors = append(colors, ColorInfo{
			clamp((r + offset) % 256),
			clamp((g + offset) % 256),
			clamp((b + offset) % 256),
		})
	}
	for _, c := range colors {
		printColorLine(c)
	}
}

func printColorLine(c ColorInfo) {
	fmt.Printf("\x1b[48;2;%d;%d;%dm", c.R, c.G, c.B)
	fmt.Print("            ")
	fmt.Print("\x1b[0m")
	fmt.Printf("  | RGB - [ %3d, %3d, %3d ] | HEX - [#%02X%02X%02X]\n", c.R, c.G, c.B, c.R, c.G, c.B)
}
