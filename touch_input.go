//go:build !test

package main

import (
	"image"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

func (g *Game) handleTouchGestures(oldX, oldY float64) {
	// Touch gestures
	touchIDs := ebiten.TouchIDs()
	justPressedIDs := inpututil.AppendJustPressedTouchIDs(nil)
	if g.ssPending > 0 || g.skipClickTicks > 0 {
		touchIDs = nil
		justPressedIDs = nil
	}
	if len(touchIDs) > 0 {
		g.touchUsed = true
	}
	for _, id := range justPressedIDs {
		x, y := ebiten.TouchPosition(id)
		g.touchStartX = x
		g.touchStartY = y
		g.touchMoved = false
		g.touchActive = true
		g.touchUI = false
		if g.showGeyserList || g.showShotMenu || g.showAstMenu {
			g.touchUI = true
		} else {
			pt := image.Rect(x, y, x+1, y+1)
			if g.helpRect().Overlaps(pt) ||
				g.geyserRect().Overlaps(pt) || g.optionsRect().Overlaps(pt) {
				g.touchUI = true
			} else {
				if g.legend != nil && g.showLegend {
					if g.biomeLegendRect().Overlaps(pt) {
						g.touchUI = true
					}
				}
				useNumbers := g.useNumbers && g.zoom < LegendZoomThreshold && !g.screenshotMode
				if !g.touchUI && useNumbers && g.legendImage != nil && g.showLegend {
					if g.itemLegendRect().Overlaps(pt) {
						g.touchUI = true
					}
				}
			}
		}
		if g.touchUI {
			g.updateHover(x, y)
		}
	}
	switch len(touchIDs) {
	case 1:
		id := touchIDs[0]
		x, y := ebiten.TouchPosition(id)
		if last, ok := g.touches[id]; ok {
			dx := x - last.x
			dy := y - last.y
			if g.touchUI {
				if g.showGeyserList {
					g.adjustGeyserScroll(-float64(dy))
				} else {
					if g.legend != nil && g.showLegend {
						pt := image.Rect(g.touchStartX, g.touchStartY, g.touchStartX+1, g.touchStartY+1)
						if g.biomeLegendRect().Overlaps(pt) {
							g.biomeScroll -= float64(dy)
							if g.biomeScroll < 0 {
								g.biomeScroll = 0
							}
							if max := g.maxBiomeScroll(); g.biomeScroll > max {
								g.biomeScroll = max
							}
							g.updateHover(x, y)
						}
					}
					useNumbers := g.useNumbers && g.zoom < LegendZoomThreshold && !g.screenshotMode
					if useNumbers && g.legendImage != nil && g.showLegend {
						pt := image.Rect(g.touchStartX, g.touchStartY, g.touchStartX+1, g.touchStartY+1)
						if g.itemLegendRect().Overlaps(pt) {
							g.itemScroll -= float64(dy)
							if g.itemScroll < 0 {
								g.itemScroll = 0
							}
							if max := g.maxItemScroll(); g.itemScroll > max {
								g.itemScroll = max
							}
							g.updateHover(x, y)
						}
					}
					g.needsRedraw = true
				}
			} else {
				g.camX += float64(dx)
				g.camY += float64(dy)
				if abs(dx) > TouchDragThreshold || abs(dy) > TouchDragThreshold {
					g.touchMoved = true
				}
			}
		} else {
			g.touchStartX = x
			g.touchStartY = y
			g.touchMoved = false
			g.touchActive = true
			g.touchUI = false
			if g.showGeyserList || g.showShotMenu || g.showAstMenu {
				g.touchUI = true
			} else {
				pt := image.Rect(x, y, x+1, y+1)
				if g.helpRect().Overlaps(pt) || g.screenshotRect().Overlaps(pt) ||
					g.geyserRect().Overlaps(pt) || g.optionsRect().Overlaps(pt) {
					g.touchUI = true
				} else {
					if g.legend != nil && g.showLegend {
						pt := image.Rect(x, y, x+1, y+1)
						if g.biomeLegendRect().Overlaps(pt) {
							g.touchUI = true
						}
					}
					useNumbers := g.useNumbers && g.zoom < LegendZoomThreshold && !g.screenshotMode
					if !g.touchUI && useNumbers && g.legendImage != nil && g.showLegend {
						pt := image.Rect(x, y, x+1, y+1)
						if g.itemLegendRect().Overlaps(pt) {
							g.touchUI = true
						}
					}
				}
			}
			if g.touchUI {
				g.updateHover(x, y)
			}
		}
		g.touches = map[ebiten.TouchID]touchPoint{id: {x: x, y: y}}
		g.pinchDist = 0
	case 2:
		id1, id2 := touchIDs[0], touchIDs[1]
		x1, y1 := ebiten.TouchPosition(id1)
		x2, y2 := ebiten.TouchPosition(id2)
		dx := float64(x2 - x1)
		dy := float64(y2 - y1)
		dist := math.Hypot(dx, dy)
		midX := float64(x1+x2) / 2
		midY := float64(y1+y2) / 2
		if g.pinchDist != 0 {
			factor := dist / g.pinchDist
			oldZoom := g.zoom
			g.zoom *= factor
			if g.zoom < g.minZoom {
				g.zoom = g.minZoom
			}
			if g.zoom > MaxZoom {
				g.zoom = MaxZoom
			}
			worldX := (midX - g.camX) / oldZoom
			worldY := (midY - g.camY) / oldZoom
			g.camX = midX - worldX*g.zoom
			g.camY = midY - worldY*g.zoom
		}
		g.pinchDist = dist
		g.touches = map[ebiten.TouchID]touchPoint{
			id1: {x: x1, y: y1},
			id2: {x: x2, y: y2},
		}
		g.touchActive = false
		g.touchMoved = false
	default:
		if g.touchActive && !g.touchMoved {
			mx, my := g.touchStartX, g.touchStartY
			pt := image.Rect(mx, my, mx+1, my+1)
			if g.showGeyserList {
				if g.geyserCloseRect().Overlaps(pt) {
					g.showGeyserList = false
					g.needsRedraw = true
				}
			} else if g.showShotMenu {
				if !g.clickScreenshotMenu(mx, my) {
					if !g.screenshotMenuRect().Overlaps(pt) && !g.screenshotRect().Overlaps(pt) {
						g.showShotMenu = false
						g.needsRedraw = true
					}
				}
			} else if g.showOptions {
				if !g.clickOptionsMenu(mx, my) {
					if !g.optionsMenuRect().Overlaps(pt) && !g.optionsRect().Overlaps(pt) {
						g.showOptions = false
						g.needsRedraw = true
					}
				}
			} else if g.screenshotRect().Overlaps(pt) {
				g.showShotMenu = true
				g.needsRedraw = true
			} else if g.optionsRect().Overlaps(pt) {
				g.showOptions = true
				g.needsRedraw = true
			} else if g.asteroidInfoRect().Overlaps(pt) {
				g.showAstMenu = true
				g.needsRedraw = true
			} else if g.helpRect().Overlaps(pt) {
				g.showHelp = !g.showHelp
				g.needsRedraw = true
			} else if g.showHelp && g.helpCloseRect().Overlaps(pt) {
				g.showHelp = false
				g.needsRedraw = true
			} else if g.geyserRect().Overlaps(pt) {
				g.camX = oldX
				g.camY = oldY
				g.dragging = false
				g.showGeyserList = true
				g.needsRedraw = true
			} else if g.touchUI {
				g.updateHover(mx, my)
				g.clickLegend(mx, my)
			} else {
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
		g.touchUI = false
		g.touches = nil
		g.pinchDist = 0
		g.touchActive = false
		g.touchMoved = false
	}

	if len(touchIDs) > 0 {
		g.showInfo = false
		g.infoPinned = false
	}
}
