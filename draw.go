//go:build !test

package main

import (
	"fmt"
	"image/color"
	"math"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

func (g *Game) Draw(screen *ebiten.Image) {
	if g.drawGeyserListScreen(screen) {
		return
	}
	if g.drawLoadingScreen(screen) {
		return
	}
	if g.needsRedraw {
		screen.Fill(backgroundColor)
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
		useNumbers := g.useNumbers && g.zoom < LegendZoomThreshold && !g.screenshotMode && g.showItemNames
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
			formatted, _ := formatLabel(name)
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
					}
					if g.showItemNames || useNumbers {
						if g.showItemNames || useNumbers {
							labels = append(labels, label{formatted, int(x), int(y+h/2) + 2, labelClr})
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
			}
			if g.showItemNames || useNumbers {
				if g.showItemNames || useNumbers {
					labels = append(labels, label{formatted, int(x), int(y) + 4, labelClr})
				}
			}
		}
		for _, poi := range g.pois {
			x := math.Round((float64(poi.X) * 2 * g.zoom) + g.camX)
			y := math.Round((float64(poi.Y) * 2 * g.zoom) + g.camY)

			name := displayPOI(poi.ID)
			formatted, _ := formatLabel(name)
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
					}
					labels = append(labels, label{formatted, int(x), int(y+h/2) + 2, labelClr})
					continue
				}
			}

			vector.DrawFilledRect(screen, float32(x-2), float32(y-2), 4, 4, dotClr, true)
			if hover {
				vector.StrokeRect(screen, float32(x-3), float32(y-3), 6, 6, 2, dotClr, false)
			}
			if useNumbers {
				formatted = strconv.Itoa(g.legendMap["p"+name])
			}
			labels = append(labels, label{formatted, int(x), int(y) + 4, labelClr})
		}

		for _, l := range labels {
			if l.clr.A != 0 {
				drawTextWithBGBorderScale(screen, l.text, l.x, l.y, l.clr, 1, true)
			} else {
				drawTextWithBGScale(screen, l.text, l.x, l.y, 1, true)
			}
		}

		if g.coord != "" && !g.screenshotMode && g.showItemNames {
			label := g.coord
			aName := g.asteroidID
			if aName == "" {
				aName = "Unknown"
			}
			astName := fmt.Sprintf("Asteroid: %s", aName)

			x := g.width / 2
			astBase := seedBaseline() + notoFont.Metrics().Height.Ceil() + 4
			drawTextWithBGScale(screen, astName, x, astBase, 1, true)
			ar := g.asteroidArrowRect()
			drawDownArrow(screen, ar, g.showAstMenu)
			drawTextWithBGScale(screen, label, x, seedBaseline(), 1, true)
		}

		if g.showLegend {
			if g.legend == nil {
				g.legend, g.legendBiomes = buildLegendImage(g.biomes)
			}
			opLegend := &ebiten.DrawImageOptions{}
			opLegend.GeoM.Translate(0, -g.biomeScroll)
			screen.DrawImage(g.legend, opLegend)
			if g.selectedBiome >= 0 {
				spacing := float64(rowSpacing())
				y0 := math.Round((10 + spacing + spacing*float64(g.selectedBiome)) - g.biomeScroll)
				h := math.Round(spacing)
				w := math.Round(float64(g.legend.Bounds().Dx()))
				vector.StrokeRect(screen, 0.5, float32(y0)-4, float32(w)-1, float32(h)-1, 2, highlightColor, false)
			}
			if useNumbers && !g.screenshotMode {
				g.drawNumberLegend(screen)
			}
		}

		if !g.screenshotMode {
			tray := g.bottomTrayRect()
			vector.DrawFilledRect(screen, float32(tray.Min.X), float32(tray.Min.Y), float32(tray.Dx()), float32(tray.Dy()), bottomTrayColor, false)
			size := g.iconSize()
			sr := g.screenshotRect()
			scx := float32(sr.Min.X + size/2)
			scy := float32(sr.Min.Y + size/2)
			vector.DrawFilledCircle(screen, scx, scy, float32(size)/2, overlayColor, true)
			if cam, ok := g.icons["../icons/camera.png"]; ok && cam != nil {
				op := &ebiten.DrawImageOptions{Filter: g.filterMode()}
				scale := float64(size) / math.Max(float64(cam.Bounds().Dx()), float64(cam.Bounds().Dy()))
				op.GeoM.Scale(scale, scale)
				w := float64(cam.Bounds().Dx()) * scale
				h := float64(cam.Bounds().Dy()) * scale
				op.GeoM.Translate(float64(sr.Min.X)+(float64(size)-w)/2, float64(sr.Min.Y)+(float64(size)-h)/2)
				screen.DrawImage(cam, op)
			}
			vector.StrokeCircle(screen, scx, scy, float32(size)/2, 1, buttonBorderColor, true)

			hr := g.helpRect()
			cx := float32(hr.Min.X + size/2)
			cy := float32(hr.Min.Y + size/2)
			vector.DrawFilledCircle(screen, cx, cy, float32(size)/2, overlayColor, true)
			if helpImg, ok := g.icons["../icons/help.png"]; ok && helpImg != nil {
				op := &ebiten.DrawImageOptions{Filter: g.filterMode()}
				sc := float64(size) / math.Max(float64(helpImg.Bounds().Dx()), float64(helpImg.Bounds().Dy()))
				op.GeoM.Scale(sc, sc)
				w := float64(helpImg.Bounds().Dx()) * sc
				h := float64(helpImg.Bounds().Dy()) * sc
				op.GeoM.Translate(float64(hr.Min.X)+(float64(size)-w)/2, float64(hr.Min.Y)+(float64(size)-h)/2)
				screen.DrawImage(helpImg, op)
			}
			vector.StrokeCircle(screen, cx, cy, float32(size)/2, 1, buttonBorderColor, true)

			or := g.optionsRect()
			ocx := float32(or.Min.X + size/2)
			ocy := float32(or.Min.Y + size/2)
			vector.DrawFilledCircle(screen, ocx, ocy, float32(size)/2, overlayColor, true)
			if gear, ok := g.icons["../icons/gear.png"]; ok && gear != nil {
				op := &ebiten.DrawImageOptions{Filter: g.filterMode()}
				sc := float64(size) / math.Max(float64(gear.Bounds().Dx()), float64(gear.Bounds().Dy()))
				op.GeoM.Scale(sc, sc)
				w := float64(gear.Bounds().Dx()) * sc
				h := float64(gear.Bounds().Dy()) * sc
				op.GeoM.Translate(float64(or.Min.X)+(float64(size)-w)/2, float64(or.Min.Y)+(float64(size)-h)/2)
				screen.DrawImage(gear, op)
			}
			vector.StrokeCircle(screen, ocx, ocy, float32(size)/2, 1, buttonBorderColor, true)

			gr := g.geyserRect()
			gcx := float32(gr.Min.X + size/2)
			gcy := float32(gr.Min.Y + size/2)
			vector.DrawFilledCircle(screen, gcx, gcy, float32(size)/2, overlayColor, true)
			if icon, ok := g.icons["geyser_water.png"]; ok && icon != nil {
				op := &ebiten.DrawImageOptions{Filter: g.filterMode()}
				sc := float64(size) / math.Max(float64(icon.Bounds().Dx()), float64(icon.Bounds().Dy()))
				op.GeoM.Scale(sc, sc)
				w := float64(icon.Bounds().Dx()) * sc
				h := float64(icon.Bounds().Dy()) * sc
				op.GeoM.Translate(float64(gr.Min.X)+(float64(size)-w)/2, float64(gr.Min.Y)+(float64(size)-h)/2)
				screen.DrawImage(icon, op)
			}
			vector.StrokeCircle(screen, gcx, gcy, float32(size)/2, 1, buttonBorderColor, true)

			if g.hoverIcon != hoverNone {
				switch g.hoverIcon {
				case hoverScreenshot:
					g.drawTooltip(screen, "Screenshot", sr, 1)
				case hoverHelp:
					g.drawTooltip(screen, "Help", hr, 1)
				case hoverOptions:
					g.drawTooltip(screen, "Options", or, 1)
				case hoverGeysers:
					g.drawTooltip(screen, "Geyser List", gr, 1)
				}
			}
		}

		highlightLabels := []label{}
		for _, gy := range highlightGeysers {
			x := math.Round((float64(gy.X) * 2 * g.zoom) + g.camX)
			y := math.Round((float64(gy.Y) * 2 * g.zoom) + g.camY)

			name := displayGeyser(gy.ID)
			formatted, _ := formatLabel(name)
			dotClr := highlightColor
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
					}
					if g.showItemNames {
						highlightLabels = append(highlightLabels, label{formatted, int(x), int(y+h/2) + 2, labelClr})
					}
					continue
				}
			}

			vector.DrawFilledRect(screen, float32(x-2), float32(y-2), 4, 4, dotClr, true)
			vector.StrokeRect(screen, float32(x-3), float32(y-3), 6, 6, 2, dotClr, false)
			if useNumbers {
				formatted = strconv.Itoa(g.legendMap["g"+name])
			}
			if g.showItemNames {
				highlightLabels = append(highlightLabels, label{formatted, int(x), int(y) + 4, labelClr})
			}
		}

		for _, poi := range highlightPOIs {
			x := math.Round((float64(poi.X) * 2 * g.zoom) + g.camX)
			y := math.Round((float64(poi.Y) * 2 * g.zoom) + g.camY)

			name := displayPOI(poi.ID)
			formatted, _ := formatLabel(name)
			dotClr := highlightColor
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
					}
					if g.showItemNames {
						highlightLabels = append(highlightLabels, label{formatted, int(x), int(y+h/2) + 2, labelClr})
					}
					continue
				}
			}

			vector.DrawFilledRect(screen, float32(x-2), float32(y-2), 4, 4, dotClr, true)
			vector.StrokeRect(screen, float32(x-3), float32(y-1), 6, 6, 2, dotClr, false)
			if useNumbers {
				formatted = strconv.Itoa(g.legendMap["p"+name])
			}
			if g.showItemNames {
				highlightLabels = append(highlightLabels, label{formatted, int(x), int(y) + 4, labelClr})
			}
		}

		if len(highlightLabels) > 0 {
			for _, l := range highlightLabels {
				if l.clr.A != 0 {
					drawTextWithBGBorderScale(screen, l.text, l.x, l.y, l.clr, 1, true)
				} else {
					drawTextWithBGScale(screen, l.text, l.x, l.y, 1, true)
				}
			}
		}

		g.drawUI(screen)
	}
	g.captureScreen(screen)

}
