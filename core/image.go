package core

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"log"
)

func ImageToBase64(img image.Image) string {
	data := &bytes.Buffer{}
	if err := jpeg.Encode(data, img, nil); err != nil {
		log.Println("Error encoding avatar: ", err)
		return ""
	}

	return fmt.Sprintf("data:image/png;base64,%s", base64.StdEncoding.EncodeToString(data.Bytes()))
}

// ColorToImage creates a 16x16 image filled with the specified color.
func ColorToImage(c color.Color) image.Image {
	// Create a new rectangle with the desired dimensions.
	rect := image.Rect(0, 0, 16, 16)
	// Create a new RGBA image with the specified rectangle.
	img := image.NewRGBA(rect)
	// Use the draw.Draw function to fill the image with the specified color.
	draw.Draw(img, img.Bounds(), &image.Uniform{C: c}, image.Point{}, draw.Src)
	return img
}
