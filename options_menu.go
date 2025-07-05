//go:build !test

package main

import (
	"fmt"
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
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
	vector.DrawFilledRect(dst, float32(rect.Min.X), float32(rect.Min.Y), float32(rect.Dx()), float32(rect.Dy()), colorRGBA(0, 0, 0, 200), false)
	drawTextWithBG(dst, OptionsMenuTitle, rect.Min.X+2, rect.Min.Y+2)
	y := rect.Min.Y + 2 + OptionsMenuSpacing

	drawToggle := func(label string, enabled bool) {
		if enabled {
			drawTextWithBGBorder(dst, label, rect.Min.X+2, y, color.RGBA{255, 255, 255, 255})
		} else {
			drawTextWithBG(dst, label, rect.Min.X+2, y)
		}
		y += OptionsMenuSpacing
	}

	drawToggle("Textures", g.textures)
	drawToggle("Vsync", g.vsync)
	drawToggle("Show Item Names", g.showItemNames)
	drawToggle("Show Legends", g.showLegend)
	drawToggle("Use Item Numbers", g.useNumbers)

	label := "Icon Size"
	drawTextWithBG(dst, label, rect.Min.X+2, y)
	w, _ := textDimensions(label)
	bx := rect.Min.X + 2 + w + 6
	minus := image.Rect(bx, y-2, bx+16, y-2+20)
	plus := image.Rect(bx+20, y-2, bx+36, y-2+20)
	vector.DrawFilledRect(dst, float32(minus.Min.X), float32(minus.Min.Y), float32(minus.Dx()), float32(minus.Dy()), colorRGBA(0, 0, 0, 180), false)
	vector.StrokeRect(dst, float32(minus.Min.X)+0.5, float32(minus.Min.Y)+0.5, float32(minus.Dx())-1, float32(minus.Dy())-1, 2, color.RGBA{255, 255, 255, 255}, false)
	drawPlusMinus(dst, minus, true)
	vector.DrawFilledRect(dst, float32(plus.Min.X), float32(plus.Min.Y), float32(plus.Dx()), float32(plus.Dy()), colorRGBA(0, 0, 0, 180), false)
	vector.StrokeRect(dst, float32(plus.Min.X)+0.5, float32(plus.Min.Y)+0.5, float32(plus.Dx())-1, float32(plus.Dy())-1, 2, color.RGBA{255, 255, 255, 255}, false)
	drawPlusMinus(dst, plus, false)
	y += OptionsMenuSpacing

	drawToggle("Smart Rendering", g.smartRender)

	fps := fmt.Sprintf("FPS: %.1f", ebiten.ActualFPS())
	drawTextWithBG(dst, fps, rect.Min.X+2, y)
	y += OptionsMenuSpacing

	drawTextWithBG(dst, "Version: "+ClientVersion, rect.Min.X+2, y)
	y += OptionsMenuSpacing

	drawTextWithBGBorder(dst, "Close", rect.Min.X+2, y, color.RGBA{255, 255, 255, 255})
}

func (g *Game) clickOptionsMenu(mx, my int) bool {
	rect := g.optionsMenuRect()
	if !rect.Overlaps(image.Rect(mx, my, mx+1, my+1)) {
		return false
	}
	y := rect.Min.Y + 2 + OptionsMenuSpacing

	// Textures
	r := image.Rect(rect.Min.X, y-2, rect.Min.X+rect.Dx(), y-2+20)
	if r.Overlaps(image.Rect(mx, my, mx+1, my+1)) {
		g.textures = !g.textures
		g.needsRedraw = true
		return true
	}
	y += OptionsMenuSpacing

	// Vsync
	r = image.Rect(rect.Min.X, y-2, rect.Min.X+rect.Dx(), y-2+20)
	if r.Overlaps(image.Rect(mx, my, mx+1, my+1)) {
		g.vsync = !g.vsync
		ebiten.SetVsyncEnabled(g.vsync)
		g.needsRedraw = true
		return true
	}
	y += OptionsMenuSpacing

	// Show Item Names
	r = image.Rect(rect.Min.X, y-2, rect.Min.X+rect.Dx(), y-2+20)
	if r.Overlaps(image.Rect(mx, my, mx+1, my+1)) {
		g.showItemNames = !g.showItemNames
		g.needsRedraw = true
		return true
	}
	y += OptionsMenuSpacing

	// Show Legends
	r = image.Rect(rect.Min.X, y-2, rect.Min.X+rect.Dx(), y-2+20)
	if r.Overlaps(image.Rect(mx, my, mx+1, my+1)) {
		g.showLegend = !g.showLegend
		g.needsRedraw = true
		return true
	}
	y += OptionsMenuSpacing

	// Use Item Numbers
	r = image.Rect(rect.Min.X, y-2, rect.Min.X+rect.Dx(), y-2+20)
	if r.Overlaps(image.Rect(mx, my, mx+1, my+1)) {
		g.useNumbers = !g.useNumbers
		g.needsRedraw = true
		return true
	}
	y += OptionsMenuSpacing

	// Icon Size buttons
	labelW, _ := textDimensions("Icon Size")
	bx := rect.Min.X + 2 + labelW + 6
	minus := image.Rect(bx, y-2, bx+16, y-2+20)
	plus := image.Rect(bx+20, y-2, bx+36, y-2+20)
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
	r = image.Rect(rect.Min.X, y-2, rect.Min.X+rect.Dx(), y-2+20)
	if r.Overlaps(image.Rect(mx, my, mx+1, my+1)) {
		g.smartRender = !g.smartRender
		g.needsRedraw = true
		return true
	}
	y += OptionsMenuSpacing

	// FPS (not clickable)
	y += OptionsMenuSpacing

	// Version (not clickable)
	y += OptionsMenuSpacing

	// Close
	r = image.Rect(rect.Min.X, y-2, rect.Min.X+rect.Dx(), y-2+20)
	if r.Overlaps(image.Rect(mx, my, mx+1, my+1)) {
		g.showOptions = false
		g.needsRedraw = true
		return true
	}

	return true
}
