package pkg

import (
	"fmt"
	"os"
)

func SavePaletteToHTML(filename string, colors []ColorInfo) error {
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
