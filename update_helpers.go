package main

import (
	"image"
	"image/color"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

func (g *Game) checkRedrawTriggers() {
	if g.fitOnLoad {
		g.centerAndFit()
		g.needsRedraw = true
		g.fitOnLoad = false
	}
	if !g.smartRender {
		g.needsRedraw = true
	}
	minimized := ebiten.IsWindowMinimized()
	if g.wasMinimized && !minimized {
		g.needsRedraw = true
	}
	g.wasMinimized = minimized
	if time.Since(g.lastDraw) >= time.Second*3 {
		g.needsRedraw = true
	}
}

func (g *Game) processScreenshot() {
	if g.ssPending > 0 {
		if g.ssPending == 1 {
			g.saveScreenshot()
			g.ssSaved = time.Now()
			g.showShotMenu = false
			g.noColor = false
			g.skipClickTicks = 1
		}
		g.ssPending--
		g.needsRedraw = true
	} else if g.skipClickTicks > 0 {
		g.skipClickTicks--
	}
}

func (g *Game) handleGeyserListInput() bool {
	if !g.showGeyserList {
		return false
	}
	_, wheelY := ebiten.Wheel()
	if wheelY != 0 {
		g.adjustGeyserScroll(-float64(wheelY) * 10)
	}
	mx, my := ebiten.CursorPosition()
	if mx >= 0 && mx < g.width && my >= 0 && my < g.height {
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			if g.geyserCloseRect().Overlaps(image.Rect(mx, my, mx+1, my+1)) {
				g.showGeyserList = false
				g.needsRedraw = true
			}
		}
	}
	for _, id := range inpututil.AppendJustPressedTouchIDs(nil) {
		x, y := ebiten.TouchPosition(id)
		if x >= 0 && x < g.width && y >= 0 && y < g.height {
			if g.geyserCloseRect().Overlaps(image.Rect(x, y, x+1, y+1)) {
				g.showGeyserList = false
				g.needsRedraw = true
			}
		}
	}
	return true
}

func (g *Game) handleAsteroidMenuInput() bool {
	if !g.showAstMenu {
		return false
	}
	_, wheelY := ebiten.Wheel()
	if wheelY != 0 {
		g.adjustAsteroidScroll(-float64(wheelY) * 10)
	}
	mx, my := ebiten.CursorPosition()
	if mx >= 0 && mx < g.width && my >= 0 && my < g.height {
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			if !g.clickAsteroidMenu(mx, my) {
				if !g.asteroidMenuRect().Overlaps(image.Rect(mx, my, mx+1, my+1)) && !g.asteroidInfoRect().Overlaps(image.Rect(mx, my, mx+1, my+1)) {
					g.showAstMenu = false
					g.needsRedraw = true
				}
			}
		}
	}
	for _, id := range inpututil.AppendJustPressedTouchIDs(nil) {
		x, y := ebiten.TouchPosition(id)
		if x >= 0 && x < g.width && y >= 0 && y < g.height {
			if !g.clickAsteroidMenu(x, y) {
				if !g.asteroidMenuRect().Overlaps(image.Rect(x, y, x+1, y+1)) && !g.asteroidInfoRect().Overlaps(image.Rect(x, y, x+1, y+1)) {
					g.showAstMenu = false
					g.needsRedraw = true
				}
			}
		}
	}
	return true
}

func (g *Game) drawGeyserListScreen(dst *ebiten.Image) bool {
	if !g.showGeyserList {
		return false
	}
	dst.Fill(color.Black)
	g.drawGeyserList(dst)
	g.needsRedraw = false
	g.lastDraw = time.Now()
	return true
}

func (g *Game) drawLoadingScreen(dst *ebiten.Image) bool {
	if !(g.loading || (len(g.biomes) == 0 && g.status != "")) {
		return false
	}
	dst.Fill(backgroundColor)
	msg := g.status
	scale := 1.0
	if msg == "" {
		msg = "Fetching asteroid data..."
		scale = 2.0
	}
	_, h := textDimensions(msg)
	x := g.width / 2
	y := g.height/2 - int(float64(h)*scale/2)
	if g.statusError {
		drawTextWithBGBorderScale(dst, msg, x, y, errorBorderColor, scale, true)
	} else {
		drawTextWithBGScale(dst, msg, x, y, scale, true)
	}
	g.lastDraw = time.Now()
	return true
}
