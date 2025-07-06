package main

import (
	"image"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

func drawTextWithBG(dst *ebiten.Image, text string, x, y int) {
	w, h := textDimensions(text)
	vector.DrawFilledRect(dst, float32(x-2), float32(y-2), float32(w+4), float32(h+4), color.RGBA{0, 0, 0, 128}, false)
	drawText(dst, text, x, y)
}

func drawTextWithBGScale(dst *ebiten.Image, text string, x, y int, scale float64) {
	if scale == 1.0 {
		drawTextWithBG(dst, text, x, y)
		return
	}
	w, h := textDimensions(text)
	w += 4
	h += 4
	img := ebiten.NewImage(w, h)
	vector.DrawFilledRect(img, 0, 0, float32(w), float32(h), color.RGBA{0, 0, 0, 128}, false)
	drawText(img, text, 2, 2)
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(scale, scale)
	op.GeoM.Translate(float64(x-2), float64(y-2))
	dst.DrawImage(img, op)
}

func drawTextWithBGBorder(dst *ebiten.Image, text string, x, y int, border color.Color) {
	w, h := textDimensions(text)
	bx := x - 2
	by := y - 2
	bw := w + 4
	bh := h + 4
	vector.DrawFilledRect(dst, float32(bx-1), float32(by-1), float32(bw+2), float32(bh+2), border, false)
	vector.DrawFilledRect(dst, float32(bx), float32(by), float32(bw), float32(bh), color.RGBA{0, 0, 0, 128}, false)
	drawText(dst, text, x, y)
}

func drawTextWithBGBorderScale(dst *ebiten.Image, text string, x, y int, border color.Color, scale float64) {
	if scale == 1.0 {
		drawTextWithBGBorder(dst, text, x, y, border)
		return
	}
	w, h := textDimensions(text)
	w += 4
	h += 4
	img := ebiten.NewImage(w+2, h+2)
	vector.DrawFilledRect(img, 0, 0, float32(w+2), float32(h+2), border, false)
	vector.DrawFilledRect(img, 1, 1, float32(w), float32(h), color.RGBA{0, 0, 0, 128}, false)
	drawText(img, text, 3, 3)
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(scale, scale)
	op.GeoM.Translate(float64(x-2), float64(y-2))
	dst.DrawImage(img, op)
}

func drawTextScale(dst *ebiten.Image, text string, x, y int, scale float64) {
	if scale == 1.0 {
		drawText(dst, text, x, y)
		return
	}
	w, h := textDimensions(text)
	img := ebiten.NewImage(w, h)
	drawText(img, text, 0, 0)
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(scale, scale)
	op.GeoM.Translate(float64(x), float64(y))
	dst.DrawImage(img, op)
}

func (g *Game) drawInfoPanel(dst *ebiten.Image, text string, icon *ebiten.Image, x, y int, scale float64) {
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
	vector.StrokeRect(img, 0.5, 0.5, float32(w)-1, float32(h)-1, 1, buttonBorderColor, false)
	if icon != nil {
		opIcon := &ebiten.DrawImageOptions{Filter: g.filterMode()}
		scaleIcon := float64(InfoIconSize) / math.Max(float64(icon.Bounds().Dx()), float64(icon.Bounds().Dy()))
		opIcon.GeoM.Scale(scaleIcon, scaleIcon)
		opIcon.GeoM.Translate(4, float64(h-iconH)/2)
		img.DrawImage(icon, opIcon)
	}
	drawText(img, text, iconW+gap+4, 4)
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(scale, scale)
	op.GeoM.Translate(float64(x-4), float64(y-4))
	dst.DrawImage(img, op)
}

func (g *Game) drawInfoRow(dst *ebiten.Image, text string, icon *ebiten.Image, x, y int, scale float64) {
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
	if scale == 1.0 {
		if icon != nil {
			opIcon := &ebiten.DrawImageOptions{Filter: g.filterMode()}
			sc := float64(InfoIconSize) / math.Max(float64(icon.Bounds().Dx()), float64(icon.Bounds().Dy()))
			opIcon.GeoM.Scale(sc, sc)
			opIcon.GeoM.Translate(float64(x), float64(y+(h-iconH)/2))
			dst.DrawImage(icon, opIcon)
		}
		drawText(dst, text, x+iconW+gap, y+(h-txtH)/2)
		return
	}

	img := ebiten.NewImage(w, h)
	if icon != nil {
		opIcon := &ebiten.DrawImageOptions{Filter: g.filterMode()}
		sc := float64(InfoIconSize) / math.Max(float64(icon.Bounds().Dx()), float64(icon.Bounds().Dy()))
		opIcon.GeoM.Scale(sc, sc)
		opIcon.GeoM.Translate(0, float64(h-iconH)/2)
		img.DrawImage(icon, opIcon)
	}
	drawText(img, text, iconW+gap, (h-txtH)/2)
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(scale, scale)
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
