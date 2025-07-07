package main

import (
	"image"
	"runtime"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

func (g *Game) Update() error {
	const panSpeed = PanSpeed

	g.checkRedrawTriggers()
	g.processScreenshot()

	oldX, oldY, oldZoom := g.camX, g.camY, g.zoom

	if g.handleGeyserListInput() {
		g.handleTouchGestures(oldX, oldY)
		return nil
	}

	if g.handleAsteroidMenuInput() {
		return nil
	}

	// Keyboard panning
	if ebiten.IsKeyPressed(ebiten.KeyLeft) || ebiten.IsKeyPressed(ebiten.KeyA) {
		g.camX += panSpeed
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) || ebiten.IsKeyPressed(ebiten.KeyD) {
		g.camX -= panSpeed
	}
	if ebiten.IsKeyPressed(ebiten.KeyUp) || ebiten.IsKeyPressed(ebiten.KeyW) {
		g.camY += panSpeed
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) || ebiten.IsKeyPressed(ebiten.KeyS) {
		g.camY -= panSpeed
	}

	mxTmp, myTmp := ebiten.CursorPosition()
	mousePressed := ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft)
	mouseJustPressed := inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft)
	if g.ssPending > 0 || g.skipClickTicks > 0 {
		mousePressed = false
		mouseJustPressed = false
	}

	justPressed := mouseJustPressed

	// Mouse dragging
	if mousePressed {
		if g.dragging {
			g.camX += float64(mxTmp - g.lastX)
			g.camY += float64(myTmp - g.lastY)
		}
		g.lastX = mxTmp
		g.lastY = myTmp
		g.dragging = true
	} else {
		g.dragging = false
	}

	g.handleTouchGestures(oldX, oldY)

	// Zoom with keyboard
	zoomFactor := 1.0
	if ebiten.IsKeyPressed(ebiten.KeyEqual) || ebiten.IsKeyPressed(ebiten.KeyKPAdd) {
		zoomFactor *= KeyZoomFactor
	}
	if ebiten.IsKeyPressed(ebiten.KeyMinus) || ebiten.IsKeyPressed(ebiten.KeyKPSubtract) {
		zoomFactor /= KeyZoomFactor
	}

	// Zoom with mouse wheel
	_, wheelY := ebiten.Wheel()
	if runtime.GOARCH == "wasm" && wheelY != 0 {
		now := time.Now()
		if now.Sub(g.lastWheel) < WheelThrottle {
			wheelY = 0
		} else {
			g.lastWheel = now
		}
	}
	if wheelY != 0 {
		if g.showGeyserList {
			g.adjustGeyserScroll(-float64(wheelY) * 10)
		} else {
			handled := false
			useNumbers := g.useNumbers && g.zoom < LegendZoomThreshold && !g.screenshotMode
			if g.legend != nil && g.showLegend {
				lw := g.legend.Bounds().Dx()
				lh := g.legend.Bounds().Dy()
				if lh > g.height && mxTmp >= 0 && mxTmp < lw {
					g.biomeScroll -= float64(wheelY) * 10
					if g.biomeScroll < 0 {
						g.biomeScroll = 0
					}
					max := g.maxBiomeScroll()
					if g.biomeScroll > max {
						g.biomeScroll = max
					}
					g.needsRedraw = true
					handled = true
				}
			}
			if !handled && useNumbers && g.legendImage != nil && g.showLegend {
				lw := g.legendImage.Bounds().Dx()
				lh := g.legendImage.Bounds().Dy()
				x0 := g.width - lw - 12
				if lh+10 > g.height && mxTmp >= x0 && mxTmp < x0+lw {
					g.itemScroll -= float64(wheelY) * 10
					if g.itemScroll < 0 {
						g.itemScroll = 0
					}
					max := g.maxItemScroll()
					if g.itemScroll > max {
						g.itemScroll = max
					}
					g.needsRedraw = true
					handled = true
				}
			}
			if !handled {
				if wheelY > 0 {
					zoomFactor *= WheelZoomFactor
				} else {
					zoomFactor /= WheelZoomFactor
				}
			}
		}
	}

	if zoomFactor != 1.0 {
		oldZoom := g.zoom
		g.zoom *= zoomFactor
		if g.zoom < g.minZoom {
			g.zoom = g.minZoom
		}
		if g.zoom > MaxZoom {
			g.zoom = MaxZoom
		}
		cx, cy := float64(g.width)/2, float64(g.height)/2
		worldX := (cx - g.camX) / oldZoom
		worldY := (cy - g.camY) / oldZoom
		g.camX = cx - worldX*g.zoom
		g.camY = cy - worldY*g.zoom
	}

	mx, my := -1, -1
	if !g.screenshotMode {
		mx, my = mxTmp, myTmp
		if mx < 0 || mx >= g.width || my < 0 || my >= g.height {
			mx, my = -1, -1
		} else if mx != g.lastMouseX || my != g.lastMouseY {
			g.touchUsed = false
		}
		g.lastMouseX, g.lastMouseY = mx, my
		if justPressed && g.helpRect().Overlaps(image.Rect(mx, my, mx+1, my+1)) {
			if g.showHelp {
				if time.Since(g.lastHelpClick) >= MenuToggleDelay {
					g.showHelp = false
				}
			} else {
				g.showHelp = true
			}
			g.lastHelpClick = time.Now()
			g.needsRedraw = true
		} else if g.showHelp && justPressed && g.helpCloseRect().Overlaps(image.Rect(mx, my, mx+1, my+1)) {
			g.showHelp = false
			g.needsRedraw = true
		} else if g.showShotMenu {
			if justPressed {
				if g.screenshotRect().Overlaps(image.Rect(mx, my, mx+1, my+1)) {
					if time.Since(g.lastShotClick) >= MenuToggleDelay {
						g.showShotMenu = false
						g.noColor = false
					}
					g.lastShotClick = time.Now()
					g.needsRedraw = true
				} else if !g.clickScreenshotMenu(mx, my) {
					if !g.screenshotMenuRect().Overlaps(image.Rect(mx, my, mx+1, my+1)) && !g.screenshotRect().Overlaps(image.Rect(mx, my, mx+1, my+1)) {
						g.showShotMenu = false
						g.noColor = false
						g.needsRedraw = true
					}
				}
			}
		} else if g.showOptions {
			if justPressed {
				if g.optionsRect().Overlaps(image.Rect(mx, my, mx+1, my+1)) {
					if time.Since(g.lastOptionsClick) >= MenuToggleDelay {
						g.showOptions = false
					}
					g.lastOptionsClick = time.Now()
					g.needsRedraw = true
				} else if !g.clickOptionsMenu(mx, my) {
					if !g.optionsMenuRect().Overlaps(image.Rect(mx, my, mx+1, my+1)) && !g.optionsRect().Overlaps(image.Rect(mx, my, mx+1, my+1)) {
						g.showOptions = false
						g.needsRedraw = true
					}
				}
			}
		} else if justPressed && g.screenshotRect().Overlaps(image.Rect(mx, my, mx+1, my+1)) {
			if g.showShotMenu {
				if time.Since(g.lastShotClick) >= MenuToggleDelay {
					g.showShotMenu = false
					g.noColor = false
				}
			} else {
				g.closeMenus()
				g.showShotMenu = true
				g.noColor = g.ssNoColor
			}
			g.lastShotClick = time.Now()
			g.needsRedraw = true
		} else if justPressed && g.optionsRect().Overlaps(image.Rect(mx, my, mx+1, my+1)) {
			if g.showOptions {
				if time.Since(g.lastOptionsClick) >= MenuToggleDelay {
					g.showOptions = false
				}
			} else {
				g.closeMenus()
				g.showOptions = true
			}
			g.lastOptionsClick = time.Now()
			g.needsRedraw = true
		} else if justPressed && g.asteroidInfoRect().Overlaps(image.Rect(mx, my, mx+1, my+1)) {
			if g.showAstMenu {
				if time.Since(g.lastAsteroidClick) >= MenuToggleDelay {
					g.showAstMenu = false
				}
			} else {
				g.closeMenus()
				g.showAstMenu = true
			}
			g.lastAsteroidClick = time.Now()
			g.needsRedraw = true
		} else if justPressed && g.clickLegend(mx, my) {
			// handled in clickLegend
		} else if justPressed && g.geyserRect().Overlaps(image.Rect(mx, my, mx+1, my+1)) {
			if g.showGeyserList {
				if time.Since(g.lastGeyserClick) >= MenuToggleDelay {
					g.showGeyserList = false
				}
			} else {
				g.closeMenus()
				g.camX = oldX
				g.camY = oldY
				g.dragging = false
				g.showGeyserList = true
			}
			g.lastGeyserClick = time.Now()
			g.needsRedraw = true
		} else if justPressed {
			if info, ix, iy, icon, found := g.itemAt(mx, my); found {
				g.camX += float64(g.width/2 - ix)
				g.camY += float64(g.height/2 - iy)
				g.clampCamera()

				if max := g.maxBiomeScroll(); g.biomeScroll > max {
					g.biomeScroll = max
				}
				if g.biomeScroll < 0 {
					g.biomeScroll = 0
				}
				if max := g.maxItemScroll(); g.itemScroll > max {
					g.itemScroll = max
				}
				if g.itemScroll < 0 {
					g.itemScroll = 0
				}

				g.infoText = info
				g.infoIcon = icon
				g.showInfo = true
				g.infoPinned = true
				g.needsRedraw = true
			}
		}
	}

	if !g.screenshotMode && !g.touchUsed {
		g.updateHover(mx, my)
		g.updateIconHover(mx, my)
	} else if g.hoverIcon != hoverNone {
		g.hoverIcon = hoverNone
		g.needsRedraw = true
	}

	if g.dragging {
		g.showInfo = false
		g.infoPinned = false
	}

	prevShow := g.showInfo
	prevPinned := g.infoPinned
	prevText := g.infoText
	prevIcon := g.infoIcon

	info, ix, iy, icon, found := "", 0, 0, (*ebiten.Image)(nil), false
	if g.mobile {
		cx := g.width / 2
		cy := g.height / 2
		info, ix, iy, icon, found = g.itemAt(cx, cy)
	} else if !g.touchUsed {
		info, ix, iy, icon, found = g.itemAt(mx, my)
	}
	if found {
		g.infoText = info
		g.infoX = ix
		g.infoY = iy
		g.infoIcon = icon
		g.showInfo = true
		if mousePressed {
			g.infoPinned = true
		}
	} else if !g.infoPinned || mousePressed {
		g.showInfo = false
		g.infoPinned = false
		g.infoIcon = nil
	}

	if g.showInfo != prevShow || g.infoPinned != prevPinned || g.infoText != prevText || g.infoIcon != prevIcon {
		g.needsRedraw = true
	}

	if !g.itemPanelVisible() && g.selectedItem != -1 {
		g.selectedItem = -1
		g.needsRedraw = true
	}

	g.clampCamera()

	if g.camX != oldX || g.camY != oldY || g.zoom != oldZoom {
		g.needsRedraw = true
	}

	if g.captured {
		return ebiten.Termination
	}

	return nil
}
