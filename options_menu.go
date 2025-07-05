//go:build !test

package main

import (
	"fmt"
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

func (g *Game) optionsRect() image.Rectangle {
	size := g.iconSize()
	x := g.width - size*6 - HelpMargin*6
	y := g.height - size - HelpMargin
	return image.Rect(x, y, x+size, y+size)
}

func (g *Game) optionsMenuSize() (int, int) {
	labels := []string{
		OptionsMenuTitle,
		"Textures",
		"Vsync",
		"Show Item Names",
		"Show Legends",
		"Use Item Numbers",
		"Icon Size [-] [+]",
		"Smart Rendering",
		"Half Resolution",
		"Auto Low-Res",
		"Linear Filtering",
		"FPS: 60.0",
		"Version: " + ClientVersion,
		"Close",
	}
	maxW := 0
	for _, s := range labels {
		if len(s) > maxW {
			maxW = len(s)
		}
	}
	w := maxW*LabelCharWidth + 4
	h := (len(labels)+1)*OptionsMenuSpacing + 4
	return w, h
}

func (g *Game) optionsMenuRect() image.Rectangle {
	w, h := g.optionsMenuSize()
	scale := g.uiScale()
	w = int(float64(w) * scale)
	h = int(float64(h) * scale)
	x := g.optionsRect().Min.X - w - int(10*scale)
	if x < 0 {
		x = 0
	}
	y := g.optionsRect().Min.Y - h
	if y < 0 {
		y = 0
	}
	return image.Rect(x, y, x+w, y+h)
}

func (g *Game) drawOptionsMenu(dst *ebiten.Image) {
	scale := g.uiScale()
	rect := g.optionsMenuRect()
	w, h := g.optionsMenuSize()
	img := ebiten.NewImage(w, h)
	drawFrame(img, image.Rect(0, 0, w, h))
	ebitenutil.DebugPrintAt(img, OptionsMenuTitle, 6, 6)
	y := 6 + OptionsMenuSpacing

	drawToggle := func(label string, enabled bool) {
		btn := image.Rect(4, y-4, w-4, y-4+22)
		drawButton(img, btn, enabled)
		ebitenutil.DebugPrintAt(img, label, btn.Min.X+6, btn.Min.Y+4)
		y += OptionsMenuSpacing
	}

	drawToggle("Textures", g.textures)
	drawToggle("Vsync", g.vsync)
	drawToggle("Show Item Names", g.showItemNames)
	drawToggle("Show Legends", g.showLegend)
	drawToggle("Use Item Numbers", g.useNumbers)

	label := "Icon Size"
	ebitenutil.DebugPrintAt(img, label, 6, y)
	tw, _ := textDimensions(label)
	bx := 6 + tw + 6
	minus := image.Rect(bx, y-4, bx+20, y-4+22)
	plus := image.Rect(bx+24, y-4, bx+44, y-4+22)
	drawButton(img, minus, false)
	drawPlusMinus(img, minus, true)
	drawButton(img, plus, false)
	drawPlusMinus(img, plus, false)
	y += OptionsMenuSpacing

	drawToggle("Smart Rendering", g.smartRender)
	drawToggle("Half Resolution", g.halfRes)
	drawToggle("Auto Low-Res", g.autoLowRes)
	drawToggle("Linear Filtering", g.linearFilter)

	fps := fmt.Sprintf("FPS: %.1f", ebiten.ActualFPS())
	ebitenutil.DebugPrintAt(img, fps, 6, y)
	y += OptionsMenuSpacing

	ebitenutil.DebugPrintAt(img, "Version: "+ClientVersion, 6, y)
	y += OptionsMenuSpacing

	btn := image.Rect(4, y-4, w-4, y-4+22)
	drawButton(img, btn, true)
	ebitenutil.DebugPrintAt(img, "Close", btn.Min.X+6, btn.Min.Y+4)

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(scale, scale)
	op.GeoM.Translate(float64(rect.Min.X), float64(rect.Min.Y))
	dst.DrawImage(img, op)
}

func (g *Game) clickOptionsMenu(mx, my int) bool {
	rect := g.optionsMenuRect()
	if !rect.Overlaps(image.Rect(mx, my, mx+1, my+1)) {
		return false
	}
	scale := g.uiScale()
	x := int(float64(mx-rect.Min.X) / scale)
	y := int(float64(my-rect.Min.Y) / scale)
	mx = x
	my = y
	w, _ := g.optionsMenuSize()
	y = 6 + OptionsMenuSpacing

	// Textures
	r := image.Rect(4, y-4, w-4, y-4+22)
	if r.Overlaps(image.Rect(mx, my, mx+1, my+1)) {
		g.textures = !g.textures
		g.needsRedraw = true
		return true
	}
	y += OptionsMenuSpacing

	// Vsync
	r = image.Rect(4, y-4, w-4, y-4+22)
	if r.Overlaps(image.Rect(mx, my, mx+1, my+1)) {
		g.vsync = !g.vsync
		ebiten.SetVsyncEnabled(g.vsync)
		g.needsRedraw = true
		return true
	}
	y += OptionsMenuSpacing

	// Show Item Names
	r = image.Rect(4, y-4, w-4, y-4+22)
	if r.Overlaps(image.Rect(mx, my, mx+1, my+1)) {
		g.showItemNames = !g.showItemNames
		g.needsRedraw = true
		return true
	}
	y += OptionsMenuSpacing

	// Show Legends
	r = image.Rect(4, y-4, w-4, y-4+22)
	if r.Overlaps(image.Rect(mx, my, mx+1, my+1)) {
		g.showLegend = !g.showLegend
		g.needsRedraw = true
		return true
	}
	y += OptionsMenuSpacing

	// Use Item Numbers
	r = image.Rect(4, y-4, w-4, y-4+22)
	if r.Overlaps(image.Rect(mx, my, mx+1, my+1)) {
		g.useNumbers = !g.useNumbers
		g.needsRedraw = true
		return true
	}
	y += OptionsMenuSpacing

	// Icon Size buttons
	labelW, _ := textDimensions("Icon Size")
	bx := 6 + labelW + 6
	minus := image.Rect(bx, y-4, bx+20, y-4+22)
	plus := image.Rect(bx+24, y-4, bx+44, y-4+22)
	if minus.Overlaps(image.Rect(mx, my, mx+1, my+1)) {
		if g.iconScale > 0.25 {
			g.iconScale -= 0.25
		}
		g.needsRedraw = true
		return true
	}
	if plus.Overlaps(image.Rect(mx, my, mx+1, my+1)) {
		g.iconScale += 0.25
		g.needsRedraw = true
		return true
	}
	y += OptionsMenuSpacing

	// Smart Rendering
	r = image.Rect(4, y-4, w-4, y-4+22)
	if r.Overlaps(image.Rect(mx, my, mx+1, my+1)) {
		g.smartRender = !g.smartRender
		g.needsRedraw = true
		return true
	}
	y += OptionsMenuSpacing

	// Half Resolution
	r = image.Rect(4, y-4, w-4, y-4+22)
	if r.Overlaps(image.Rect(mx, my, mx+1, my+1)) {
		g.halfRes = !g.halfRes
		g.needsRedraw = true
		return true
	}
	y += OptionsMenuSpacing

	// Auto Low-Res
	r = image.Rect(4, y-4, w-4, y-4+22)
	if r.Overlaps(image.Rect(mx, my, mx+1, my+1)) {
		g.autoLowRes = !g.autoLowRes
		g.needsRedraw = true
		return true
	}
	y += OptionsMenuSpacing

	// Linear Filtering
	r = image.Rect(4, y-4, w-4, y-4+22)
	if r.Overlaps(image.Rect(mx, my, mx+1, my+1)) {
		g.linearFilter = !g.linearFilter
		g.needsRedraw = true
		return true
	}
	y += OptionsMenuSpacing

	// FPS (not clickable)
	y += OptionsMenuSpacing

	// Version (not clickable)
	y += OptionsMenuSpacing

	// Close
	r = image.Rect(4, y-4, w-4, y-4+22)
	if r.Overlaps(image.Rect(mx, my, mx+1, my+1)) {
		g.showOptions = false
		g.needsRedraw = true
		return true
	}

	return true
}
