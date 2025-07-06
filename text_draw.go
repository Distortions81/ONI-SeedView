//go:build !test

package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
)

func drawText(dst *ebiten.Image, str string, x, y int) {
	if notoFont == nil {
		return
	}
	text.Draw(dst, str, notoFont, x, y+int(notoFont.Metrics().Ascent.Ceil()), color.White)
}

func drawTextColor(dst *ebiten.Image, str string, x, y int, clr color.Color) {
	if notoFont == nil {
		return
	}
	text.Draw(dst, str, notoFont, x, y+int(notoFont.Metrics().Ascent.Ceil()), clr)
}
