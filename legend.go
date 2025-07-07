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
			maxW = uiScaled(20+10) + w
		}
	}

	spacing := rowSpacing()
	width := 30 + maxW + 5
	height := spacing*(len(names)+2) + 7

	img := ebiten.NewImage(width, height)
	img.Fill(legendBGColor)

	y := 10
	drawTextWithBG(img, "Biomes", 5, y, false)
	y += spacing
	for _, name := range names {
		clr, ok := biomeColors[name]
		if !ok {
			clr = color.RGBA{60, 60, 60, 255}
		}
		vector.DrawFilledRect(img, 5, float32(y), float32(uiScaledF(20)), float32(uiScaledF(10)), clr, false)
		drawTextWithBGBorder(img, displayBiome(name), uiScaled(20+10), y, clr, false)
		y += spacing
	}

	drawTextWithBGBorder(img, "Clear", 5, y, buttonBorderColor, false)

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
	g.geyserItems = nil
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
		img.Fill(legendBGColor)
		y := 10
		drawTextWithBG(img, "Items", 5, y, false)
		y += spacing
		for i, e := range g.legendEntries {
			clr := color.RGBA{}
			if i < len(g.legendColors) {
				clr = g.legendColors[i]
			}
			drawTextWithBGBorder(img, e, 5, y, clr, false)
			y += spacing
		}
		drawTextWithBGBorder(img, "Clear", 5, y, buttonBorderColor, false)
		g.legendImage = img
	}
	w := float64(g.legendImage.Bounds().Dx())
	x := float64(dst.Bounds().Dx()) - w - 12
	y := 10.0 - g.itemScroll
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(math.Round(x), math.Round(y))
	dst.DrawImage(g.legendImage, op)
	if g.selectedItem >= 0 {
		spacing := float64(rowSpacing())
		hy := y + (10 + spacing*float64(g.selectedItem+1))
		hh := spacing
		vector.StrokeRect(dst, float32(math.Round(x))+0.5, float32(math.Round(hy))-4, float32(math.Round(w))-1, float32(math.Round(hh))-1, 2, highlightColor, false)
	}
}

func (g *Game) buildGeyserItems() {
	if g.geyserItems != nil && len(g.geyserItems) == len(g.geysers) {
		return
	}
	g.geyserItems = make([]geyserListItem, len(g.geysers))
	for i, gy := range g.geysers {
		ic := (*ebiten.Image)(nil)
		if n := iconForGeyser(gy.ID); n != "" {
			ic = g.icons[n]
		}
		txt := displayGeyser(gy.ID) + "\n" + formatGeyserInfo(gy)
		w, h := infoRowSize(txt, ic)
		g.geyserItems[i] = geyserListItem{text: txt, icon: ic, w: w, h: h}
	}
}

func (g *Game) drawGeyserList(dst *ebiten.Image) {
	vector.DrawFilledRect(dst, 0, 0, float32(g.width), float32(g.height), color.RGBA{0, 0, 0, 255}, false)
	cr := g.geyserCloseRect()
	drawCloseButton(dst, cr)

	spacing := uiScaled(10)
	g.buildGeyserItems()
	maxW := 0
	for _, it := range g.geyserItems {
		if it.w > maxW {
			maxW = it.w
		}
	}
	cols := 1
	if maxW+spacing > 0 {
		cols = g.width / (maxW + spacing)
		if cols < 1 {
			cols = 1
		}
	}
	rows := (len(g.geyserItems) + cols - 1) / cols
	rowHeights := make([]int, rows)
	for i, it := range g.geyserItems {
		r := i / cols
		if it.h > rowHeights[r] {
			rowHeights[r] = it.h
		}
	}
	y := spacing - int(g.geyserScroll)
	idx := 0
	for r := 0; r < rows; r++ {
		rowH := rowHeights[r]
		x := spacing
		for c := 0; c < cols && idx < len(g.geyserItems); c++ {
			it := g.geyserItems[idx]
			if y+rowH > 0 && y < g.height && x+it.w > 0 && x < g.width {
				g.drawInfoRow(dst, it.text, it.icon, x, y)
			}
			x += maxW + spacing
			idx++
		}
		y += rowH + spacing
	}

	if max := g.maxGeyserScroll(); max > 0 {
		barW := uiScaled(ScrollBarWidth)
		barX := g.width - barW - uiScaled(2)
		h := float64(g.height)
		barH := h * h / (h + max)
		barY := (g.geyserScroll / max) * (h - barH)
		left := barX - uiScaled(1)
		right := barX + barW + uiScaled(1)
		vector.StrokeLine(dst, float32(left), 0, float32(left), float32(g.height), 1, scrollBarTrackColor, true)
		vector.StrokeLine(dst, float32(right), 0, float32(right), float32(g.height), 1, scrollBarTrackColor, true)
		vector.DrawFilledRect(dst, float32(barX), float32(barY), float32(barW), float32(barH), scrollBarColor, false)
		r := float32(barW) / 2
		vector.DrawFilledCircle(dst, float32(barX)+r, float32(barY), r, scrollBarColor, true)
		vector.DrawFilledCircle(dst, float32(barX)+r, float32(barY)+float32(barH), r, scrollBarColor, true)
	}
}

func (g *Game) maxGeyserScroll() float64 {
	spacing := uiScaled(10)
	g.buildGeyserItems()
	maxW := 0
	for _, it := range g.geyserItems {
		if it.w > maxW {
			maxW = it.w
		}
	}
	cols := 1
	if maxW+spacing > 0 {
		cols = g.width / (maxW + spacing)
		if cols < 1 {
			cols = 1
		}
	}
	rows := (len(g.geyserItems) + cols - 1) / cols
	if rows == 0 {
		return 0
	}
	rowHeights := make([]int, rows)
	for i, it := range g.geyserItems {
		r := i / cols
		if it.h > rowHeights[r] {
			rowHeights[r] = it.h
		}
	}
	total := spacing
	for _, h := range rowHeights {
		total += h + spacing
	}
	extra := rowHeights[rows-1] + spacing
	max := float64(total - g.height + extra)
	if max < 0 {
		max = 0
	}
	return max
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
	h := float64(g.legend.Bounds().Dy())
	extra := float64(rowSpacing() * LegendScrollExtraRows)
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
	h := float64(g.legendImage.Bounds().Dy())
	extra := float64(rowSpacing() * LegendScrollExtraRows)
	max := h + 10 - float64(g.height) + extra
	if max < 0 {
		max = 0
	}
	return max
}

func (g *Game) itemPanelVisible() bool {
	return g.showLegend && g.useNumbers && g.zoom < LegendZoomThreshold && !g.screenshotMode
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
	if g.legend != nil && len(g.legendBiomes) > 0 {
		w := g.legend.Bounds().Dx()
		if mx >= 0 && mx < w {
			spacing := float64(rowSpacing())
			baseY := int(10+spacing) - int(g.biomeScroll)
			rowH := int(spacing)
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
	useNumbers := g.useNumbers && g.zoom < LegendZoomThreshold && !g.screenshotMode
	if useNumbers && g.legendImage != nil {
		w := g.legendImage.Bounds().Dx()
		x0 := g.width - w - 12
		if mx >= x0 && mx < x0+w {
			spacing := float64(rowSpacing())
			baseY := int(10+spacing) - int(g.itemScroll)
			rowH := int(spacing)
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
	if g.bottomTrayRect().Overlaps(image.Rect(mx, my, mx+1, my+1)) {
		return false
	}
	handled := false
	if g.legend != nil && len(g.legendBiomes) > 0 {
		w := g.legend.Bounds().Dx()
		if mx >= 0 && mx < w {
			spacing := float64(rowSpacing())
			baseY := int(10+spacing) - int(g.biomeScroll)
			rowH := int(spacing)
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
	useNumbers := g.useNumbers && g.zoom < LegendZoomThreshold && !g.screenshotMode
	if useNumbers && g.legendImage != nil {
		w := g.legendImage.Bounds().Dx()
		x0 := g.width - w - 12
		if mx >= x0 && mx < x0+w {
			spacing := float64(rowSpacing())
			baseY := int(10+spacing) - int(g.itemScroll)
			rowH := int(spacing)
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
