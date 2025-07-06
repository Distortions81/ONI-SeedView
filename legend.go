package main

import (
	"fmt"
	"image"
	"image/color"
	"math"
	"sort"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

func buildLegendImage(biomes []BiomePath) (*ebiten.Image, []string) {
	set := make(map[string]struct{})
	for _, b := range biomes {
		set[b.Name] = struct{}{}
	}
	names := make([]string, 0, len(set))
	for n := range set {
		names = append(names, n)
	}
	sort.Strings(names)

	maxW := 0
	for _, name := range names {
		w, _ := textDimensions(displayBiome(name))
		if w > maxW {
			maxW = w
		}
	}

	spacing := rowSpacing()
	width := 30 + maxW + 5
	height := spacing*(len(names)+2) + 7

	img := ebiten.NewImage(width, height)
	img.Fill(color.RGBA{0, 0, 0, 77})

	y := 10
	drawTextWithBG(img, "Biomes", 5, y)
	y += spacing
	for _, name := range names {
		clr, ok := biomeColors[name]
		if !ok {
			clr = color.RGBA{60, 60, 60, 255}
		}
		vector.DrawFilledRect(img, 5, float32(y), 20, 10, clr, false)
		drawTextWithBGBorder(img, displayBiome(name), 30, y, clr)
		y += spacing
	}

	drawTextWithBGBorder(img, "Clear", 5, y, buttonBorderColor)

	return img, names
}

func (g *Game) initObjectLegend() {
	if g.legendMap != nil {
		return
	}
	g.legendMap = make(map[string]int)
	g.legendColors = nil
	counter := 1
	for _, gy := range g.geysers {
		name := displayGeyser(gy.ID)
		if _, ok := g.legendMap["g"+name]; !ok {
			g.legendMap["g"+name] = counter
			g.legendEntries = append(g.legendEntries, fmt.Sprintf("%d: %s", counter, name))
			g.legendColors = append(g.legendColors, uniqueColor(counter))
			counter++
		}
	}
	for _, poi := range g.pois {
		name := displayPOI(poi.ID)
		if _, ok := g.legendMap["p"+name]; !ok {
			g.legendMap["p"+name] = counter
			g.legendEntries = append(g.legendEntries, fmt.Sprintf("%d: %s", counter, name))
			g.legendColors = append(g.legendColors, uniqueColor(counter))
			counter++
		}
	}
}

func (g *Game) invalidateLegends() {
	g.legend = nil
	g.legendImage = nil
}

func (g *Game) drawNumberLegend(dst *ebiten.Image) {
	if len(g.legendEntries) == 0 {
		return
	}
	if g.legendImage == nil {
		maxW := 0
		for _, e := range g.legendEntries {
			w, _ := textDimensions(e)
			if w > maxW {
				maxW = w
			}
		}
		spacing := rowSpacing()
		width := maxW + 10
		height := spacing*(len(g.legendEntries)+2) + 7
		img := ebiten.NewImage(width, height)
		img.Fill(color.RGBA{0, 0, 0, 77})
		y := 10
		drawTextWithBG(img, "Items", 5, y)
		y += spacing
		for i, e := range g.legendEntries {
			clr := color.RGBA{}
			if i < len(g.legendColors) {
				clr = g.legendColors[i]
			}
			drawTextWithBGBorder(img, e, 5, y, clr)
			y += spacing
		}
		drawTextWithBGBorder(img, "Clear", 5, y, buttonBorderColor)
		g.legendImage = img
	}
	scale := g.uiScale()
	w := float64(g.legendImage.Bounds().Dx()) * scale
	x := float64(dst.Bounds().Dx()) - w - 12
	y := 10.0 - g.itemScroll
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(scale, scale)
	op.GeoM.Translate(math.Round(x), math.Round(y))
	dst.DrawImage(g.legendImage, op)
	lh := float64(g.legendImage.Bounds().Dy()) * scale
	strokeFrame(dst, image.Rect(int(math.Round(x)), int(math.Round(y)), int(math.Round(x+w)), int(math.Round(y+lh))))
	if g.selectedItem >= 0 {
		spacing := float64(rowSpacing())
		hy := y + (10+spacing*float64(g.selectedItem+1))*scale
		hh := spacing * scale
		vector.StrokeRect(dst, float32(math.Round(x))+0.5, float32(math.Round(hy))-4, float32(math.Round(w))-1, float32(math.Round(hh))-1, 2, color.RGBA{255, 0, 0, 255}, false)
	}
}

func (g *Game) drawGeyserList(dst *ebiten.Image) {
	vector.DrawFilledRect(dst, 0, 0, float32(g.width), float32(g.height), color.RGBA{0, 0, 0, 255}, false)
	scale := g.uiScale()
	cr := g.geyserCloseRect()
	drawCloseButton(dst, cr)

	spacing := int(10 * scale)
	type item struct {
		text string
		icon *ebiten.Image
		w    int
		h    int
	}
	items := make([]item, len(g.geysers))
	maxW := 0
	for i, gy := range g.geysers {
		ic := (*ebiten.Image)(nil)
		if n := iconForGeyser(gy.ID); n != "" {
			ic = g.icons[n]
		}
		txt := displayGeyser(gy.ID) + "\n" + formatGeyserInfo(gy)
		w, h := infoRowSize(txt, ic)
		sw := int(float64(w) * scale)
		sh := int(float64(h) * scale)
		if sw > maxW {
			maxW = sw
		}
		items[i] = item{text: txt, icon: ic, w: sw, h: sh}
	}
	cols := 1
	if maxW+spacing > 0 {
		cols = g.width / (maxW + spacing)
		if cols < 1 {
			cols = 1
		}
	}
	rows := (len(items) + cols - 1) / cols
	rowHeights := make([]int, rows)
	for i, it := range items {
		r := i / cols
		if it.h > rowHeights[r] {
			rowHeights[r] = it.h
		}
	}
	y := spacing - int(g.geyserScroll)
	idx := 0
	for r := 0; r < rows; r++ {
		x := spacing
		for c := 0; c < cols && idx < len(items); c++ {
			it := items[idx]
			g.drawInfoRow(dst, it.text, it.icon, x, y, scale)
			x += maxW + spacing
			idx++
		}
		y += rowHeights[r] + spacing
	}
}

func (g *Game) maxGeyserScroll() float64 {
	scale := g.uiScale()
	spacing := int(10 * scale)
	type item struct {
		w int
		h int
	}
	items := make([]item, len(g.geysers))
	maxW := 0
	for i, gy := range g.geysers {
		ic := (*ebiten.Image)(nil)
		if n := iconForGeyser(gy.ID); n != "" {
			ic = g.icons[n]
		}
		txt := displayGeyser(gy.ID) + "\n" + formatGeyserInfo(gy)
		w, h := infoRowSize(txt, ic)
		sw := int(float64(w) * scale)
		sh := int(float64(h) * scale)
		if sw > maxW {
			maxW = sw
		}
		items[i] = item{w: sw, h: sh}
	}
	cols := 1
	if maxW+spacing > 0 {
		cols = g.width / (maxW + spacing)
		if cols < 1 {
			cols = 1
		}
	}
	rows := (len(items) + cols - 1) / cols
	if rows == 0 {
		return 0
	}
	rowHeights := make([]int, rows)
	for i, it := range items {
		r := i / cols
		if it.h > rowHeights[r] {
			rowHeights[r] = it.h
		}
	}
	total := spacing
	for _, h := range rowHeights {
		total += h + spacing
	}
	max := total - rowHeights[rows-1] - spacing
	if max < 0 {
		max = 0
	}
	return float64(max)
}

func (g *Game) adjustGeyserScroll(delta float64) {
	g.geyserScroll += delta
	if g.geyserScroll < 0 {
		g.geyserScroll = 0
	}
	if max := g.maxGeyserScroll(); g.geyserScroll > max {
		g.geyserScroll = max
	}
	g.needsRedraw = true
}

func (g *Game) maxBiomeScroll() float64 {
	if g.legend == nil {
		return 0
	}
	scale := g.uiScale()
	h := float64(g.legend.Bounds().Dy()) * scale
	extra := float64(rowSpacing()*LegendScrollExtraRows) * scale
	max := h - float64(g.height) + extra
	if max < 0 {
		max = 0
	}
	return max
}

func (g *Game) maxItemScroll() float64 {
	if g.legendImage == nil {
		return 0
	}
	scale := g.uiScale()
	h := float64(g.legendImage.Bounds().Dy()) * scale
	extra := float64(rowSpacing()*LegendScrollExtraRows) * scale
	max := h + 10 - float64(g.height) + extra
	if max < 0 {
		max = 0
	}
	return max
}

func (g *Game) itemPanelVisible() bool {
	return g.showLegend && g.useNumbers && g.showItemNames && g.zoom < LegendZoomThreshold && !g.screenshotMode
}

func (g *Game) updateHover(mx, my int) {
	prevBiome := g.hoverBiome
	prevItem := g.hoverItem
	g.hoverBiome = -1
	g.hoverItem = -1
	if !g.showLegend {
		if prevBiome != -1 || prevItem != -1 {
			g.needsRedraw = true
		}
		return
	}
	scale := g.uiScale()
	if g.legend != nil && len(g.legendBiomes) > 0 {
		w := int(float64(g.legend.Bounds().Dx()) * scale)
		if mx >= 0 && mx < w {
			spacing := float64(rowSpacing())
			baseY := int((10+spacing)*scale) - int(g.biomeScroll)
			rowH := int(spacing * scale)
			for i := range g.legendBiomes {
				ry := baseY + i*rowH
				if my >= ry && my < ry+rowH {
					g.hoverBiome = i
					break
				}
			}
			if g.hoverBiome == -1 {
				ry := baseY + len(g.legendBiomes)*rowH
				if my >= ry && my < ry+rowH {
					g.hoverBiome = len(g.legendBiomes)
				}
			}
		}
	}
	useNumbers := g.useNumbers && g.showItemNames && g.zoom < LegendZoomThreshold && !g.screenshotMode
	if useNumbers && g.legendImage != nil {
		w := int(float64(g.legendImage.Bounds().Dx()) * scale)
		x0 := g.width - w - 12
		if mx >= x0 && mx < x0+w {
			spacing := float64(rowSpacing())
			baseY := int((10+spacing)*scale) - int(g.itemScroll)
			rowH := int(spacing * scale)
			for i := range g.legendEntries {
				ry := baseY + i*rowH
				if my >= ry && my < ry+rowH {
					g.hoverItem = i
					break
				}
			}
			if g.hoverItem == -1 {
				ry := baseY + len(g.legendEntries)*rowH
				if my >= ry && my < ry+rowH {
					g.hoverItem = len(g.legendEntries)
				}
			}
		}
	}
	if g.hoverBiome != prevBiome || g.hoverItem != prevItem {
		g.needsRedraw = true
	}
}

func (g *Game) updateIconHover(mx, my int) {
	prev := g.hoverIcon
	g.hoverIcon = hoverNone
	if mx == -1 || my == -1 {
		if prev != g.hoverIcon {
			g.needsRedraw = true
		}
		return
	}
	pt := image.Rect(mx, my, mx+1, my+1)
	switch {
	case g.optionsRect().Overlaps(pt):
		g.hoverIcon = hoverOptions
	case g.geyserRect().Overlaps(pt):
		g.hoverIcon = hoverGeysers
	case g.screenshotRect().Overlaps(pt):
		g.hoverIcon = hoverScreenshot
	case g.helpRect().Overlaps(pt):
		g.hoverIcon = hoverHelp
	}
	if prev != g.hoverIcon {
		g.needsRedraw = true
	}
}

func (g *Game) clickLegend(mx, my int) bool {
	if !g.showLegend {
		return false
	}
	handled := false
	scale := g.uiScale()
	if g.legend != nil && len(g.legendBiomes) > 0 {
		w := int(float64(g.legend.Bounds().Dx()) * scale)
		if mx >= 0 && mx < w {
			spacing := float64(rowSpacing())
			baseY := int((10+spacing)*scale) - int(g.biomeScroll)
			rowH := int(spacing * scale)
			count := len(g.legendBiomes)
			for i := 0; i <= count; i++ {
				ry := baseY + i*rowH
				if my >= ry && my < ry+rowH {
					if i == count {
						g.selectedBiome = -1
					} else {
						g.selectedBiome = i
					}
					handled = true
					break
				}
			}
		}
	}
	useNumbers := g.useNumbers && g.showItemNames && g.zoom < LegendZoomThreshold && !g.screenshotMode
	if useNumbers && g.legendImage != nil {
		w := int(float64(g.legendImage.Bounds().Dx()) * scale)
		x0 := g.width - w - 12
		if mx >= x0 && mx < x0+w {
			spacing := float64(rowSpacing())
			baseY := int((10+spacing)*scale) - int(g.itemScroll)
			rowH := int(spacing * scale)
			count := len(g.legendEntries)
			for i := 0; i <= count; i++ {
				ry := baseY + i*rowH
				if my >= ry && my < ry+rowH {
					if i == count {
						g.selectedItem = -1
					} else {
						g.selectedItem = i
					}
					handled = true
					break
				}
			}
		}
	}
	if handled {
		g.needsRedraw = true
	}
	return handled
}
