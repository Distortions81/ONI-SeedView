//go:build !test

package main

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

func (g *Game) screenshotRect() image.Rectangle {
	size := HelpIconSize
	x := g.width - size*2 - HelpMargin*2
	y := g.height - size - HelpMargin
	return image.Rect(x, y, x+size, y+size)
}

func (g *Game) screenshotMenuRect() image.Rectangle {
	labels := []string{"Low", "Medium", "High", "Extreme", "Save Screenshot"}
	maxW := 0
	for _, s := range labels {
		if len(s) > maxW {
			maxW = len(s)
		}
	}
	w := maxW*LabelCharWidth + 4
	h := len(labels)*ScreenshotMenuSpacing + 4
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
	vector.DrawFilledRect(dst, float32(rect.Min.X), float32(rect.Min.Y), float32(rect.Dx()), float32(rect.Dy()), colorRGBA(0, 0, 0, 200), false)
	items := []string{"Low", "Medium", "High", "Extreme", "Save Screenshot"}
	y := rect.Min.Y + 2
	for i, it := range items {
		selected := i == g.ssQuality
		if selected {
			drawTextWithBGBorder(dst, it, rect.Min.X+2, y, colorRGBA(255, 255, 255, 255))
		} else {
			drawTextWithBG(dst, it, rect.Min.X+2, y)
		}
		y += ScreenshotMenuSpacing
	}
}

func (g *Game) clickScreenshotMenu(mx, my int) bool {
	rect := g.screenshotMenuRect()
	if !rect.Overlaps(image.Rect(mx, my, mx+1, my+1)) {
		return false
	}
	items := []string{"Low", "Medium", "High", "Extreme", "Save Screenshot"}
	y := rect.Min.Y + 2
	for i := range items {
		r := image.Rect(rect.Min.X, y-2, rect.Min.X+rect.Dx(), y-2+16+4)
		if r.Overlaps(image.Rect(mx, my, mx+1, my+1)) {
			switch i {
			case 0, 1, 2, 3:
				g.ssQuality = i
			case 4:
				g.saveScreenshot()
				g.showShotMenu = false
			}
			g.needsRedraw = true
			return true
		}
		y += ScreenshotMenuSpacing
	}
	return false
}

func (g *Game) saveScreenshot() {
	scale := 1.0
	switch g.ssQuality {
	case 0:
		scale = 1.0
	case 1:
		scale = 4.0
	case 2:
		scale = 8.0
	default:
		scale = 2.0
	}
	width := int(float64(g.astWidth) * 2 * scale)
	height := int(float64(g.astHeight) * 2 * scale)
	img := g.captureScreenshot(width, height, scale)
	var buf bytes.Buffer
	_ = png.Encode(&buf, img)
	name := fmt.Sprintf("%s-%s.png", g.coord, time.Now().Format("20060102-150405"))
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
	img := ebiten.NewImage(w, h)
	g.needsRedraw = true
	g.Draw(img)
	b := img.Bounds()
	pix := make([]byte, 4*b.Dx()*b.Dy())
	img.ReadPixels(pix)
	rgba := &image.RGBA{Pix: pix, Stride: 4 * b.Dx(), Rect: b}
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

func colorRGBA(r, g0, b, a uint8) color.RGBA {
	return color.RGBA{R: r, G: g0, B: b, A: a}
}
