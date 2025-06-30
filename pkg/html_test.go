package pkg

import (
	"os"
	"strings"
	"testing"
)

func TestSavePaletteToHTML(t *testing.T) {
	file := "test_output.html"
	colors := []ColorInfo{
		{255, 0, 0},
		{0, 255, 0},
		{0, 0, 255},
	}

	err := SavePaletteToHTML(file, colors)
	if err != nil {
		t.Fatalf("save failed: %v", err)
	}
	defer os.Remove(file)

	content, err := os.ReadFile(file)
	if err != nil {
		t.Fatalf("failed to read html: %v", err)
	}

	if !strings.Contains(string(content), "#FF0000") {
		t.Errorf("HTML does not contain expected red block")
	}
}
