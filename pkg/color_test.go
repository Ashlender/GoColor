package pkg

import (
	"testing"
)

func TestClamp(t *testing.T) {
	tests := []struct {
		input    int
		expected int
	}{
		{-100, 0},
		{0, 0},
		{128, 128},
		{255, 255},
		{300, 255},
	}

	for _, tt := range tests {
		got := clamp(tt.input)
		if got != tt.expected {
			t.Errorf("clamp(%d) = %d; want %d", tt.input, got, tt.expected)
		}
	}
}

func TestPaletteRuleRGB(t *testing.T) {
	r, g, b := 100, 150, 200
	colors := PaletteRuleRGB(r, g, b, 3)

	if len(colors) != 3 {
		t.Errorf("expected 3 colors, got %d", len(colors))
	}

	expected := ColorInfo{100, 150, 200}
	if colors[0] != expected {
		t.Errorf("unexpected color: got %v, want %v", colors[0], expected)
	}
}

func TestHSVtoRGB(t *testing.T) {
	tests := []struct {
		h, s, v float64
		expect  ColorInfo
	}{
		{0, 1, 1, ColorInfo{255, 0, 0}},
		{120, 1, 1, ColorInfo{0, 255, 0}},
		{240, 1, 1, ColorInfo{0, 0, 255}},
	}

	for _, tt := range tests {
		r, g, b := HSVtoRGB(tt.h, tt.s, tt.v)
		got := ColorInfo{r, g, b}
		if got != tt.expect {
			t.Errorf("HSV(%v,%v,%v) = %v; want %v", tt.h, tt.s, tt.v, got, tt.expect)
		}
	}
}
