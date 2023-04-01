package core

import (
	"image/color"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIntToColor(t *testing.T) {

	// green
	hex := 0x00FF00
	expected := color.RGBA{0, 255, 0, 255}
	got := IntToColor(hex)
	require.Equal(t, expected, got)

	// black
	hex = 0x000000
	expected = color.RGBA{0, 0, 0, 255}
	got = IntToColor(hex)
	require.Equal(t, expected, got)

	// white
	hex = 0xFFFFFF
	expected = color.RGBA{255, 255, 255, 255}
	got = IntToColor(hex)
	require.Equal(t, expected, got)

}

func TestColorToHex(t *testing.T) {

	// magenta
	col := color.RGBA{255, 0, 255, 255}
	hex := 0xFF00FF
	require.Equal(t, hex, ColorToInt(col))

	// black
	col = color.RGBA{0, 0, 0, 255}
	hex = 0x000000
	require.Equal(t, hex, ColorToInt(col))

	// white
	col = color.RGBA{255, 255, 255, 255}
	hex = 0xFFFFFF
	require.Equal(t, hex, ColorToInt(col))
}

func TestHexToColor(t *testing.T) {

	// magenta
	hex := "#ff00ff"
	col := color.RGBA{255, 0, 255, 255}

	c, err := HexToColor(hex)
	require.Nil(t, err)
	require.Equal(t, col, c)

	// black
	hex = "000000"
	col = color.RGBA{0, 0, 0, 255}

	c, err = HexToColor(hex)
	require.Nil(t, err)
	require.Equal(t, col, c)

	// white
	hex = "ffffff"
	col = color.RGBA{255, 255, 255, 255}

	c, err = HexToColor(hex)
	require.Nil(t, err)
	require.Equal(t, col, c)
}
