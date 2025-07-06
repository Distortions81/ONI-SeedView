//go:build !test

package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math"
	"os"
	"runtime"
	"strconv"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

func (g *Game) Update() error {
	const panSpeed = PanSpeed

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
	if time.Since(g.lastDraw) >= time.Second {
		g.needsRedraw = true
	}

iconsLoop:
	for g.iconResults != nil {
		select {
		case li, ok := <-g.iconResults:
			if !ok {
				g.iconResults = nil
				continue
			}
			if g.icons != nil {
				g.icons[li.name] = li.img
			}
			g.needsRedraw = true
		default:
			break iconsLoop
		}
	}

	if g.ssPending > 0 {
		if g.ssPending == 1 {
			g.saveScreenshot()
			g.ssSaved = time.Now()
			g.showShotMenu = false
			g.skipClickTicks = 1
		}
		g.ssPending--
		g.needsRedraw = true
	} else if g.skipClickTicks > 0 {
		g.skipClickTicks--
	}

	if g.showGeyserList {
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
		return nil
	}

	if g.showAstMenu {
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
		return nil
	}

	oldX, oldY, oldZoom := g.camX, g.camY, g.zoom

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
				scale := g.uiScale()
				if g.legend != nil {
					lw := int(float64(g.legend.Bounds().Dx()) * scale)
					if x >= 0 && x < lw {
						g.touchUI = true
					}
				}
				useNumbers := g.useNumbers && g.showItemNames && g.zoom < LegendZoomThreshold && !g.screenshotMode
				if !g.touchUI && useNumbers && g.legendImage != nil {
					lw := int(float64(g.legendImage.Bounds().Dx()) * scale)
					x0 := g.width - lw - 12
					if x >= x0 && x < x0+lw {
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
					scale := g.uiScale()
					if g.legend != nil {
						lw := int(float64(g.legend.Bounds().Dx()) * scale)
						if g.touchStartX >= 0 && g.touchStartX < lw {
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
					useNumbers := g.useNumbers && g.showItemNames && g.zoom < LegendZoomThreshold && !g.screenshotMode
					if useNumbers && g.legendImage != nil {
						lw := int(float64(g.legendImage.Bounds().Dx()) * scale)
						x0 := g.width - lw - 12
						if g.touchStartX >= x0 && g.touchStartX < x0+lw {
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
					scale := g.uiScale()
					if g.legend != nil {
						lw := int(float64(g.legend.Bounds().Dx()) * scale)
						if x >= 0 && x < lw {
							g.touchUI = true
						}
					}
					useNumbers := g.useNumbers && g.showItemNames && g.zoom < LegendZoomThreshold && !g.screenshotMode
					if !g.touchUI && useNumbers && g.legendImage != nil {
						lw := int(float64(g.legendImage.Bounds().Dx()) * scale)
						x0 := g.width - lw - 12
						if x >= x0 && x < x0+lw {
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
			scale := g.uiScale()
			useNumbers := g.useNumbers && g.showItemNames && g.zoom < LegendZoomThreshold && !g.screenshotMode
			if g.legend != nil {
				lw := int(float64(g.legend.Bounds().Dx()) * scale)
				lh := int(float64(g.legend.Bounds().Dy()) * scale)
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
			if !handled && useNumbers && g.legendImage != nil {
				lw := int(float64(g.legendImage.Bounds().Dx()) * scale)
				lh := int(float64(g.legendImage.Bounds().Dy()) * scale)
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
			g.showHelp = !g.showHelp
			g.needsRedraw = true
		} else if g.showHelp && justPressed && g.helpCloseRect().Overlaps(image.Rect(mx, my, mx+1, my+1)) {
			g.showHelp = false
			g.needsRedraw = true
		} else if g.showShotMenu {
			if justPressed {
				if !g.clickScreenshotMenu(mx, my) {
					if !g.screenshotMenuRect().Overlaps(image.Rect(mx, my, mx+1, my+1)) && !g.screenshotRect().Overlaps(image.Rect(mx, my, mx+1, my+1)) {
						g.showShotMenu = false
						g.needsRedraw = true
					}
				}
			}
		} else if g.showOptions {
			if justPressed {
				if !g.clickOptionsMenu(mx, my) {
					if !g.optionsMenuRect().Overlaps(image.Rect(mx, my, mx+1, my+1)) && !g.optionsRect().Overlaps(image.Rect(mx, my, mx+1, my+1)) {
						g.showOptions = false
						g.needsRedraw = true
					}
				}
			}
		} else if justPressed && g.screenshotRect().Overlaps(image.Rect(mx, my, mx+1, my+1)) {
			g.showShotMenu = true
			g.needsRedraw = true
		} else if justPressed && g.optionsRect().Overlaps(image.Rect(mx, my, mx+1, my+1)) {
			g.showOptions = true
			g.needsRedraw = true
		} else if justPressed && g.asteroidInfoRect().Overlaps(image.Rect(mx, my, mx+1, my+1)) {
			g.showAstMenu = true
			g.needsRedraw = true
		} else if justPressed && g.clickLegend(mx, my) {
			// handled in clickLegend
		} else if justPressed && g.geyserRect().Overlaps(image.Rect(mx, my, mx+1, my+1)) {
			g.camX = oldX
			g.camY = oldY
			g.dragging = false
			g.showGeyserList = true
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

func (g *Game) Draw(screen *ebiten.Image) {
	if g.showGeyserList {
		screen.Fill(color.Black)
		g.drawGeyserList(screen)
		g.needsRedraw = false
		g.lastDraw = time.Now()
		return
	}
	if g.loading || (len(g.biomes) == 0 && g.status != "") {
		screen.Fill(color.RGBA{30, 30, 30, 255})
		msg := g.status
		scale := 1.0
		if msg == "" {
			msg = "Fetching asteroid data..."
			scale = 2.0
		}
		w, h := textDimensions(msg)
		x := g.width/2 - int(float64(w)*scale/2)
		y := g.height/2 - int(float64(h)*scale/2)
		if g.statusError {
			if scale == 1.0 {
				drawTextWithBGBorder(screen, msg, x, y, errorBorderColor)
			} else {
				drawTextWithBGBorderScale(screen, msg, x, y, errorBorderColor, scale)
			}
		} else {
			if scale == 1.0 {
				drawTextWithBG(screen, msg, x, y)
			} else {
				drawTextWithBGScale(screen, msg, x, y, scale)
			}
		}
		g.lastDraw = time.Now()
		return
	}
	if g.needsRedraw {
		screen.Fill(color.RGBA{30, 30, 30, 255})
		if g.textures && g.biomeTextures != nil {
			if tex := g.biomeTextures["Space"]; tex != nil {
				clr := color.RGBA{255, 255, 255, 255}
				if c, ok := biomeColors["Space"]; ok {
					clr = c
				}
				rect := [][]Point{{
					{0, 0},
					{g.astWidth, 0},
					{g.astWidth, g.astHeight},
					{0, g.astHeight},
				}}
				drawBiomeTextured(screen, rect, tex, clr, g.camX, g.camY, g.zoom, g.filterMode())
			} else if clr, ok := biomeColors["Space"]; ok {
				vector.DrawFilledRect(screen, float32(g.camX), float32(g.camY),
					float32(float64(g.astWidth)*2*g.zoom),
					float32(float64(g.astHeight)*2*g.zoom), clr, false)
			}
		} else if clr, ok := biomeColors["Space"]; ok {
			vector.DrawFilledRect(screen, float32(g.camX), float32(g.camY),
				float32(float64(g.astWidth)*2*g.zoom),
				float32(float64(g.astHeight)*2*g.zoom), clr, false)
		}
		labels := []label{}
		var highlightGeysers []Geyser
		var highlightPOIs []PointOfInterest
		useNumbers := g.useNumbers && g.showItemNames && g.zoom < LegendZoomThreshold && !g.screenshotMode
		if g.legendMap == nil {
			g.initObjectLegend()
		}

		for _, bp := range g.biomes {
			clr, ok := biomeColors[bp.Name]
			if !ok {
				clr = color.RGBA{60, 60, 60, 255}
			}
			highlight := g.selectedBiome >= 0 && g.selectedBiome < len(g.legendBiomes) && g.legendBiomes[g.selectedBiome] == bp.Name
			texClr := clr
			if g.selectedBiome >= 0 && !highlight {
				texClr = color.RGBA{100, 100, 100, texClr.A}
			}
			if g.textures {
				if tex := g.biomeTextures[bp.Name]; tex != nil {
					drawBiomeTextured(screen, bp.Polygons, tex, texClr, g.camX, g.camY, g.zoom, g.filterMode())
				} else {
					drawBiome(screen, bp.Polygons, texClr, g.camX, g.camY, g.zoom)
				}
			} else {
				drawBiome(screen, bp.Polygons, texClr, g.camX, g.camY, g.zoom)
			}
			outlineClr := color.RGBA{255, 255, 255, 128}
			drawBiomeOutline(screen, bp.Polygons, g.camX, g.camY, g.zoom, outlineClr)
		}
		for _, gy := range g.geysers {
			x := math.Round((float64(gy.X) * 2 * g.zoom) + g.camX)
			y := math.Round((float64(gy.Y) * 2 * g.zoom) + g.camY)

			name := displayGeyser(gy.ID)
			formatted, width := formatLabel(name)
			dotClr := color.RGBA{}
			if idx, ok := g.legendMap["g"+name]; ok && idx-1 < len(g.legendColors) {
				dotClr = g.legendColors[idx-1]
			}
			labelClr := dotClr
			if !useNumbers {
				labelClr = color.RGBA{}
			}
			hover := false
			if idx, ok := g.legendMap["g"+name]; ok && g.selectedItem == idx-1 {
				hover = true
				highlightGeysers = append(highlightGeysers, gy)
				continue
			}

			if iconName := iconForGeyser(gy.ID); iconName != "" {
				if img, ok := g.icons[iconName]; ok && img != nil {
					op := &ebiten.DrawImageOptions{Filter: g.filterMode()}
					maxDim := math.Max(float64(img.Bounds().Dx()), float64(img.Bounds().Dy()))
					scale := g.zoom * IconScale * g.iconScale * float64(BaseIconPixels) / maxDim
					op.GeoM.Scale(scale, scale)
					w := float64(img.Bounds().Dx()) * scale
					h := float64(img.Bounds().Dy()) * scale
					left := math.Round(x - w/2)
					top := math.Round(y - h/2)
					op.GeoM.Translate(left, top)
					screen.DrawImage(img, op)
					if hover {
						vector.StrokeRect(screen, float32(left)+0.5, float32(top)+0.5, float32(math.Round(w))-1, float32(math.Round(h))-1, 2, dotClr, false)
					}
					if useNumbers {
						formatted = strconv.Itoa(g.legendMap["g"+name])
						width = len(formatted)
					}
					if g.showItemNames || useNumbers {
						if g.showItemNames || useNumbers {
							fs := fontScale()
							labels = append(labels, label{formatted, int(x) - int(float64(width)*fs/2), int(y+h/2) + 2, width, labelClr})
						}
					}
					continue
				}
			}

			vector.DrawFilledRect(screen, float32(x-2), float32(y-2), 4, 4, dotClr, true)
			if hover {
				vector.StrokeRect(screen, float32(x-3), float32(y-3), 6, 6, 2, dotClr, false)
			}
			if useNumbers {
				formatted = strconv.Itoa(g.legendMap["g"+name])
				width = len(formatted)
			}
			if g.showItemNames || useNumbers {
				if g.showItemNames || useNumbers {
					fs := fontScale()
					labels = append(labels, label{formatted, int(x) - int(float64(width)*fs/2), int(y) + 4, width, labelClr})
				}
			}
		}
		for _, poi := range g.pois {
			x := math.Round((float64(poi.X) * 2 * g.zoom) + g.camX)
			y := math.Round((float64(poi.Y) * 2 * g.zoom) + g.camY)

			name := displayPOI(poi.ID)
			formatted, width := formatLabel(name)
			dotClr := color.RGBA{}
			if idx, ok := g.legendMap["p"+name]; ok && idx-1 < len(g.legendColors) {
				dotClr = g.legendColors[idx-1]
			}
			labelClr := dotClr
			if !useNumbers {
				labelClr = color.RGBA{}
			}
			hover := false
			if idx, ok := g.legendMap["p"+name]; ok && g.selectedItem == idx-1 {
				hover = true
				highlightPOIs = append(highlightPOIs, poi)
				continue
			}

			if iconName := iconForPOI(poi.ID); iconName != "" {
				if img, ok := g.icons[iconName]; ok && img != nil {
					op := &ebiten.DrawImageOptions{Filter: g.filterMode()}
					maxDim := math.Max(float64(img.Bounds().Dx()), float64(img.Bounds().Dy()))
					scale := g.zoom * IconScale * g.iconScale * float64(BaseIconPixels) / maxDim
					op.GeoM.Scale(scale, scale)
					w := float64(img.Bounds().Dx()) * scale
					h := float64(img.Bounds().Dy()) * scale
					op.GeoM.Translate(x-w/2, y-h/2)
					screen.DrawImage(img, op)
					if hover {
						vector.StrokeRect(screen, float32(x-w/2), float32(y-h/2), float32(w), float32(h), 2, dotClr, false)
					}
					if useNumbers {
						formatted = strconv.Itoa(g.legendMap["p"+name])
						width = len(formatted)
					}
					fs := fontScale()
					labels = append(labels, label{formatted, int(x) - int(float64(width)*fs/2), int(y+h/2) + 2, width, labelClr})
					continue
				}
			}

			vector.DrawFilledRect(screen, float32(x-2), float32(y-2), 4, 4, dotClr, true)
			if hover {
				vector.StrokeRect(screen, float32(x-3), float32(y-3), 6, 6, 2, dotClr, false)
			}
			if useNumbers {
				formatted = strconv.Itoa(g.legendMap["p"+name])
				width = len(formatted)
			}
			fs := fontScale()
			labels = append(labels, label{formatted, int(x) - int(float64(width)*fs/2), int(y) + 4, width, labelClr})
		}

		labelScale := g.uiScale()
		fs := fontScale()
		for _, l := range labels {
			x := l.x
			if labelScale != 1.0 {
				x -= int(float64(l.width) * fs * (labelScale - 1) / 2)
			}
			if l.clr.A != 0 {
				if labelScale == 1.0 {
					drawTextWithBGBorder(screen, l.text, x, l.y, l.clr)
				} else {
					drawTextWithBGBorderScale(screen, l.text, x, l.y, l.clr, labelScale)
				}
			} else {
				if labelScale == 1.0 {
					drawTextWithBG(screen, l.text, x, l.y)
				} else {
					drawTextWithBGScale(screen, l.text, x, l.y, labelScale)
				}
			}
		}

		if g.coord != "" && !g.screenshotMode {
			label := g.coord
			scale := g.uiScale()
			aName := g.asteroidID
			if aName == "" {
				aName = "Unknown"
			}
			astName := fmt.Sprintf("Asteroid: %s", aName)
			rect := g.asteroidInfoRect()
			drawFrame(screen, rect)

			w, _ := textDimensions(astName)
			x := g.width/2 - int(float64(w)*scale/2)
			drawTextWithBGScale(screen, astName, x, int(30*scale), scale)
			ar := g.asteroidArrowRect()
			drawDownArrow(screen, ar, g.showAstMenu)

			w, _ = textDimensions(label)
			x = g.width/2 - int(float64(w)*scale/2)
			drawTextWithBGScale(screen, label, x, 10, scale)
		}

		if g.showLegend {
			if g.legend == nil {
				g.legend, g.legendBiomes = buildLegendImage(g.biomes)
			}
			opLegend := &ebiten.DrawImageOptions{}
			scale := g.uiScale()
			opLegend.GeoM.Scale(scale, scale)
			opLegend.GeoM.Translate(0, -g.biomeScroll)
			screen.DrawImage(g.legend, opLegend)
			lw := int(math.Round(float64(g.legend.Bounds().Dx()) * scale))
			lh := int(math.Round(float64(g.legend.Bounds().Dy()) * scale))
			ly := int(math.Round(-g.biomeScroll))
			strokeFrame(screen, image.Rect(0, ly, lw, ly+lh))
			if g.selectedBiome >= 0 {
				spacing := float64(rowSpacing())
				y0 := math.Round((10+spacing+spacing*float64(g.selectedBiome))*scale - g.biomeScroll)
				h := math.Round(spacing * scale)
				w := math.Round(float64(g.legend.Bounds().Dx()) * scale)
				vector.StrokeRect(screen, 0.5, float32(y0)-4, float32(w)-1, float32(h)-1, 2, color.RGBA{255, 0, 0, 255}, false)
			}
			if useNumbers && !g.screenshotMode {
				g.drawNumberLegend(screen)
			}
		}

		if !g.screenshotMode {
			tray := g.bottomTrayRect()
			vector.DrawFilledRect(screen, float32(tray.Min.X), float32(tray.Min.Y), float32(tray.Dx()), float32(tray.Dy()), bottomTrayColor, false)
			vector.StrokeRect(screen, float32(tray.Min.X)+0.5, float32(tray.Min.Y)+0.5, float32(tray.Dx())-1, float32(tray.Dy())-1, 1, buttonBorderColor, false)
			size := g.iconSize()
			sr := g.screenshotRect()
			scx := float32(sr.Min.X + size/2)
			scy := float32(sr.Min.Y + size/2)
			vector.DrawFilledCircle(screen, scx, scy, float32(size)/2, color.RGBA{0, 0, 0, 180}, true)
			vector.StrokeCircle(screen, scx, scy, float32(size)/2, 1, buttonBorderColor, true)
			if cam, ok := g.icons["../icons/camera.png"]; ok && cam != nil {
				op := &ebiten.DrawImageOptions{Filter: g.filterMode()}
				scale := float64(size) / math.Max(float64(cam.Bounds().Dx()), float64(cam.Bounds().Dy()))
				op.GeoM.Scale(scale, scale)
				w := float64(cam.Bounds().Dx()) * scale
				h := float64(cam.Bounds().Dy()) * scale
				op.GeoM.Translate(float64(sr.Min.X)+(float64(size)-w)/2, float64(sr.Min.Y)+(float64(size)-h)/2)
				screen.DrawImage(cam, op)
			}

			hr := g.helpRect()
			cx := float32(hr.Min.X + size/2)
			cy := float32(hr.Min.Y + size/2)
			vector.DrawFilledCircle(screen, cx, cy, float32(size)/2, color.RGBA{0, 0, 0, 180}, true)
			vector.StrokeCircle(screen, cx, cy, float32(size)/2, 1, buttonBorderColor, true)
			if helpImg, ok := g.icons["../icons/help.png"]; ok && helpImg != nil {
				op := &ebiten.DrawImageOptions{Filter: g.filterMode()}
				sc := float64(size) / math.Max(float64(helpImg.Bounds().Dx()), float64(helpImg.Bounds().Dy()))
				op.GeoM.Scale(sc, sc)
				w := float64(helpImg.Bounds().Dx()) * sc
				h := float64(helpImg.Bounds().Dy()) * sc
				op.GeoM.Translate(float64(hr.Min.X)+(float64(size)-w)/2, float64(hr.Min.Y)+(float64(size)-h)/2)
				screen.DrawImage(helpImg, op)
			}

			or := g.optionsRect()
			ocx := float32(or.Min.X + size/2)
			ocy := float32(or.Min.Y + size/2)
			vector.DrawFilledCircle(screen, ocx, ocy, float32(size)/2, color.RGBA{0, 0, 0, 180}, true)
			vector.StrokeCircle(screen, ocx, ocy, float32(size)/2, 1, buttonBorderColor, true)
			if gear, ok := g.icons["../icons/gear.png"]; ok && gear != nil {
				op := &ebiten.DrawImageOptions{Filter: g.filterMode()}
				sc := float64(size) / math.Max(float64(gear.Bounds().Dx()), float64(gear.Bounds().Dy()))
				op.GeoM.Scale(sc, sc)
				w := float64(gear.Bounds().Dx()) * sc
				h := float64(gear.Bounds().Dy()) * sc
				op.GeoM.Translate(float64(or.Min.X)+(float64(size)-w)/2, float64(or.Min.Y)+(float64(size)-h)/2)
				screen.DrawImage(gear, op)
			}

			gr := g.geyserRect()
			gcx := float32(gr.Min.X + size/2)
			gcy := float32(gr.Min.Y + size/2)
			vector.DrawFilledCircle(screen, gcx, gcy, float32(size)/2, color.RGBA{0, 0, 0, 180}, true)
			vector.StrokeCircle(screen, gcx, gcy, float32(size)/2, 1, buttonBorderColor, true)
			if icon, ok := g.icons["geyser_water.png"]; ok && icon != nil {
				op := &ebiten.DrawImageOptions{Filter: g.filterMode()}
				sc := float64(size) / math.Max(float64(icon.Bounds().Dx()), float64(icon.Bounds().Dy()))
				op.GeoM.Scale(sc, sc)
				w := float64(icon.Bounds().Dx()) * sc
				h := float64(icon.Bounds().Dy()) * sc
				op.GeoM.Translate(float64(gr.Min.X)+(float64(size)-w)/2, float64(gr.Min.Y)+(float64(size)-h)/2)
				screen.DrawImage(icon, op)
			}

			if g.hoverIcon != hoverNone {
				scale := g.uiScale()
				switch g.hoverIcon {
				case hoverScreenshot:
					g.drawTooltip(screen, "Screenshot", sr, scale)
				case hoverHelp:
					g.drawTooltip(screen, "Help", hr, scale)
				case hoverOptions:
					g.drawTooltip(screen, "Options", or, scale)
				case hoverGeysers:
					g.drawTooltip(screen, "Geyser List", gr, scale)
				}
			}
		}

		highlightLabels := []label{}
		if g.showItemNames {
			for _, gy := range highlightGeysers {
				x := math.Round((float64(gy.X) * 2 * g.zoom) + g.camX)
				y := math.Round((float64(gy.Y) * 2 * g.zoom) + g.camY)

				name := displayGeyser(gy.ID)
				formatted, width := formatLabel(name)
				dotClr := color.RGBA{255, 0, 0, 255}
				labelClr := dotClr
				if !useNumbers {
					labelClr = color.RGBA{}
				}

				if iconName := iconForGeyser(gy.ID); iconName != "" {
					if img, ok := g.icons[iconName]; ok && img != nil {
						op := &ebiten.DrawImageOptions{Filter: g.filterMode()}
						maxDim := math.Max(float64(img.Bounds().Dx()), float64(img.Bounds().Dy()))
						scale := g.zoom * IconScale * g.iconScale * float64(BaseIconPixels) / maxDim
						op.GeoM.Scale(scale, scale)
						w := float64(img.Bounds().Dx()) * scale
						h := float64(img.Bounds().Dy()) * scale
						left := math.Round(x - w/2)
						top := math.Round(y - h/2)
						op.GeoM.Translate(left, top)
						screen.DrawImage(img, op)
						vector.StrokeRect(screen, float32(left)+0.5, float32(top)+0.5, float32(math.Round(w))-1, float32(math.Round(h))-1, 2, dotClr, false)
						if useNumbers {
							formatted = strconv.Itoa(g.legendMap["g"+name])
							width = len(formatted)
						}
						fs := fontScale()
						highlightLabels = append(highlightLabels, label{formatted, int(x) - int(float64(width)*fs/2), int(y+h/2) + 2, width, labelClr})
						continue
					}
				}

				vector.DrawFilledRect(screen, float32(x-2), float32(y-2), 4, 4, dotClr, true)
				vector.StrokeRect(screen, float32(x-3), float32(y-3), 6, 6, 2, dotClr, false)
				if useNumbers {
					formatted = strconv.Itoa(g.legendMap["g"+name])
					width = len(formatted)
				}
				fs := fontScale()
				highlightLabels = append(highlightLabels, label{formatted, int(x) - int(float64(width)*fs/2), int(y) + 4, width, labelClr})
			}

			for _, poi := range highlightPOIs {
				x := math.Round((float64(poi.X) * 2 * g.zoom) + g.camX)
				y := math.Round((float64(poi.Y) * 2 * g.zoom) + g.camY)

				name := displayPOI(poi.ID)
				formatted, width := formatLabel(name)
				dotClr := color.RGBA{255, 0, 0, 255}
				labelClr := dotClr
				if !useNumbers {
					labelClr = color.RGBA{}
				}

				if iconName := iconForPOI(poi.ID); iconName != "" {
					if img, ok := g.icons[iconName]; ok && img != nil {
						op := &ebiten.DrawImageOptions{Filter: g.filterMode()}
						maxDim := math.Max(float64(img.Bounds().Dx()), float64(img.Bounds().Dy()))
						scale := g.zoom * IconScale * g.iconScale * float64(BaseIconPixels) / maxDim
						op.GeoM.Scale(scale, scale)
						w := float64(img.Bounds().Dx()) * scale
						h := float64(img.Bounds().Dy()) * scale
						left := math.Round(x - w/2)
						top := math.Round(y - h/2)
						op.GeoM.Translate(left, top)
						screen.DrawImage(img, op)
						vector.StrokeRect(screen, float32(left)+0.5, float32(top)+0.5, float32(math.Round(w))-1, float32(math.Round(h))-1, 2, dotClr, false)
						if useNumbers {
							formatted = strconv.Itoa(g.legendMap["p"+name])
							width = len(formatted)
						}
						fs := fontScale()
						highlightLabels = append(highlightLabels, label{formatted, int(x) - int(float64(width)*fs/2), int(y+h/2) + 2, width, labelClr})
						continue
					}
				}

				vector.DrawFilledRect(screen, float32(x-2), float32(y-2), 4, 4, dotClr, true)
				vector.StrokeRect(screen, float32(x-3), float32(y-1), 6, 6, 2, dotClr, false)
				if useNumbers {
					formatted = strconv.Itoa(g.legendMap["p"+name])
					width = len(formatted)
				}
				fs := fontScale()
				highlightLabels = append(highlightLabels, label{formatted, int(x) - int(float64(width)*fs/2), int(y) + 4, width, labelClr})
			}

		}
		if len(highlightLabels) > 0 {
			labelScale := g.uiScale()
			fs := fontScale()
			for _, l := range highlightLabels {
				x := l.x
				if labelScale != 1.0 {
					x -= int(float64(l.width) * fs * (labelScale - 1) / 2)
				}
				if l.clr.A != 0 {
					if labelScale == 1.0 {
						drawTextWithBGBorder(screen, l.text, x, l.y, l.clr)
					} else {
						drawTextWithBGBorderScale(screen, l.text, x, l.y, l.clr, labelScale)
					}
				} else {
					if labelScale == 1.0 {
						drawTextWithBG(screen, l.text, x, l.y)
					} else {
						drawTextWithBGScale(screen, l.text, x, l.y, labelScale)
					}
				}
			}
		}

		if !g.screenshotMode {
			cx := g.width / 2
			cy := g.height / 2
			crossClr := color.RGBA{255, 255, 255, 30}
			vector.StrokeLine(screen, float32(cx-CrosshairSize), float32(cy), float32(cx+CrosshairSize), float32(cy), 1, crossClr, true)
			vector.StrokeLine(screen, float32(cx), float32(cy-CrosshairSize), float32(cx), float32(cy+CrosshairSize), 1, crossClr, true)
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
	if g.screenshotPath != "" && !g.captured {
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
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	screenW := outsideWidth
	screenH := outsideHeight

	if g.width != screenW || g.height != screenH {
		// Keep the world position at the center of the screen fixed so
		// resizing doesn't shift the view.
		cxOld, cyOld := float64(g.width)/2, float64(g.height)/2
		worldX := (cxOld - g.camX) / g.zoom
		worldY := (cyOld - g.camY) / g.zoom

		g.width = screenW
		g.height = screenH

		cxNew, cyNew := float64(g.width)/2, float64(g.height)/2
		g.camX = cxNew - worldX*g.zoom
		g.camY = cyNew - worldY*g.zoom
		g.clampCamera()
		if g.astWidth > 0 && g.astHeight > 0 {
			zx := float64(g.width) / (float64(g.astWidth) * 2)
			zy := float64(g.height) / (float64(g.astHeight) * 2)
			g.minZoom = math.Min(zx, zy) * 0.25
		}

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
	}
	return screenW, screenH
}
