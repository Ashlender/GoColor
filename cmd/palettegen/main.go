package main

import (
	"GoWall/pkg"
	"flag"
	"fmt"
	"os"
	"strings"
)

func main() {

	// Flags
	var count int
	var mode string
	var conf string
	var baseColor string
	var htmlOut string

	flag.IntVar(&count, "count", 3, " --colors in the palette [2–5]")
	flag.StringVar(&mode, "mode", "rgb", " --generation mode (standard mode == rgb): [rgb/hsv]")
	flag.StringVar(&conf, "conf", "", " --palette name from palettes.json [replaces mode/count]")
	flag.StringVar(&baseColor, "base", "", " --base color in R,G,B format for [rgb] mode")
	flag.StringVar(&htmlOut, "html", "", " --name of HTML file to output the palette")
	flag.Parse()

	if conf != "" {
		colors, err := pkg.LoadPaletteFromJSON("palettes.json", conf)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error:", err)
			os.Exit(1)
		}
		printColors(colors)
		saveHTML(htmlOut, colors)
		return
	}

	if count < 2 || count > 5 {
		fmt.Fprintln(os.Stderr, "Error: Only count values [2 - 5] are supported")
		os.Exit(1)
	}

	mode = strings.ToLower(mode)
	var colors []pkg.ColorInfo
	var baseProvided bool
	var baseR, baseG, baseB int

	if baseColor != "" {
		_, err := fmt.Sscanf(baseColor, "%d,%d,%d", &baseR, &baseG, &baseB)
		if err != nil || baseR < 0 || baseR > 255 || baseG < 0 || baseG > 255 || baseB < 0 || baseB > 255 {
			fmt.Fprintln(os.Stderr, "Error: --base contains invalid [rgb] values")
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
		fmt.Fprintln(os.Stderr, "Error: --mode must be [None] or [hsv]")
		os.Exit(1)
	}

	printColors(colors)
	saveHTML(htmlOut, colors)
}

// Print function
func printColors(colors []pkg.ColorInfo) {
	for _, c := range colors {
		fmt.Printf("\x1b[48;2;%d;%d;%dm", c.R, c.G, c.B)
		fmt.Print("                ")
		fmt.Print("\x1b[0m")
		fmt.Printf("  | RGB - [ %3d, %3d, %3d ] | HEX - [#%02X%02X%02X]\n",
			c.R, c.G, c.B, c.R, c.G, c.B)
	}
}

// Saving HTML-config
func saveHTML(htmlOut string, colors []pkg.ColorInfo) {
	if htmlOut != "" {
		if err := pkg.SavePaletteToHTML(htmlOut, colors); err != nil {
			fmt.Fprintln(os.Stderr, "Error while saving HTML-config:", err)
			os.Exit(1)
		}
		fmt.Println("HTML-config saved:", htmlOut)
	}
}
