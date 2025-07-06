//go:build !test

package main

import (
	"bytes"
	"fmt"
	"golang.org/x/image/bmp"
	"image"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

func (g *Game) screenshotRect() image.Rectangle {
	size := g.iconSize()
	x := g.width - size*2 - HelpMargin*2
	y := g.height - size - HelpMargin
	return image.Rect(x, y, x+size, y+size)
}

func (g *Game) screenshotMenuSize() (int, int) {
	labels := append([]string{ScreenshotMenuTitle}, ScreenshotQualities...)
	labels = append(labels, ScreenshotBWLabel, ScreenshotSaveLabel, ScreenshotCloseLabel)
	itemCount := len(labels)
	allLabels := append([]string(nil), labels...)
	allLabels = append(allLabels, ScreenshotTakingLabel, ScreenshotSavedLabel)
	maxW := 0
	for _, s := range allLabels {
		w, _ := textDimensions(s)
		if w > maxW {
			maxW = w
		}
	}
	w := maxW + 4
	h := (itemCount+1)*menuSpacing() + 4
	return w, h
}

func (g *Game) screenshotMenuRect() image.Rectangle {
	w, h := g.screenshotMenuSize()
	x := g.screenshotRect().Min.X - w - 10
	if x < 0 {
		x = 0
	}
	y := g.screenshotRect().Min.Y - h
	if y < 0 {
		y = 0
	}
	return image.Rect(x, y, x+w, y+h)
}

func (g *Game) drawScreenshotMenu(dst *ebiten.Image) {
	rect := g.screenshotMenuRect()
	w, h := g.screenshotMenuSize()
	img := ebiten.NewImage(w, h)
	drawFrame(img, image.Rect(0, 0, w, h))
	drawText(img, ScreenshotMenuTitle, 6, 6, false)

	label := ScreenshotSaveLabel
	if g.ssPending > 0 {
		label = ScreenshotTakingLabel
	} else if time.Since(g.ssSaved) < 2*time.Second {
		label = ScreenshotSavedLabel
	}
	items := append([]string(nil), ScreenshotQualities...)
	items = append(items, ScreenshotBWLabel, label, ScreenshotCloseLabel)
	y := 6 + menuSpacing()
	for i, it := range items {
		btn := image.Rect(4, y-4, w-4, y-4+menuButtonHeight())
		selected := i == g.ssQuality
		switch i {
		case len(ScreenshotQualities):
			drawButton(img, btn, g.ssNoColor)
		case len(ScreenshotQualities) + 1:
			if g.ssPending > 0 {
				drawButton(img, btn, true)
			} else {
				drawButton(img, btn, false)
			}
		case len(ScreenshotQualities) + 2:
			drawButton(img, btn, true)
		default:
			drawButton(img, btn, selected)
		}
		lh := menuButtonHeight() - 5
		if notoFont != nil {
			lh = notoFont.Metrics().Height.Ceil()
		}
		drawText(img, it, btn.Min.X+6, btn.Min.Y+(menuButtonHeight()-lh)/2, false)
		y += menuSpacing()
		if i == len(ScreenshotQualities)-1 || i == len(ScreenshotQualities) {
			y += menuSpacing()
		}
	}
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(rect.Min.X), float64(rect.Min.Y))
	dst.DrawImage(img, op)
}

func (g *Game) clickScreenshotMenu(mx, my int) bool {
	rect := g.screenshotMenuRect()
	if !rect.Overlaps(image.Rect(mx, my, mx+1, my+1)) {
		return false
	}
	x := mx - rect.Min.X
	y := my - rect.Min.Y
	mx = x
	my = y
	items := append([]string(nil), ScreenshotQualities...)
	items = append(items, ScreenshotBWLabel, ScreenshotSaveLabel, ScreenshotCloseLabel)
	y = 6 + menuSpacing()
	w, _ := g.screenshotMenuSize()
	for i := range items {
		r := image.Rect(4, y-4, w-4, y-4+menuButtonHeight())
		if r.Overlaps(image.Rect(mx, my, mx+1, my+1)) {
			switch i {
			case 0, 1, 2:
				g.ssQuality = i
			case len(ScreenshotQualities):
				g.ssNoColor = !g.ssNoColor
			case len(ScreenshotQualities) + 1:
				if g.ssPending == 0 {
					g.ssPending = 2
				}
			case len(ScreenshotQualities) + 2:
				g.showShotMenu = false
			}
			g.needsRedraw = true
			return true
		}
		y += menuSpacing()
		if i == len(ScreenshotQualities)-1 || i == len(ScreenshotQualities) {
			y += menuSpacing()
		}
	}
	return false
}

func (g *Game) saveScreenshot() {
	scale := ScreenshotScales[g.ssQuality]
	width := int(float64(g.astWidth) * 2 * scale)
	height := int(float64(g.astHeight) * 2 * scale)
	oldBW := g.noColor
	g.noColor = g.ssNoColor
	img := g.captureScreenshot(width, height, scale)
	g.noColor = oldBW
	var buf bytes.Buffer
	_ = bmp.Encode(&buf, img)
	name := fmt.Sprintf("%s-%s.bmp", g.coord, time.Now().Format("20060102-150405"))
	_ = saveImageData(name, buf.Bytes())
}

func (g *Game) captureScreenshot(w, h int, zoom float64) *image.RGBA {
	ow, oh := g.width, g.height
	ox, oy, oz := g.camX, g.camY, g.zoom
	showHelp := g.showHelp
	showInfo := g.showInfo
	menu := g.showShotMenu
	pinned := g.infoPinned
	g.showHelp = false
	g.showInfo = false
	g.infoPinned = false
	g.showShotMenu = false
	g.width = w
	g.height = h
	g.zoom = zoom
	g.camX = 0
	g.camY = 0
	g.screenshotMode = true
	oldSize := fontSize
	setFontSize(screenshotFontSize)
	img := ebiten.NewImage(w, h)
	g.needsRedraw = true
	g.Draw(img)
	b := img.Bounds()
	pix := make([]byte, 4*b.Dx()*b.Dy())
	img.ReadPixels(pix)
	rgba := &image.RGBA{Pix: pix, Stride: 4 * b.Dx(), Rect: b}
	setFontSize(oldSize)
	g.screenshotMode = false
	g.width = ow
	g.height = oh
	g.zoom = oz
	g.camX = ox
	g.camY = oy
	g.showHelp = showHelp
	g.showInfo = showInfo
	g.infoPinned = pinned
	g.showShotMenu = menu
	g.needsRedraw = true
	return rgba
}
