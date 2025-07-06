package bmp

import (
	"image"
	"io"

	"golang.org/x/image/bmp"
)

// Encode writes the Image m to w in BMP format.
func Encode(w io.Writer, m image.Image) error {
	return bmp.Encode(w, m)
}
