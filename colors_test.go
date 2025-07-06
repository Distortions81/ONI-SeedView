//go:build test

package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"image/color"
)

var whitePixel = func() *ebiten.Image {
	img := ebiten.NewImage(1, 1)
	img.Fill(color.White)
	return img
}()
