//go:build test

package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"image/color"
)

func drawText(dst *ebiten.Image, str string, x, y int) {
	ebitenutil.DebugPrintAt(dst, str, x, y)
}

func drawTextColor(dst *ebiten.Image, str string, x, y int, clr color.Color) {
	ebitenutil.DebugPrintAt(dst, str, x, y)
}
