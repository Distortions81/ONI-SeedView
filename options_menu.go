//go:build !test

package main

import (
	"fmt"
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

func (g *Game) optionsRect() image.Rectangle {
	size := g.iconSize()
	m := g.iconMargin()
	x := g.width - size*4 - m*4
	y := g.height - size - m
	return image.Rect(x, y, x+size, y+size)
}

func (g *Game) optionsMenuSize() (int, int) {
	fontLabel := fmt.Sprintf("Font Size [-] [+] %.0fpt", fontSize)
	dpiLabel := fmt.Sprintf("DPI Scale [-] [+] %.0f%%", g.dpiAdjust*100)
	labels := []string{
		OptionsMenuTitle,
		"Show Item Names",
		"Show Legends",
		"Use Item Numbers",
		fontLabel,
		"Icon Size [-] [+]",
		"Textures",
		"Vsync",
		"Power Saver",
		"Linear Filtering",
		"High DPI",
		dpiLabel,
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
	w := maxW + 4
	h := (len(labels)+1)*menuSpacing() + 4
	return w, h
}

func (g *Game) optionsMenuRect() image.Rectangle {
	w, h := g.optionsMenuSize()
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
	w, h := g.optionsMenuSize()
	img := ebiten.NewImage(w, h)
	drawFrame(img, image.Rect(0, 0, w, h))
	drawText(img, OptionsMenuTitle, 6, 6, false)
	y := 6 + menuSpacing()

	drawToggle := func(label string, enabled bool) {
		btn := image.Rect(4, y-4, w-4, y-4+menuButtonHeight())
		drawButton(img, btn, enabled)
		lh := menuButtonHeight() - 5
		if notoFont != nil {
			lh = notoFont.Metrics().Height.Ceil()
		}
		drawText(img, label, btn.Min.X+6, btn.Min.Y+(menuButtonHeight()-lh)/2, false)
		y += menuSpacing()
	}

	drawToggle("Show Item Names", g.showItemNames)
	drawToggle("Show Legends", g.showLegend)
	drawToggle("Use Item Numbers", g.useNumbers)

	label := "Font Size"
	drawText(img, label, 6, y, false)
	tw, _ := textDimensions(label)
	bx := 6 + tw + 6
	minus := image.Rect(bx, y-4, bx+20, y-4+menuButtonHeight())
	plus := image.Rect(bx+24, y-4, bx+44, y-4+menuButtonHeight())
	drawButton(img, minus, false)
	drawPlusMinus(img, minus, true)
	drawButton(img, plus, false)
	drawPlusMinus(img, plus, false)
	sizeStr := fmt.Sprintf("%.0fpt", fontSize)
	drawText(img, sizeStr, plus.Max.X+6, y, false)
	y += menuSpacing()

	label = "Icon Size"
	drawText(img, label, 6, y, false)
	tw, _ = textDimensions(label)
	bx = 6 + tw + 6
	minus = image.Rect(bx, y-4, bx+20, y-4+menuButtonHeight())
	plus = image.Rect(bx+24, y-4, bx+44, y-4+menuButtonHeight())
	drawButton(img, minus, false)
	drawPlusMinus(img, minus, true)
	drawButton(img, plus, false)
	drawPlusMinus(img, plus, false)
	y += menuSpacing()

	drawToggle("Textures", g.textures)
	drawToggle("Vsync", g.vsync)
	drawToggle("Power Saver", g.smartRender)
	drawToggle("Linear Filtering", g.linearFilter)
	drawToggle("High DPI", g.highDPI)

	label = "DPI Scale"
	drawText(img, label, 6, y, false)
	tw, _ = textDimensions(label)
	bx = 6 + tw + 6
	minus = image.Rect(bx, y-4, bx+20, y-4+menuButtonHeight())
	plus = image.Rect(bx+24, y-4, bx+44, y-4+menuButtonHeight())
	drawButton(img, minus, false)
	drawPlusMinus(img, minus, true)
	drawButton(img, plus, false)
	drawPlusMinus(img, plus, false)
	scaleStr := fmt.Sprintf("%.0f%%", g.dpiAdjust*100)
	drawText(img, scaleStr, plus.Max.X+6, y, false)
	y += menuSpacing()

	fps := fmt.Sprintf("FPS: %.1f", ebiten.ActualFPS())
	drawText(img, fps, 6, y, false)
	y += menuSpacing()

	drawText(img, "Version: "+ClientVersion, 6, y, false)
	y += menuSpacing()

	btn := image.Rect(4, y-4, w-4, y-4+menuButtonHeight())
	drawButton(img, btn, true)
	lh := menuButtonHeight() - 5
	if notoFont != nil {
		lh = notoFont.Metrics().Height.Ceil()
	}
	drawText(img, "Close", btn.Min.X+6, btn.Min.Y+(menuButtonHeight()-lh)/2, false)

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
	y = 6 + menuSpacing()

	// Show Item Names
	r := image.Rect(4, y-4, w-4, y-4+menuButtonHeight())
	if r.Overlaps(image.Rect(mx, my, mx+1, my+1)) {
		g.showItemNames = !g.showItemNames
		g.needsRedraw = true
		return true
	}
	y += menuSpacing()

	// Show Legends
	r = image.Rect(4, y-4, w-4, y-4+menuButtonHeight())
	if r.Overlaps(image.Rect(mx, my, mx+1, my+1)) {
		g.showLegend = !g.showLegend
		g.needsRedraw = true
		return true
	}
	y += menuSpacing()

	// Use Item Numbers
	r = image.Rect(4, y-4, w-4, y-4+menuButtonHeight())
	if r.Overlaps(image.Rect(mx, my, mx+1, my+1)) {
		g.useNumbers = !g.useNumbers
		g.needsRedraw = true
		return true
	}
	y += menuSpacing()

	// Font Size buttons
	labelW, _ := textDimensions("Font Size")
	bx := 6 + labelW + 6
	minus := image.Rect(bx, y-4, bx+20, y-4+menuButtonHeight())
	plus := image.Rect(bx+24, y-4, bx+44, y-4+menuButtonHeight())
	if minus.Overlaps(image.Rect(mx, my, mx+1, my+1)) {
		decreaseFontSize()
		if max := g.maxBiomeScroll(); max == 0 {
			g.biomeScroll = 0
		} else if g.biomeScroll > max {
			g.biomeScroll = max
		}
		if max := g.maxItemScroll(); max == 0 {
			g.itemScroll = 0
		} else if g.itemScroll > max {
			g.itemScroll = max
		}
		g.needsRedraw = true
		return true
	}
	if plus.Overlaps(image.Rect(mx, my, mx+1, my+1)) {
		increaseFontSize()
		if max := g.maxBiomeScroll(); max == 0 {
			g.biomeScroll = 0
		} else if g.biomeScroll > max {
			g.biomeScroll = max
		}
		if max := g.maxItemScroll(); max == 0 {
			g.itemScroll = 0
		} else if g.itemScroll > max {
			g.itemScroll = max
		}
		g.needsRedraw = true
		return true
	}
	y += menuSpacing()

	// Icon Size buttons
	labelW, _ = textDimensions("Icon Size")
	bx = 6 + labelW + 6
	minus = image.Rect(bx, y-4, bx+20, y-4+menuButtonHeight())
	plus = image.Rect(bx+24, y-4, bx+44, y-4+menuButtonHeight())
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

	// Textures
	r = image.Rect(4, y-4, w-4, y-4+menuButtonHeight())
	if r.Overlaps(image.Rect(mx, my, mx+1, my+1)) {
		g.textures = !g.textures
		g.needsRedraw = true
		return true
	}
	y += menuSpacing()

	// Vsync
	r = image.Rect(4, y-4, w-4, y-4+menuButtonHeight())
	if r.Overlaps(image.Rect(mx, my, mx+1, my+1)) {
		g.vsync = !g.vsync
		ebiten.SetVsyncEnabled(g.vsync)
		g.needsRedraw = true
		return true
	}
	y += menuSpacing()

	// Power Saver
	r = image.Rect(4, y-4, w-4, y-4+menuButtonHeight())
	if r.Overlaps(image.Rect(mx, my, mx+1, my+1)) {
		g.smartRender = !g.smartRender
		g.needsRedraw = true
		return true
	}
	y += menuSpacing()

	// Linear Filtering
	r = image.Rect(4, y-4, w-4, y-4+menuButtonHeight())
	if r.Overlaps(image.Rect(mx, my, mx+1, my+1)) {
		g.linearFilter = !g.linearFilter
		g.needsRedraw = true
		return true
	}
	y += menuSpacing()

	// High DPI
	r = image.Rect(4, y-4, w-4, y-4+menuButtonHeight())
	if r.Overlaps(image.Rect(mx, my, mx+1, my+1)) {
		g.highDPI = !g.highDPI
		g.updateDPIScale()
		if max := g.maxBiomeScroll(); max == 0 {
			g.biomeScroll = 0
		} else if g.biomeScroll > max {
			g.biomeScroll = max
		}
		if max := g.maxItemScroll(); max == 0 {
			g.itemScroll = 0
		} else if g.itemScroll > max {
			g.itemScroll = max
		}
		g.needsRedraw = true
		return true
	}
	y += menuSpacing()

	// DPI Scale buttons
	labelW, _ = textDimensions("DPI Scale")
	bx = 6 + labelW + 6
	minus = image.Rect(bx, y-4, bx+20, y-4+menuButtonHeight())
	plus = image.Rect(bx+24, y-4, bx+44, y-4+menuButtonHeight())
	if minus.Overlaps(image.Rect(mx, my, mx+1, my+1)) {
		if g.dpiAdjust > 0.5 {
			g.dpiAdjust -= 0.25
			g.updateDPIScale()
		}
		g.needsRedraw = true
		return true
	}
	if plus.Overlaps(image.Rect(mx, my, mx+1, my+1)) {
		g.dpiAdjust += 0.25
		g.updateDPIScale()
		g.needsRedraw = true
		return true
	}
	y += menuSpacing()

	// FPS (not clickable)
	y += menuSpacing()

	// Version (not clickable)
	y += menuSpacing()

	// Close
	r = image.Rect(4, y-4, w-4, y-4+menuButtonHeight())
	if r.Overlaps(image.Rect(mx, my, mx+1, my+1)) {
		g.showOptions = false
		g.needsRedraw = true
		return true
	}

	return true
}
