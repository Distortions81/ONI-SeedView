//go:build !test

package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math"
	"os"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

func (g *Game) drawUI(screen *ebiten.Image) {
	if !g.screenshotMode {
		cx := g.width / 2
		cy := g.height / 2
		vector.StrokeLine(screen, float32(cx-CrosshairSize), float32(cy), float32(cx+CrosshairSize), float32(cy), 1, crosshairColor, true)
		vector.StrokeLine(screen, float32(cx), float32(cy-CrosshairSize), float32(cx), float32(cy+CrosshairSize), 1, crosshairColor, true)
		worldX := int(math.Round(((float64(cx) - g.camX) / g.zoom) / 2))
		worldY := int(math.Round(((float64(cy) - g.camY) / g.zoom) / 2))
		coords := fmt.Sprintf("X: %d Y: %d", worldX, worldY)
		scale := g.uiScale()
		drawTextWithBGScale(screen, coords, 5, g.height-int(20*scale), scale)
	}

	if g.showInfo {
		scale := g.uiScale()
		w, h := textDimensions(g.infoText)
		iconW, iconH := 0, 0
		if g.infoIcon != nil {
			iconW = InfoIconSize
			iconH = InfoIconSize
		}
		panelW := w + iconW + 4
		panelH := h
		if iconH > h {
			panelH = iconH
		}
		tx := g.width/2 - int(float64(panelW)*scale)/2
		ty := g.height - int(float64(panelH)*scale) - 30
		g.drawInfoPanel(screen, g.infoText, g.infoIcon, tx, ty, scale)
	}

	if g.showShotMenu {
		g.drawScreenshotMenu(screen)
	}
	if g.showOptions {
		g.drawOptionsMenu(screen)
	}
	if g.showAstMenu {
		g.drawAsteroidMenu(screen)
	}
	if g.showHelp && !g.screenshotMode {
		scale := g.uiScale()
		rect := g.helpMenuRect()
		drawFrame(screen, rect)
		drawTextScale(screen, helpMessage, rect.Min.X+2, rect.Min.Y+2, scale)
		cr := g.helpCloseRect()
		drawCloseButton(screen, cr)
	}

	g.needsRedraw = false
	g.lastDraw = time.Now()
	if g.noColor {
		if g.grayImage == nil || g.grayImage.Bounds().Dx() != g.width || g.grayImage.Bounds().Dy() != g.height {
			g.grayImage = ebiten.NewImage(g.width, g.height)
		}
		g.grayImage.Clear()
		g.grayImage.DrawImage(screen, nil)
		screen.Clear()
		op := &ebiten.DrawImageOptions{}
		op.ColorM = grayColorM
		screen.DrawImage(g.grayImage, op)
	}
}

func (g *Game) captureScreen(screen *ebiten.Image) {
	if g.screenshotPath == "" || g.captured {
		return
	}
	b := screen.Bounds()
	pixels := make([]byte, 4*b.Dx()*b.Dy())
	screen.ReadPixels(pixels)
	img := &image.RGBA{Pix: pixels, Stride: 4 * b.Dx(), Rect: b}
	if f, err := os.Create(g.screenshotPath); err == nil {
		_ = png.Encode(f, img)
		f.Close()
	}
	g.captured = true
}
