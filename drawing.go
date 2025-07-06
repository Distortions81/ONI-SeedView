package main

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

func drawTextWithBG(dst *ebiten.Image, text string, x, y int, center bool) {
	w, h := textDimensions(text)
	if center {
		x -= w / 2
	}
	vector.DrawFilledRect(dst, float32(x-2), float32(y-2), float32(w+4), float32(h+4), color.RGBA{0, 0, 0, 128}, false)
	drawText(dst, text, x, y, false)
}

// drawTextWithBG draws text with a translucent background.
// Previously this supported arbitrary scaling, but the UI no longer scales,
// so this helper simply draws at the current font size.
func drawTextWithBGScale(dst *ebiten.Image, text string, x, y int, _ float64, center bool) {
	drawTextWithBG(dst, text, x, y, center)
}

func drawTextWithBGBorder(dst *ebiten.Image, text string, x, y int, border color.Color, center bool) {
	w, h := textDimensions(text)
	if center {
		x -= w / 2
	}
	bx := x - 2
	by := y - 2
	bw := w + 4
	bh := h + 4
	vector.DrawFilledRect(dst, float32(bx-1), float32(by-1), float32(bw+2), float32(bh+2), border, false)
	vector.DrawFilledRect(dst, float32(bx), float32(by), float32(bw), float32(bh), color.RGBA{0, 0, 0, 128}, false)
	drawText(dst, text, x, y, false)
}

func drawTextWithBGBorderScale(dst *ebiten.Image, text string, x, y int, border color.Color, _ float64, center bool) {
	drawTextWithBGBorder(dst, text, x, y, border, center)
}

func drawTextScale(dst *ebiten.Image, text string, x, y int, _ float64, center bool) {
	drawText(dst, text, x, y, center)
}

func (g *Game) drawInfoPanel(dst *ebiten.Image, text string, icon *ebiten.Image, x, y int) {
	txtW, txtH := textDimensions(text)
	iconW, iconH := 0, 0
	if icon != nil {
		iconW = InfoIconSize
		iconH = InfoIconSize
	}
	gap := 4
	w := txtW + iconW + gap + 8
	h := txtH
	if iconH > txtH {
		h = iconH
	}
	h += 8
	img := ebiten.NewImage(w, h)
	vector.DrawFilledRect(img, 0, 0, float32(w), float32(h), color.RGBA{0, 0, 0, InfoPanelAlpha}, false)
	vector.StrokeRect(img, 0.5, 0.5, float32(w)-1, float32(h)-1, 1*g.vectorScale(), buttonBorderColor, false)
	if icon != nil {
		opIcon := &ebiten.DrawImageOptions{Filter: g.filterMode()}
		scaleIcon := float64(InfoIconSize) / math.Max(float64(icon.Bounds().Dx()), float64(icon.Bounds().Dy()))
		opIcon.GeoM.Scale(scaleIcon, scaleIcon)
		opIcon.GeoM.Translate(4, float64(h-iconH)/2)
		img.DrawImage(icon, opIcon)
	}
	drawText(img, text, iconW+gap+4, 4, false)
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(x-4), float64(y-4))
	dst.DrawImage(img, op)
}

func (g *Game) drawInfoRow(dst *ebiten.Image, text string, icon *ebiten.Image, x, y int) {
	txtW, txtH := textDimensions(text)
	iconW, iconH := 0, 0
	if icon != nil {
		iconW = InfoIconSize
		iconH = InfoIconSize
	}
	gap := 4
	w := txtW + iconW + gap
	h := txtH
	if iconH > txtH {
		h = iconH
	}
	img := ebiten.NewImage(w, h)
	if icon != nil {
		opIcon := &ebiten.DrawImageOptions{Filter: g.filterMode()}
		sc := float64(InfoIconSize) / math.Max(float64(icon.Bounds().Dx()), float64(icon.Bounds().Dy()))
		opIcon.GeoM.Scale(sc, sc)
		opIcon.GeoM.Translate(0, float64(h-iconH)/2)
		img.DrawImage(icon, opIcon)
	}
	drawText(img, text, iconW+gap, (h-txtH)/2, false)
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(x), float64(y))
	dst.DrawImage(img, op)
}

func infoPanelSize(text string, icon *ebiten.Image) (int, int) {
	txtW, txtH := textDimensions(text)
	iconW, iconH := 0, 0
	if icon != nil {
		iconW = InfoIconSize
		iconH = InfoIconSize
	}
	gap := 4
	w := txtW + iconW + gap + 8
	h := txtH
	if iconH > txtH {
		h = iconH
	}
	h += 8
	return w, h
}

func infoRowSize(text string, icon *ebiten.Image) (int, int) {
	txtW, txtH := textDimensions(text)
	iconW, iconH := 0, 0
	if icon != nil {
		iconW = InfoIconSize
		iconH = InfoIconSize
	}
	gap := 4
	w := txtW + iconW + gap
	h := txtH
	if iconH > txtH {
		h = iconH
	}
	return w, h
}

func drawBiome(dst *ebiten.Image, polys [][]Point, clr color.Color, camX, camY, zoom float64) {
	if len(polys) == 0 {
		return
	}
	var p vector.Path
	for _, pts := range polys {
		if len(pts) == 0 {
			continue
		}
		p.MoveTo(float32(pts[0].X*2), float32(pts[0].Y*2))
		for _, pt := range pts[1:] {
			p.LineTo(float32(pt.X*2), float32(pt.Y*2))
		}
		p.Close()
	}
	vs, is := p.AppendVerticesAndIndicesForFilling(nil, nil)
	r, g0, b, a := clr.RGBA()
	for i := range vs {
		x := float64(vs[i].DstX)*zoom + camX
		y := float64(vs[i].DstY)*zoom + camY
		vs[i].DstX = float32(x)
		vs[i].DstY = float32(y)
		vs[i].SrcX = 0
		vs[i].SrcY = 0
		vs[i].ColorR = float32(r) / 0xffff
		vs[i].ColorG = float32(g0) / 0xffff
		vs[i].ColorB = float32(b) / 0xffff
		vs[i].ColorA = float32(a) / 0xffff
	}
	op := &ebiten.DrawTrianglesOptions{
		AntiAlias:      true,
		ColorScaleMode: ebiten.ColorScaleModePremultipliedAlpha,
		FillRule:       ebiten.FillRuleEvenOdd,
	}
	dst.DrawTriangles(vs, is, whitePixel, op)
}

func drawBiomeTextured(dst *ebiten.Image, polys [][]Point, tex *ebiten.Image, clr color.Color, camX, camY, zoom float64, filter ebiten.Filter) {
	if len(polys) == 0 || tex == nil {
		return
	}
	var p vector.Path
	for _, pts := range polys {
		if len(pts) == 0 {
			continue
		}
		p.MoveTo(float32(pts[0].X*2), float32(pts[0].Y*2))
		for _, pt := range pts[1:] {
			p.LineTo(float32(pt.X*2), float32(pt.Y*2))
		}
		p.Close()
	}
	vs, is := p.AppendVerticesAndIndicesForFilling(nil, nil)
	r, g0, b, a := clr.RGBA()
	texW, texH := tex.Bounds().Dx(), tex.Bounds().Dy()
	for i := range vs {
		worldX := float64(vs[i].DstX) / 2
		worldY := float64(vs[i].DstY) / 2
		x := float64(vs[i].DstX)*zoom + camX
		y := float64(vs[i].DstY)*zoom + camY
		vs[i].DstX = float32(x)
		vs[i].DstY = float32(y)
		vs[i].SrcX = float32(worldX*BiomeTextureScale) * float32(texW)
		vs[i].SrcY = float32(worldY*BiomeTextureScale) * float32(texH)
		vs[i].ColorR = float32(r) / 0xffff
		vs[i].ColorG = float32(g0) / 0xffff
		vs[i].ColorB = float32(b) / 0xffff
		vs[i].ColorA = float32(a) / 0xffff
	}
	op := &ebiten.DrawTrianglesOptions{
		AntiAlias:      true,
		ColorScaleMode: ebiten.ColorScaleModePremultipliedAlpha,
		FillRule:       ebiten.FillRuleEvenOdd,
		Address:        ebiten.AddressRepeat,
		Filter:         filter,
	}
	dst.DrawTriangles(vs, is, tex, op)
}

func drawBiomeOutline(dst *ebiten.Image, polys [][]Point, camX, camY, zoom float64, clr color.Color) {
	for _, pts := range polys {
		if len(pts) < 2 {
			continue
		}
		for i := 0; i < len(pts); i++ {
			a := pts[i]
			b := pts[(i+1)%len(pts)]
			x0 := float32(math.Round(float64(a.X*2)*zoom + camX))
			y0 := float32(math.Round(float64(a.Y*2)*zoom + camY))
			x1 := float32(math.Round(float64(b.X*2)*zoom + camX))
			y1 := float32(math.Round(float64(b.Y*2)*zoom + camY))
			vector.StrokeLine(dst, x0, y0, x1, y1, 1, clr, true)
		}
	}
}
