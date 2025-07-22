package main

import (
	"bytes"
	"fmt"
	"image"
	"time"

	"golang.org/x/image/bmp"

	"github.com/hajimehoshi/ebiten/v2"
)

func desaturateImage(img *image.RGBA) {
	const r = 0.299
	const g = 0.587
	const b = 0.114
	pix := img.Pix
	for i := 0; i < len(pix); i += 4 {
		gray := uint8(float64(pix[i])*r + float64(pix[i+1])*g + float64(pix[i+2])*b)
		pix[i] = gray
		pix[i+1] = gray
		pix[i+2] = gray
	}
}

func (g *Game) screenshotRect() image.Rectangle {
	size := g.iconSize()
	x := g.width - size*2 - uiScaled(HelpMargin*2)
	y := g.height - size - uiScaled(HelpMargin)
	return image.Rect(x, y, x+size, y+size)
}

func (g *Game) screenshotMenuSize() (int, int) {
	labels := append([]string{ScreenshotMenuTitle}, ScreenshotQualities...)
	labels = append(labels, ScreenshotBWLabel, ScreenshotSaveLabel, ScreenshotCancelLabel)
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
	w := maxW + uiScaled(4)
	// Extra spacing after quality options and the B&W toggle
	h := (itemCount+2)*menuSpacing() + uiScaled(6)
	return w, h
}

func (g *Game) screenshotMenuRect() image.Rectangle {
	w, h := g.screenshotMenuSize()
	pad := uiScaled(2)
	x := g.screenshotRect().Min.X - w - uiScaled(10)
	if x < pad {
		x = pad
	}
	if x+w > g.width-pad {
		x = g.width - w - pad
	}
	y := g.screenshotRect().Min.Y - h
	if y < pad {
		y = pad
	}
	if y+h > g.height-pad {
		y = g.height - h - pad
	}
	return image.Rect(x, y, x+w, y+h)
}

func (g *Game) saveScreenshot() {
	scale := ScreenshotScales[g.ssQuality]
	width := int(float64(g.astWidth) * 2 * scale)
	height := int(float64(g.astHeight) * 2 * scale)
	img := g.captureScreenshot(width, height, scale)
	if g.ssNoColor {
		desaturateImage(img)
	}
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
	oldScale := uiScale
	setUIScale(4.0)
	img := ebiten.NewImage(w, h)
	g.needsRedraw = true
	g.Draw(img)
	b := img.Bounds()
	pix := make([]byte, 4*b.Dx()*b.Dy())
	img.ReadPixels(pix)
	rgba := &image.RGBA{Pix: pix, Stride: 4 * b.Dx(), Rect: b}
	setUIScale(oldScale)
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
