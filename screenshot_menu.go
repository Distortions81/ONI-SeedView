package main

import (
	"bytes"
	"fmt"
	"image"
	"time"

	"golang.org/x/image/bmp"

	"github.com/hajimehoshi/ebiten/v2"
)

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

func (g *Game) drawScreenshotMenu(dst *ebiten.Image) {
	rect := g.screenshotMenuRect()
	w, h := g.screenshotMenuSize()
	img := ebiten.NewImage(w, h)
	drawFrame(img, image.Rect(0, 0, w, h))
	pad := uiScaled(6)
	drawText(img, ScreenshotMenuTitle, pad, pad, false)

	label := ScreenshotSaveLabel
	if g.ssPending > 0 {
		label = ScreenshotTakingLabel
	} else if time.Since(g.ssSaved) < 3*time.Second {
		label = ScreenshotSavedLabel
	}
	items := append([]string(nil), ScreenshotQualities...)
	items = append(items, ScreenshotBWLabel, label, ScreenshotCancelLabel)
	y := pad + menuSpacing()
	for i, it := range items {
		btn := image.Rect(uiScaled(4), y-uiScaled(4), w-uiScaled(4), y-uiScaled(4)+menuButtonHeight())
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
		drawText(img, it, btn.Min.X+pad, btn.Min.Y+(menuButtonHeight()-lh)/2, false)
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
	items = append(items, ScreenshotBWLabel, ScreenshotSaveLabel, ScreenshotCancelLabel)
	y = uiScaled(6) + menuSpacing()
	w, _ := g.screenshotMenuSize()
	for i := range items {
		r := image.Rect(uiScaled(4), y-uiScaled(4), w-uiScaled(4), y-uiScaled(4)+menuButtonHeight())
		if r.Overlaps(image.Rect(mx, my, mx+1, my+1)) {
			switch i {
			case 0, 1, 2:
				g.ssQuality = i
			case len(ScreenshotQualities):
				g.ssNoColor = !g.ssNoColor
				g.noColor = g.ssNoColor
			case len(ScreenshotQualities) + 1:
				if g.ssPending == 0 {
					g.ssPending = 2
				}
			case len(ScreenshotQualities) + 2:
				g.showShotMenu = false
				g.noColor = false
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

	var out image.Image = img
	if g.ssNoColor {
		out = rgbaToGray(img)
	}

	var buf bytes.Buffer
	_ = bmp.Encode(&buf, out)
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

func rgbaToGray(src *image.RGBA) *image.Gray {
	dst := image.NewGray(src.Bounds())
	for y := 0; y < src.Rect.Dy(); y++ {
		sp := src.Pix[y*src.Stride:]
		dp := dst.Pix[y*dst.Stride:]
		for x := 0; x < src.Rect.Dx(); x++ {
			i := x * 4
			r := sp[i]
			gr := sp[i+1]
			b := sp[i+2]
			yv := (299*uint16(r) + 587*uint16(gr) + 114*uint16(b) + 500) / 1000
			dp[x] = uint8(yv)
		}
	}
	return dst
}
