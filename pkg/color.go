package pkg

import (
	"math"
	"math/rand"
	"time"
)

type ColorInfo struct {
	R, G, B int
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

// Generation based on UNIX time
func PaletteGenerator() (r, g, b int) {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(256), rand.Intn(256), rand.Intn(256)
}

func PaletteRuleRGB(r, g, b, count int) []ColorInfo {
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

// Generation based on UNIX time
func PaletteRuleHSV(count int) []ColorInfo {
	var colors []ColorInfo
	rand.Seed(time.Now().UnixNano())

	for i := 0; i < count; i++ {
		h := float64(i) * (360.0 / float64(count))
		s := 0.6 + rand.Float64()*0.4
		v := 0.7 + rand.Float64()*0.3
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
	case h < 60:
		r1, g1, b1 = c, x, 0
	case h < 120:
		r1, g1, b1 = x, c, 0
	case h < 180:
		r1, g1, b1 = 0, c, x
	case h < 240:
		r1, g1, b1 = 0, x, c
	case h < 300:
		r1, g1, b1 = x, 0, c
	default:
		r1, g1, b1 = c, 0, x
	}
	return clamp(int((r1 + m) * 255)), clamp(int((g1 + m) * 255)), clamp(int((b1 + m) * 255))
}
