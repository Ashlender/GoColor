package main

import (
	"fmt"
	"github.com/fatih/color"
	"math/rand"
	"time"
)

func main() {

	screenDrawer()

	paletteRule(paletteGenerator())

	color.BgRGB(paletteGenerator()).Println("            ")

	paletteRGB()
}

func paletteGenerator() (r, g, b int) {

	//rand based on unix-time
	rand.Seed(time.Now().UnixNano())
	r = rand.Intn(255) + 1
	g = rand.Intn(255) + 1
	b = rand.Intn(255) + 1

	return r, g, b
}

// генерит доп цвета по правилу пятиугольника
func paletteRule(r, g, d int) {

	pentaColorArray := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	fmt.Println(pentaColorArray)

}

// берет инфу из правила и составляет палитру
func paletteRGB() {

	color.BgRGB(255, 128, 0).Println("            ", testRgbToHex)
	color.BgRGB(230, 42, 42).Println("            ", testRgbToHex)
}

// ргб в хекс код
func hexCode() string {

	hex := ""

	return hex
}

// Draw the TUI
func screenDrawer() {

	fmt.Println("Welcome to GoWall!")
}

// ==============TEST==================================================
func testRgbToHex() {
	r := 255
	g := 0
	b := 212
	fmt.Printf("#%02x%02x%02x", r, g, b)
}
