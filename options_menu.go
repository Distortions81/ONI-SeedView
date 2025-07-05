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
	x := g.width - size*5 - HelpMargin*5
	y := g.height - size - HelpMargin
	return image.Rect(x, y, x+size, y+size)
}

func (g *Game) optionsMenuRect() image.Rectangle {
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
	x := g.optionsRect().Min.X - w - 10
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
	rect := g.optionsMenuRect()
	drawFrame(dst, rect)
	ebitenutil.DebugPrintAt(dst, OptionsMenuTitle, rect.Min.X+6, rect.Min.Y+6)
	y := rect.Min.Y + 6 + OptionsMenuSpacing

	drawToggle := func(label string, enabled bool) {
		btn := image.Rect(rect.Min.X+4, y-4, rect.Max.X-4, y-4+22)
		drawButton(dst, btn, enabled)
		ebitenutil.DebugPrintAt(dst, label, btn.Min.X+6, btn.Min.Y+4)
		y += OptionsMenuSpacing
	}

	drawToggle("Textures", g.textures)
	drawToggle("Vsync", g.vsync)
	drawToggle("Show Item Names", g.showItemNames)
	drawToggle("Show Legends", g.showLegend)
	drawToggle("Use Item Numbers", g.useNumbers)

	label := "Icon Size"
	ebitenutil.DebugPrintAt(dst, label, rect.Min.X+6, y)
	w, _ := textDimensions(label)
	bx := rect.Min.X + 6 + w + 6
	minus := image.Rect(bx, y-4, bx+20, y-4+22)
	plus := image.Rect(bx+24, y-4, bx+44, y-4+22)
	drawButton(dst, minus, false)
	drawPlusMinus(dst, minus, true)
	drawButton(dst, plus, false)
	drawPlusMinus(dst, plus, false)
	y += OptionsMenuSpacing

	drawToggle("Smart Rendering", g.smartRender)
	drawToggle("Half Resolution", g.halfRes)
	drawToggle("Auto Low-Res", g.autoLowRes)
	drawToggle("Linear Filtering", g.linearFilter)

	fps := fmt.Sprintf("FPS: %.1f", ebiten.ActualFPS())
	ebitenutil.DebugPrintAt(dst, fps, rect.Min.X+6, y)
	y += OptionsMenuSpacing

	ebitenutil.DebugPrintAt(dst, "Version: "+ClientVersion, rect.Min.X+6, y)
	y += OptionsMenuSpacing

	btn := image.Rect(rect.Min.X+4, y-4, rect.Max.X-4, y-4+22)
	drawButton(dst, btn, true)
	ebitenutil.DebugPrintAt(dst, "Close", btn.Min.X+6, btn.Min.Y+4)
}

func (g *Game) clickOptionsMenu(mx, my int) bool {
	rect := g.optionsMenuRect()
	if !rect.Overlaps(image.Rect(mx, my, mx+1, my+1)) {
		return false
	}
	y := rect.Min.Y + 6 + OptionsMenuSpacing

	// Textures
	r := image.Rect(rect.Min.X+4, y-4, rect.Max.X-4, y-4+22)
	if r.Overlaps(image.Rect(mx, my, mx+1, my+1)) {
		g.textures = !g.textures
		g.needsRedraw = true
		return true
	}
	y += OptionsMenuSpacing

	// Vsync
	r = image.Rect(rect.Min.X+4, y-4, rect.Max.X-4, y-4+22)
	if r.Overlaps(image.Rect(mx, my, mx+1, my+1)) {
		g.vsync = !g.vsync
		ebiten.SetVsyncEnabled(g.vsync)
		g.needsRedraw = true
		return true
	}
	y += OptionsMenuSpacing

	// Show Item Names
	r = image.Rect(rect.Min.X+4, y-4, rect.Max.X-4, y-4+22)
	if r.Overlaps(image.Rect(mx, my, mx+1, my+1)) {
		g.showItemNames = !g.showItemNames
		g.needsRedraw = true
		return true
	}
	y += OptionsMenuSpacing

	// Show Legends
	r = image.Rect(rect.Min.X+4, y-4, rect.Max.X-4, y-4+22)
	if r.Overlaps(image.Rect(mx, my, mx+1, my+1)) {
		g.showLegend = !g.showLegend
		g.needsRedraw = true
		return true
	}
	y += OptionsMenuSpacing

	// Use Item Numbers
	r = image.Rect(rect.Min.X+4, y-4, rect.Max.X-4, y-4+22)
	if r.Overlaps(image.Rect(mx, my, mx+1, my+1)) {
		g.useNumbers = !g.useNumbers
		g.needsRedraw = true
		return true
	}
	y += OptionsMenuSpacing

	// Icon Size buttons
	labelW, _ := textDimensions("Icon Size")
	bx := rect.Min.X + 6 + labelW + 6
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
	r = image.Rect(rect.Min.X+4, y-4, rect.Max.X-4, y-4+22)
	if r.Overlaps(image.Rect(mx, my, mx+1, my+1)) {
		g.smartRender = !g.smartRender
		g.needsRedraw = true
		return true
	}
	y += OptionsMenuSpacing

	// Half Resolution
	r = image.Rect(rect.Min.X+4, y-4, rect.Max.X-4, y-4+22)
	if r.Overlaps(image.Rect(mx, my, mx+1, my+1)) {
		g.halfRes = !g.halfRes
		g.needsRedraw = true
		return true
	}
	y += OptionsMenuSpacing

	// Auto Low-Res
	r = image.Rect(rect.Min.X+4, y-4, rect.Max.X-4, y-4+22)
	if r.Overlaps(image.Rect(mx, my, mx+1, my+1)) {
		g.autoLowRes = !g.autoLowRes
		g.needsRedraw = true
		return true
	}
	y += OptionsMenuSpacing

	// Linear Filtering
	r = image.Rect(rect.Min.X+4, y-4, rect.Max.X-4, y-4+22)
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
	r = image.Rect(rect.Min.X+4, y-4, rect.Max.X-4, y-4+22)
	if r.Overlaps(image.Rect(mx, my, mx+1, my+1)) {
		g.showOptions = false
		g.needsRedraw = true
		return true
	}

	return true
}
