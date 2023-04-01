package core

import (
	"image/color"
	"strconv"
	"strings"
)

// IntToColor converts a hex int to a Go Color
func IntToColor(hex int) color.Color {
	r := uint8(hex >> 16 & 0xFF)
	g := uint8(hex >> 8 & 0xFF)
	b := uint8(hex & 0xFF)
	return color.RGBA{r, g, b, 255}
}

// ColorToInt converts a Go Color to a hex int
func ColorToInt(c color.Color) int {
	rgba := color.RGBAModel.Convert(c).(color.RGBA)
	hex := int(rgba.R)<<16 | int(rgba.G)<<8 | int(rgba.B)
	return hex
}

// HexToColor coverts hex string to color
func HexToColor(s string) (color.Color, error) {
	s = strings.ReplaceAll(s, "#", "")

	hex, err := strconv.ParseInt(s, 16, 32)
	if err != nil {
		return nil, err
	}
	return IntToColor(int(hex)), nil

}
