package main

import (
	"os"
	"path/filepath"
	"testing"

	"GoWall/pkg"
)

func TestMainFunctionality(t *testing.T) {
	// Test input data
	r, g, b := 100, 150, 200
	numColors := 4
	outputFile := "test_output.html"

	// Generating a palette
	colors := pkg.PaletteRuleRGB(r, g, b, numColors)
	if len(colors) != numColors {
		t.Fatalf("needed %d colors, got %d", numColors, len(colors))
	}

	// Saving HTML config
	err := pkg.SavePaletteToHTML(outputFile, colors)
	if err != nil {
		t.Fatalf("Error while saving HTML: %v", err)
	}

	// Checking file existence
	if _, err := os.Stat(outputFile); os.IsNotExist(err) {
		t.Fatalf("file %s not exist", outputFile)
	}

	os.Remove(outputFile)
}

func TestLoadPaletteFromJSON(t *testing.T) {
	paletteFile := filepath.Join(".", "palettes.json")
	name := "solarized"

	colors, err := pkg.LoadPaletteFromJSON(paletteFile, name)
	if err != nil {
		t.Fatalf("Error JSON: %v", err)
	}

	if len(colors) == 0 {
		t.Errorf("palette %q empty", name)
	}
}
