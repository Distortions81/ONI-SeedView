//go:build !test

package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
)

func drawText(dst *ebiten.Image, str string, x, y int, center bool) {
	if notoFont == nil {
		return
	}
	if center {
		w, _ := textDimensions(str)
		x -= w / 2
	}
	text.Draw(dst, str, notoFont, x, y+int(notoFont.Metrics().Ascent.Ceil()), color.White)
}

func drawTextColor(dst *ebiten.Image, str string, x, y int, center bool, clr color.Color) {
	if notoFont == nil {
		return
	}
	if center {
		w, _ := textDimensions(str)
		x -= w / 2
	}
	text.Draw(dst, str, notoFont, x, y+int(notoFont.Metrics().Ascent.Ceil()), clr)
}
