//go:build test

package main

import (
	"image/color"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
)

// Stubbed text helpers for tests. They avoid drawing but keep signatures.
func drawText(dst *ebiten.Image, str string, x, y int) {}

func drawTextColor(dst *ebiten.Image, str string, x, y int, clr color.Color) {}

func textDimensions(str string) (int, int) {
	lines := strings.Split(str, "\n")
	width := 0
	for _, l := range lines {
		if len(l) > width {
			width = len(l)
		}
	}
	return width * LabelCharWidth, len(lines) * 20
}
