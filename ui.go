package main

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// drawFrame draws a simple framed rectangle with the global frame color.
func drawFrame(dst *ebiten.Image, rect image.Rectangle) {
	vector.DrawFilledRect(dst, float32(rect.Min.X), float32(rect.Min.Y), float32(rect.Dx()), float32(rect.Dy()), frameColor, false)
	vector.StrokeRect(dst, float32(rect.Min.X)+0.5, float32(rect.Min.Y)+0.5, float32(rect.Dx())-1, float32(rect.Dy())-1, 1, buttonBorderColor, false)
}

// strokeFrame draws only the border using the global frame color without
// filling the rectangle.
func strokeFrame(dst *ebiten.Image, rect image.Rectangle) {
	vector.StrokeRect(dst, float32(rect.Min.X)+0.5, float32(rect.Min.Y)+0.5, float32(rect.Dx())-1, float32(rect.Dy())-1, 1, buttonBorderColor, false)
}

// drawButton draws a rounded button. When active, the cyan color is used.
func drawButton(dst *ebiten.Image, rect image.Rectangle, active bool) {
	clr := buttonInactiveColor
	if active {
		clr = buttonActiveColor
	}
	radius := float32(4)
	drawRoundedRect(dst, rect, radius, clr, buttonBorderColor)
}

// drawCloseButton draws a close box with an X.
func drawCloseButton(dst *ebiten.Image, rect image.Rectangle) {
	drawButton(dst, rect, false)
	pad := float32(rect.Dx()) * 0.25
	vector.StrokeLine(dst, float32(rect.Min.X)+pad, float32(rect.Min.Y)+pad, float32(rect.Max.X)-pad, float32(rect.Max.Y)-pad, 2, buttonBorderColor, true)
	vector.StrokeLine(dst, float32(rect.Min.X)+pad, float32(rect.Max.Y)-pad, float32(rect.Max.X)-pad, float32(rect.Min.Y)+pad, 2, buttonBorderColor, true)
}

// drawRoundedRect fills and strokes a rectangle with a corner radius.
func drawRoundedRect(dst *ebiten.Image, rect image.Rectangle, radius float32, fill, border color.Color) {
	x0 := float32(rect.Min.X)
	y0 := float32(rect.Min.Y)
	w := float32(rect.Dx())
	h := float32(rect.Dy())
	r := radius
	var p vector.Path
	p.MoveTo(x0+r, y0)
	p.LineTo(x0+w-r, y0)
	p.QuadTo(x0+w, y0, x0+w, y0+r)
	p.LineTo(x0+w, y0+h-r)
	p.QuadTo(x0+w, y0+h, x0+w-r, y0+h)
	p.LineTo(x0+r, y0+h)
	p.QuadTo(x0, y0+h, x0, y0+h-r)
	p.LineTo(x0, y0+r)
	p.QuadTo(x0, y0, x0+r, y0)
	p.Close()
	vs, is := p.AppendVerticesAndIndicesForFilling(nil, nil)
	r0, g0, b0, a0 := fill.RGBA()
	for i := range vs {
		vs[i].ColorR = float32(r0) / 0xffff
		vs[i].ColorG = float32(g0) / 0xffff
		vs[i].ColorB = float32(b0) / 0xffff
		vs[i].ColorA = float32(a0) / 0xffff
	}
	op := &ebiten.DrawTrianglesOptions{AntiAlias: true, ColorScaleMode: ebiten.ColorScaleModePremultipliedAlpha}
	dst.DrawTriangles(vs, is, whitePixel, op)

	vector.StrokeLine(dst, x0+r, y0, x0+w-r, y0, 1, border, true)
	vector.StrokeLine(dst, x0+w, y0+r, x0+w, y0+h-r, 1, border, true)
	vector.StrokeLine(dst, x0+r, y0+h, x0+w-r, y0+h, 1, border, true)
	vector.StrokeLine(dst, x0, y0+r, x0, y0+h-r, 1, border, true)
	vector.StrokeLine(dst, x0+r, y0, x0, y0+r, 1, border, true)
	vector.StrokeLine(dst, x0+w-r, y0, x0+w, y0+r, 1, border, true)
	vector.StrokeLine(dst, x0, y0+h-r, x0+r, y0+h, 1, border, true)
	vector.StrokeLine(dst, x0+w, y0+h-r, x0+w-r, y0+h, 1, border, true)
}
