package pokedraw

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	_ "image/png"
	"strings"
)

func DisplayImage(data []byte) error {
	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return err
	}

	ramp := " .=+#@"
	max := img.Bounds().Max
	scaleX, scaleY := 16, 8

	var asciiArt []string
	for y := 0; y < max.Y; y += scaleX {
		var row strings.Builder
		for x := 0; x < max.X; x += scaleY {
			c := avgPixel(img, x, y, scaleX, scaleY)
			r, g, b := getRGB(img.At(x, y))
			symbol := string(ramp[len(ramp)*c/65536])
			row.WriteString(colorizeSymbol(symbol, r, g, b))
		}
		asciiArt = append(asciiArt, row.String())
	}

	trimmedArt := trimAsciiArt(asciiArt)
	for _, line := range trimmedArt {
		if line == trimmedArt[0] {
			continue
		}
		fmt.Println(line)
	}

	return nil
}

func avgPixel(img image.Image, x, y, w, h int) int {
	cnt, sum, max := 0, 0, img.Bounds().Max
	for i := x; i < x+w && i < max.X; i++ {
		for j := y; j < y+h && j < max.Y; j++ {
			sum += grayscale(img.At(i, j))
			cnt++
		}
	}
	return sum / cnt
}

func grayscale(c color.Color) int {
	r, g, b, _ := c.RGBA()
	return int((r + g + b) / 3)
}

func getRGB(c color.Color) (int, int, int) {
	r, g, b, _ := c.RGBA()
	return int(r / 256), int(g / 256), int(b / 256)
}

func colorizeSymbol(symbol string, r, g, b int) string {
	return fmt.Sprintf("\x1b[38;2;%d;%d;%dm%s\x1b[0m", r, g, b, symbol)
}

func trimAsciiArt(asciiArt []string) []string {
	top, bottom := 0, len(asciiArt)-1
	left, right := -1, -1

	for i, line := range asciiArt {
		if strings.TrimSpace(line) != "" {
			if top == 0 {
				top = i
			}
			bottom = i
		}
	}

	for _, line := range asciiArt[top : bottom+1] {
		for i, char := range line {
			if char != ' ' {
				if left == -1 || i < left {
					left = i
				}
				if i > right {
					right = i
				}
			}
		}
	}

	if left == -1 || right == -1 {
		return asciiArt[top:bottom]
	}

	if left < 0 {
		left = 0
	}
	if right >= len(asciiArt[0]) {
		right = len(asciiArt[0]) - 1
	}

	var trimmedArt []string
	for _, line := range asciiArt[top : bottom+1] {
		if right < len(line) {
			trimmedArt = append(trimmedArt, line[left:right+1])
		} else {
			trimmedArt = append(trimmedArt, line[left:])
		}
	}

	return trimmedArt
}
