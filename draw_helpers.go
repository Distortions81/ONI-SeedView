package main

import (
	"fmt"
	"image"
	"image/color"
	"math"
	"os"
	"time"

	"golang.org/x/image/bmp"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

func (g *Game) drawUI(screen *ebiten.Image) {
	if !g.screenshotMode {
		cx := g.width / 2
		cy := g.height / 2
		crossClr := color.RGBA{255, 255, 255, 30}
		size := uiScaled(CrosshairSize)
		thickness := uiScaled(1)
		vector.StrokeLine(screen, float32(cx-size), float32(cy), float32(cx+size), float32(cy), float32(thickness), crossClr, true)
		vector.StrokeLine(screen, float32(cx), float32(cy-size), float32(cx), float32(cy+size), float32(thickness), crossClr, true)
		if g.showItemNames {
			worldX := int(math.Round(((float64(cx) - g.camX) / g.zoom) / 2))
			worldY := int(math.Round(((float64(cy) - g.camY) / g.zoom) / 2))
			coords := fmt.Sprintf("X: %d Y: %d", worldX, worldY)
			lh := notoFont.Metrics().Height.Ceil()
			y := g.height - lh - uiScaled(8)
			drawTextWithBGScale(screen, coords, uiScaled(5), y, 1, false)
		}
	}

	if g.showInfo {
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
		tx := g.width/2 - panelW/2
		ty := g.height - panelH - uiScaled(30)
		g.drawInfoPanel(screen, g.infoText, g.infoIcon, tx, ty)
	}
}

func (g *Game) drawOverlayMenus(screen *ebiten.Image) {
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
		rect := g.helpMenuRect()
		drawFrame(screen, rect)
		drawTextScale(screen, helpMessage, rect.Min.X+2, rect.Min.Y+2, 1, false)
		cr := g.helpCloseRect()
		drawCloseButton(screen, cr)
	}

	g.needsRedraw = false
	g.lastDraw = time.Now()
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
		_ = bmp.Encode(f, img)
		_ = f.Close()
	}
	g.captured = true
}
