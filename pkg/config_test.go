package pkg

import (
	"encoding/json"
	"os"
	"testing"
)

func TestLoadPaletteFromJSON(t *testing.T) {
	const file = "temp_test_palette.json"
	data := `[{"name": "test", "colors": [[1, 2, 3], [4, 5, 6]]}]`
	err := os.WriteFile(file, []byte(data), 0644)
	if err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}
	defer os.Remove(file)

	colors, err := LoadPaletteFromJSON(file, "test")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(colors) != 2 || colors[0] != (ColorInfo{1, 2, 3}) {
		t.Errorf("incorrect palette: got %v", colors)
	}
}

func TestMarshalUnmarshalConfigPalette(t *testing.T) {
	p := ConfigPalette{
		Name: "sample",
		Colors: [][3]int{
			{10, 20, 30},
			{40, 50, 60},
		},
	}

	data, err := json.Marshal(p)
	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}

	var out ConfigPalette
	if err := json.Unmarshal(data, &out); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}

	if out.Name != p.Name || len(out.Colors) != len(p.Colors) {
		t.Errorf("unmarshal mismatch: got %v, want %v", out, p)
	}
}
