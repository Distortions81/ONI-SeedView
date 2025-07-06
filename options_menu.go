package main

import (
	"fmt"
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

func (g *Game) optionsRect() image.Rectangle {
	size := g.iconSize()
	x := g.width - size*4 - uiScaled(HelpMargin*4)
	y := g.height - size - uiScaled(HelpMargin)
	return image.Rect(x, y, x+size, y+size)
}

func (g *Game) optionsMenuSize() (int, int) {
	uiLabel := fmt.Sprintf("UI Scale [-] [+] %.0f%%", uiScale*100)
	labels := []string{
		OptionsMenuTitle,
		"Show Item Names",
		"Show Legends",
		"Use Item Numbers",
		"Icon Size [-] [+]",
		uiLabel,
		"Textures",
		"Vsync",
		"Power Saver",
		"Linear Filtering",
		"HiDPI Mode",
		"FPS: 60.0",
		"Version: " + ClientVersion,
		"Close",
	}
	maxW := 0
	for _, s := range labels {
		w, _ := textDimensions(s)
		if w > maxW {
			maxW = w
		}
	}
	w := maxW + uiScaled(4)
	h := (len(labels)+1)*menuSpacing() + uiScaled(4)
	return w, h
}

func (g *Game) optionsMenuRect() image.Rectangle {
	w, h := g.optionsMenuSize()
	x := g.optionsRect().Min.X - w - uiScaled(10)
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
	w, h := g.optionsMenuSize()
	img := ebiten.NewImage(w, h)
	drawFrame(img, image.Rect(0, 0, w, h))
	pad := uiScaled(6)
	drawText(img, OptionsMenuTitle, pad, pad, false)
	y := pad + menuSpacing()

	drawToggle := func(label string, enabled bool) {
		btn := image.Rect(uiScaled(4), y-uiScaled(4), w-uiScaled(4), y-uiScaled(4)+menuButtonHeight())
		drawButton(img, btn, enabled)
		lh := menuButtonHeight() - 5
		if notoFont != nil {
			lh = notoFont.Metrics().Height.Ceil()
		}
		drawText(img, label, btn.Min.X+pad, btn.Min.Y+(menuButtonHeight()-lh)/2, false)
		y += menuSpacing()
	}

	drawToggle("Show Item Names", g.showItemNames)
	drawToggle("Show Legends", g.showLegend)
	drawToggle("Use Item Numbers", g.useNumbers)

	label := "Icon Size"
	drawText(img, label, pad, y, false)
	tw, _ := textDimensions(label)
	bx := pad + tw + pad
	minus := image.Rect(bx, y-uiScaled(4), bx+uiScaled(20), y-uiScaled(4)+menuButtonHeight())
	plus := image.Rect(bx+uiScaled(24), y-uiScaled(4), bx+uiScaled(44), y-uiScaled(4)+menuButtonHeight())
	drawButton(img, minus, false)
	drawPlusMinus(img, minus, true)
	drawButton(img, plus, false)
	drawPlusMinus(img, plus, false)
	y += menuSpacing()

	// UI Scale buttons
	label = "UI Scale"
	drawText(img, label, pad, y, false)
	tw, _ = textDimensions(label)
	bx = pad + tw + pad
	minus = image.Rect(bx, y-uiScaled(4), bx+uiScaled(20), y-uiScaled(4)+menuButtonHeight())
	plus = image.Rect(bx+uiScaled(24), y-uiScaled(4), bx+uiScaled(44), y-uiScaled(4)+menuButtonHeight())
	drawButton(img, minus, false)
	drawPlusMinus(img, minus, true)
	drawButton(img, plus, false)
	drawPlusMinus(img, plus, false)
	scaleStr := fmt.Sprintf("%.0f%%", uiScale*100)
	drawText(img, scaleStr, plus.Max.X+pad, y, false)
	y += menuSpacing()

	drawToggle("Textures", g.textures)
	drawToggle("Vsync", g.vsync)
	drawToggle("Power Saver", g.smartRender)
	drawToggle("Linear Filtering", g.linearFilter)
	drawToggle("HiDPI Mode", g.hidpi)

	fps := fmt.Sprintf("FPS: %.1f", ebiten.ActualFPS())
	drawText(img, fps, pad, y, false)
	y += menuSpacing()

	drawText(img, "Version: "+ClientVersion, pad, y, false)
	y += menuSpacing()

	btn := image.Rect(uiScaled(4), y-uiScaled(4), w-uiScaled(4), y-uiScaled(4)+menuButtonHeight())
	drawButton(img, btn, true)
	lh := menuButtonHeight() - 5
	if notoFont != nil {
		lh = notoFont.Metrics().Height.Ceil()
	}
	drawText(img, "Close", btn.Min.X+pad, btn.Min.Y+(menuButtonHeight()-lh)/2, false)

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(rect.Min.X), float64(rect.Min.Y))
	dst.DrawImage(img, op)
}

func (g *Game) clickOptionsMenu(mx, my int) bool {
	rect := g.optionsMenuRect()
	if !rect.Overlaps(image.Rect(mx, my, mx+1, my+1)) {
		return false
	}
	x := mx - rect.Min.X
	y := my - rect.Min.Y
	mx = x
	my = y
	w, _ := g.optionsMenuSize()
	y = uiScaled(6) + menuSpacing()

	// Show Item Names
	r := image.Rect(uiScaled(4), y-uiScaled(4), w-uiScaled(4), y-uiScaled(4)+menuButtonHeight())
	if r.Overlaps(image.Rect(mx, my, mx+1, my+1)) {
		g.showItemNames = !g.showItemNames
		g.needsRedraw = true
		return true
	}
	y += menuSpacing()

	// Show Legends
	r = image.Rect(uiScaled(4), y-uiScaled(4), w-uiScaled(4), y-uiScaled(4)+menuButtonHeight())
	if r.Overlaps(image.Rect(mx, my, mx+1, my+1)) {
		g.showLegend = !g.showLegend
		g.needsRedraw = true
		return true
	}
	y += menuSpacing()

	// Use Item Numbers
	r = image.Rect(uiScaled(4), y-uiScaled(4), w-uiScaled(4), y-uiScaled(4)+menuButtonHeight())
	if r.Overlaps(image.Rect(mx, my, mx+1, my+1)) {
		g.useNumbers = !g.useNumbers
		g.needsRedraw = true
		return true
	}
	y += menuSpacing()

	// Icon Size buttons
	labelW, _ := textDimensions("Icon Size")
	bx := uiScaled(6) + labelW + uiScaled(6)
	minus := image.Rect(bx, y-uiScaled(4), bx+uiScaled(20), y-uiScaled(4)+menuButtonHeight())
	plus := image.Rect(bx+uiScaled(24), y-uiScaled(4), bx+uiScaled(44), y-uiScaled(4)+menuButtonHeight())
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
	y += menuSpacing()

	// Icon Size buttons
	labelW, _ = textDimensions("Icon Size")
	bx = uiScaled(6) + labelW + uiScaled(6)
	minus = image.Rect(bx, y-uiScaled(4), bx+uiScaled(20), y-uiScaled(4)+menuButtonHeight())
	plus = image.Rect(bx+uiScaled(24), y-uiScaled(4), bx+uiScaled(44), y-uiScaled(4)+menuButtonHeight())
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
	y += menuSpacing()

	// UI Scale buttons
	labelW, _ = textDimensions("UI Scale")
	bx = uiScaled(6) + labelW + uiScaled(6)
	minus = image.Rect(bx, y-uiScaled(4), bx+uiScaled(20), y-uiScaled(4)+menuButtonHeight())
	plus = image.Rect(bx+uiScaled(24), y-uiScaled(4), bx+uiScaled(44), y-uiScaled(4)+menuButtonHeight())
	if minus.Overlaps(image.Rect(mx, my, mx+1, my+1)) {
		if uiScale > 0.5 {
			setUIScale(uiScale - 0.1)
		}
		g.needsRedraw = true
		return true
	}
	if plus.Overlaps(image.Rect(mx, my, mx+1, my+1)) {
		setUIScale(uiScale + 0.1)
		g.needsRedraw = true
		return true
	}
	y += menuSpacing()

	// Textures
	r = image.Rect(uiScaled(4), y-uiScaled(4), w-uiScaled(4), y-uiScaled(4)+menuButtonHeight())
	if r.Overlaps(image.Rect(mx, my, mx+1, my+1)) {
		g.textures = !g.textures
		g.needsRedraw = true
		return true
	}
	y += menuSpacing()

	// Vsync
	r = image.Rect(uiScaled(4), y-uiScaled(4), w-uiScaled(4), y-uiScaled(4)+menuButtonHeight())
	if r.Overlaps(image.Rect(mx, my, mx+1, my+1)) {
		g.vsync = !g.vsync
		ebiten.SetVsyncEnabled(g.vsync)
		g.needsRedraw = true
		return true
	}
	y += menuSpacing()

	// Power Saver
	r = image.Rect(uiScaled(4), y-uiScaled(4), w-uiScaled(4), y-uiScaled(4)+menuButtonHeight())
	if r.Overlaps(image.Rect(mx, my, mx+1, my+1)) {
		g.smartRender = !g.smartRender
		g.needsRedraw = true
		return true
	}
	y += menuSpacing()

	// Linear Filtering
	r = image.Rect(uiScaled(4), y-uiScaled(4), w-uiScaled(4), y-uiScaled(4)+menuButtonHeight())
	if r.Overlaps(image.Rect(mx, my, mx+1, my+1)) {
		g.linearFilter = !g.linearFilter
		g.needsRedraw = true
		return true
	}
	y += menuSpacing()

	// HiDPI Mode
	r = image.Rect(uiScaled(4), y-uiScaled(4), w-uiScaled(4), y-uiScaled(4)+menuButtonHeight())
	if r.Overlaps(image.Rect(mx, my, mx+1, my+1)) {
		g.hidpi = !g.hidpi
		g.needsRedraw = true
		return true
	}
	y += menuSpacing()

	// FPS (not clickable)
	y += menuSpacing()

	// Version (not clickable)
	y += menuSpacing()

	// Close
	r = image.Rect(uiScaled(4), y-uiScaled(4), w-uiScaled(4), y-uiScaled(4)+menuButtonHeight())
	if r.Overlaps(image.Rect(mx, my, mx+1, my+1)) {
		g.showOptions = false
		g.needsRedraw = true
		return true
	}

	return true
}
