//go:build !test

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math"
	"os"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"time"

	_ "embed"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const helpMessage = "Controls:\n" +
	"Arrow keys/WASD or drag to pan\n" +
	"Mouse wheel or +/- to zoom\n" +
	"Pinch to zoom on touch"

var errorBorderColor = color.RGBA{R: 244, G: 67, B: 54, A: 255}

func drawTextWithBG(dst *ebiten.Image, text string, x, y int) {
	lines := strings.Split(text, "\n")
	width := 0
	for _, l := range lines {
		if len(l) > width {
			width = len(l)
		}
	}
	height := len(lines) * 16
	vector.DrawFilledRect(dst, float32(x-2), float32(y-2), float32(width*6+4), float32(height+4), color.RGBA{0, 0, 0, 128}, false)
	ebitenutil.DebugPrintAt(dst, text, x, y)
}

func drawTextWithBGScale(dst *ebiten.Image, text string, x, y int, scale float64) {
	if scale == 1.0 {
		drawTextWithBG(dst, text, x, y)
		return
	}
	lines := strings.Split(text, "\n")
	width := 0
	for _, l := range lines {
		if len(l) > width {
			width = len(l)
		}
	}
	height := len(lines) * 16
	w := width*6 + 4
	h := height + 4
	img := ebiten.NewImage(w, h)
	vector.DrawFilledRect(img, 0, 0, float32(w), float32(h), color.RGBA{0, 0, 0, 128}, false)
	ebitenutil.DebugPrintAt(img, text, 2, 2)
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(scale, scale)
	op.GeoM.Translate(float64(x-2), float64(y-2))
	dst.DrawImage(img, op)
}

func drawTextWithBGBorder(dst *ebiten.Image, text string, x, y int, border color.Color) {
	lines := strings.Split(text, "\n")
	width := 0
	for _, l := range lines {
		if len(l) > width {
			width = len(l)
		}
	}
	height := len(lines) * 16
	bx := x - 2
	by := y - 2
	bw := width*6 + 4
	bh := height + 4
	vector.DrawFilledRect(dst, float32(bx-1), float32(by-1), float32(bw+2), float32(bh+2), border, false)
	vector.DrawFilledRect(dst, float32(bx), float32(by), float32(bw), float32(bh), color.RGBA{0, 0, 0, 128}, false)
	ebitenutil.DebugPrintAt(dst, text, x, y)
}

func drawTextWithBGBorderScale(dst *ebiten.Image, text string, x, y int, border color.Color, scale float64) {
	if scale == 1.0 {
		drawTextWithBGBorder(dst, text, x, y, border)
		return
	}
	lines := strings.Split(text, "\n")
	width := 0
	for _, l := range lines {
		if len(l) > width {
			width = len(l)
		}
	}
	height := len(lines) * 16
	w := width*6 + 4
	h := height + 4
	img := ebiten.NewImage(w+2, h+2)
	vector.DrawFilledRect(img, 0, 0, float32(w+2), float32(h+2), border, false)
	vector.DrawFilledRect(img, 1, 1, float32(w), float32(h), color.RGBA{0, 0, 0, 128}, false)
	ebitenutil.DebugPrintAt(img, text, 3, 3)
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(scale, scale)
	op.GeoM.Translate(float64(x-2), float64(y-2))
	dst.DrawImage(img, op)
}

func drawInfoPanel(dst *ebiten.Image, text string, icon *ebiten.Image, x, y int, scale float64) {
	lines := strings.Split(text, "\n")
	width := 0
	for _, l := range lines {
		if len(l) > width {
			width = len(l)
		}
	}
	txtW := width * LabelCharWidth
	txtH := len(lines) * 16
	iconW, iconH := 0, 0
	if icon != nil {
		iconW = InfoIconSize
		iconH = InfoIconSize
	}
	gap := 4
	w := txtW + iconW + gap + 4
	h := txtH
	if iconH > txtH {
		h = iconH
	}
	h += 4
	img := ebiten.NewImage(w, h)
	vector.DrawFilledRect(img, 0, 0, float32(w), float32(h), color.RGBA{0, 0, 0, InfoPanelAlpha}, false)
	if icon != nil {
		opIcon := &ebiten.DrawImageOptions{Filter: ebiten.FilterLinear}
		scaleIcon := float64(InfoIconSize) / math.Max(float64(icon.Bounds().Dx()), float64(icon.Bounds().Dy()))
		opIcon.GeoM.Scale(scaleIcon, scaleIcon)
		opIcon.GeoM.Translate(2, float64(h-iconH)/2)
		img.DrawImage(icon, opIcon)
	}
	ebitenutil.DebugPrintAt(img, text, iconW+gap+2, 2)
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(scale, scale)
	op.GeoM.Translate(float64(x-2), float64(y-2))
	dst.DrawImage(img, op)
}

func textDimensions(text string) (int, int) {
	lines := strings.Split(text, "\n")
	width := 0
	for _, l := range lines {
		if len(l) > width {
			width = len(l)
		}
	}
	return width * LabelCharWidth, len(lines) * 16
}

func infoPanelSize(text string, icon *ebiten.Image) (int, int) {
	lines := strings.Split(text, "\n")
	width := 0
	for _, l := range lines {
		if len(l) > width {
			width = len(l)
		}
	}
	txtW := width * LabelCharWidth
	txtH := len(lines) * 16
	iconW, iconH := 0, 0
	if icon != nil {
		iconW = InfoIconSize
		iconH = InfoIconSize
	}
	gap := 4
	w := txtW + iconW + gap + 4
	h := txtH
	if iconH > txtH {
		h = iconH
	}
	h += 4
	return w, h
}

func drawBiome(dst *ebiten.Image, polys [][]Point, clr color.Color, camX, camY, zoom float64) {
	if len(polys) == 0 {
		return
	}
	var p vector.Path
	for _, pts := range polys {
		if len(pts) == 0 {
			continue
		}
		p.MoveTo(float32(pts[0].X*2), float32(pts[0].Y*2))
		for _, pt := range pts[1:] {
			p.LineTo(float32(pt.X*2), float32(pt.Y*2))
		}
		p.Close()
	}
	vs, is := p.AppendVerticesAndIndicesForFilling(nil, nil)
	r, g, b, a := clr.RGBA()
	for i := range vs {
		x := float64(vs[i].DstX)*zoom + camX
		y := float64(vs[i].DstY)*zoom + camY
		vs[i].DstX = float32(x)
		vs[i].DstY = float32(y)
		vs[i].SrcX = 0
		vs[i].SrcY = 0
		vs[i].ColorR = float32(r) / 0xffff
		vs[i].ColorG = float32(g) / 0xffff
		vs[i].ColorB = float32(b) / 0xffff
		vs[i].ColorA = float32(a) / 0xffff
	}
	op := &ebiten.DrawTrianglesOptions{
		AntiAlias:      true,
		ColorScaleMode: ebiten.ColorScaleModePremultipliedAlpha,
		FillRule:       ebiten.FillRuleEvenOdd,
	}
	dst.DrawTriangles(vs, is, whitePixel, op)
}

func drawBiomeTextured(dst *ebiten.Image, polys [][]Point, tex *ebiten.Image, clr color.Color, camX, camY, zoom float64) {
	if len(polys) == 0 || tex == nil {
		return
	}
	var p vector.Path
	for _, pts := range polys {
		if len(pts) == 0 {
			continue
		}
		p.MoveTo(float32(pts[0].X*2), float32(pts[0].Y*2))
		for _, pt := range pts[1:] {
			p.LineTo(float32(pt.X*2), float32(pt.Y*2))
		}
		p.Close()
	}
	vs, is := p.AppendVerticesAndIndicesForFilling(nil, nil)
	r, g, b, a := clr.RGBA()
	texW, texH := tex.Bounds().Dx(), tex.Bounds().Dy()
	for i := range vs {
		worldX := float64(vs[i].DstX) / 2
		worldY := float64(vs[i].DstY) / 2
		x := float64(vs[i].DstX)*zoom + camX
		y := float64(vs[i].DstY)*zoom + camY
		vs[i].DstX = float32(x)
		vs[i].DstY = float32(y)
		vs[i].SrcX = float32(worldX*BiomeTextureScale) * float32(texW)
		vs[i].SrcY = float32(worldY*BiomeTextureScale) * float32(texH)
		vs[i].ColorR = float32(r) / 0xffff
		vs[i].ColorG = float32(g) / 0xffff
		vs[i].ColorB = float32(b) / 0xffff
		vs[i].ColorA = float32(a) / 0xffff
	}
	op := &ebiten.DrawTrianglesOptions{
		AntiAlias:      true,
		ColorScaleMode: ebiten.ColorScaleModePremultipliedAlpha,
		FillRule:       ebiten.FillRuleEvenOdd,
		Address:        ebiten.AddressRepeat,
		Filter:         ebiten.FilterLinear,
	}
	dst.DrawTriangles(vs, is, tex, op)
}

func drawBiomeOutline(dst *ebiten.Image, polys [][]Point, camX, camY, zoom float64, clr color.Color) {
	for _, pts := range polys {
		if len(pts) < 2 {
			continue
		}
		for i := 0; i < len(pts); i++ {
			a := pts[i]
			b := pts[(i+1)%len(pts)]
			x0 := float32(math.Round(float64(a.X*2)*zoom + camX))
			y0 := float32(math.Round(float64(a.Y*2)*zoom + camY))
			x1 := float32(math.Round(float64(b.X*2)*zoom + camX))
			y1 := float32(math.Round(float64(b.Y*2)*zoom + camY))
			vector.StrokeLine(dst, x0, y0, x1, y1, 1, clr, true)
		}
	}
}

// buildLegendImage returns an image of biome colors with a single background.
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

	maxLen := 0
	for _, name := range names {
		if l := len(displayBiome(name)); l > maxLen {
			maxLen = l
		}
	}

	width := 30 + maxLen*6 + 5
	height := LegendRowSpacing*(len(names)+1) + 7

	img := ebiten.NewImage(width, height)
	img.Fill(color.RGBA{0, 0, 0, 77})

	y := 10
	drawTextWithBG(img, "Biomes", 5, y)
	y += LegendRowSpacing
	for _, name := range names {
		clr, ok := biomeColors[name]
		if !ok {
			clr = color.RGBA{60, 60, 60, 255}
		}
		vector.DrawFilledRect(img, 5, float32(y), 20, 10, clr, false)
		drawTextWithBGBorder(img, displayBiome(name), 30, y, clr)
		y += LegendRowSpacing
	}

	return img, names
}

// initObjectLegend prepares the mapping of object names to numeric labels and
// caches the legend text.
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

// drawNumberLegend draws the cached legend image on the top right corner with
// a single semi-transparent background.
func (g *Game) drawNumberLegend(dst *ebiten.Image) {
	if len(g.legendEntries) == 0 {
		return
	}
	if g.legendImage == nil {
		maxLen := 0
		for _, e := range g.legendEntries {
			if len(e) > maxLen {
				maxLen = len(e)
			}
		}
		width := maxLen*6 + 10
		height := LegendRowSpacing*(len(g.legendEntries)+1) + 7
		img := ebiten.NewImage(width, height)
		img.Fill(color.RGBA{0, 0, 0, 77})
		y := 10
		drawTextWithBG(img, "Items", 5, y)
		y += LegendRowSpacing
		for i, e := range g.legendEntries {
			clr := color.RGBA{}
			if i < len(g.legendColors) {
				clr = g.legendColors[i]
			}
			drawTextWithBGBorder(img, e, 5, y, clr)
			y += LegendRowSpacing
		}
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
	if g.hoverItem >= 0 {
		hy := y + float64(10+LegendRowSpacing*(g.hoverItem+1))*scale
		hh := float64(LegendRowSpacing) * scale
		vector.StrokeRect(dst, float32(math.Round(x))+0.5, float32(math.Round(hy))-4, float32(math.Round(w))-1, float32(math.Round(hh))-1, 2, color.RGBA{255, 0, 0, 255}, false)
	}
}

func (g *Game) drawGeyserList(dst *ebiten.Image) {
	vector.DrawFilledRect(dst, 0, 0, float32(g.width), float32(g.height), color.RGBA{0, 0, 0, 255}, false)
	scale := g.uiScale()
	cr := g.geyserCloseRect()
	vector.DrawFilledRect(dst, float32(cr.Min.X), float32(cr.Min.Y), float32(cr.Dx()), float32(cr.Dy()), color.RGBA{0, 0, 0, 200}, false)
	vector.StrokeRect(dst, float32(cr.Min.X)+0.5, float32(cr.Min.Y)+0.5, float32(cr.Dx())-1, float32(cr.Dy())-1, 2, color.RGBA{255, 255, 255, 255}, false)
	vector.StrokeLine(dst, float32(cr.Min.X+4), float32(cr.Min.Y+4), float32(cr.Max.X-4), float32(cr.Max.Y-4), 2, color.RGBA{255, 255, 255, 255}, false)
	vector.StrokeLine(dst, float32(cr.Min.X+4), float32(cr.Max.Y-4), float32(cr.Max.X-4), float32(cr.Min.Y+4), 2, color.RGBA{255, 255, 255, 255}, false)

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
		w, h := infoPanelSize(txt, ic)
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
			drawInfoPanel(dst, it.text, it.icon, x, y, scale)
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
		w, h := infoPanelSize(txt, ic)
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

func (g *Game) maxBiomeScroll() float64 {
	if g.legend == nil {
		return 0
	}
	scale := g.uiScale()
	h := float64(g.legend.Bounds().Dy()) * scale
	extra := float64(LegendRowSpacing*LegendScrollExtraRows) * scale
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
	extra := float64(LegendRowSpacing*LegendScrollExtraRows) * scale
	max := h + 10 - float64(g.height) + extra
	if max < 0 {
		max = 0
	}
	return max
}

func (g *Game) updateHover(mx, my int) {
	prevBiome := g.hoverBiome
	prevItem := g.hoverItem
	g.hoverBiome = -1
	g.hoverItem = -1
	scale := g.uiScale()
	if g.legend != nil && len(g.legendBiomes) > 0 {
		w := int(float64(g.legend.Bounds().Dx()) * scale)
		if mx >= 0 && mx < w {
			baseY := int(float64(10+LegendRowSpacing)*scale) - int(g.biomeScroll)
			rowH := int(float64(LegendRowSpacing) * scale)
			for i := range g.legendBiomes {
				ry := baseY + i*rowH
				if my >= ry && my < ry+rowH {
					g.hoverBiome = i
					break
				}
			}
		}
	}
	useNumbers := !g.mobile && g.zoom < LegendZoomThreshold && !g.screenshotMode
	if useNumbers && g.legendImage != nil {
		w := int(float64(g.legendImage.Bounds().Dx()) * scale)
		x0 := g.width - w - 12
		if mx >= x0 && mx < x0+w {
			baseY := int(float64(10+LegendRowSpacing)*scale) - int(g.itemScroll)
			rowH := int(float64(LegendRowSpacing) * scale)
			for i := range g.legendEntries {
				ry := baseY + i*rowH
				if my >= ry && my < ry+rowH {
					g.hoverItem = i
					break
				}
			}
		}
	}
	if g.hoverBiome != prevBiome || g.hoverItem != prevItem {
		g.needsRedraw = true
	}
}

// Game implements ebiten.Game and displays geysers with their names.
type Game struct {
	geysers           []Geyser
	pois              []PointOfInterest
	biomes            []BiomePath
	icons             map[string]*ebiten.Image
	biomeTextures     map[string]*ebiten.Image
	width             int
	height            int
	astWidth          int
	astHeight         int
	camX              float64
	camY              float64
	zoom              float64
	minZoom           float64
	dragging          bool
	lastX             int
	lastY             int
	touches           map[ebiten.TouchID]touchPoint
	pinchDist         float64
	needsRedraw       bool
	screenshotPath    string
	captured          bool
	legend            *ebiten.Image
	legendMap         map[string]int
	legendEntries     []string
	legendColors      []color.RGBA
	legendImage       *ebiten.Image
	legendBiomes      []string
	hoverBiome        int
	hoverItem         int
	showGeyserList    bool
	geyserScroll      float64
	biomeScroll       float64
	itemScroll        float64
	showHelp          bool
	lastWheel         time.Time
	loading           bool
	status            string
	iconResults       chan loadedIcon
	coord             string
	mobile            bool
	showInfo          bool
	infoPinned        bool
	infoText          string
	infoX             int
	infoY             int
	infoIcon          *ebiten.Image
	lastMouseX        int
	lastMouseY        int
	mousePrev         bool
	touchUsed         bool
	touchActive       bool
	touchStartX       int
	touchStartY       int
	touchMoved        bool
	touchUI           bool
	showShotMenu      bool
	screenshotMode    bool
	ssQuality         int
	ssSaved           time.Time
	ssPending         int
	skipClickTicks    int
	lastDraw          time.Time
	wasMinimized      bool
	magnify           bool
	asteroidID        int
	asteroidSpecified bool
	statusError       bool
}

type label struct {
	text  string
	x     int
	y     int
	width int
	clr   color.RGBA
}

type touchPoint struct {
	x int
	y int
}

type loadedIcon struct {
	name string
	img  *ebiten.Image
}

// uiScale returns a multiplier for UI elements based on the current window
// height. Larger windows and high quality screenshots use bigger values so
// that text and icons remain readable.
func (g *Game) uiScale() float64 {
	if g.mobile {
		return 1.0
	}
	if g.screenshotMode {
		if g.height > 1700 {
			return 4.0
		}
		if g.height > 850 {
			return 2.0
		}
		return 1.0
	}
	if g.magnify {
		return 2.0
	}
	return 1.0
}

func (g *Game) iconSize() int {
	return int(float64(HelpIconSize) * g.uiScale())
}

func drawPlusMinus(dst *ebiten.Image, rect image.Rectangle, minus bool) {
	cx := float32(rect.Min.X + rect.Dx()/2)
	cy := float32(rect.Min.Y + rect.Dy()/2)
	length := float32(rect.Dx()) * 0.5
	thickness := float32(rect.Dx()) / 8
	vector.StrokeLine(dst, cx-length/2, cy, cx+length/2, cy, thickness, color.RGBA{255, 255, 255, 255}, true)
	if !minus {
		vector.StrokeLine(dst, cx, cy-length/2, cx, cy+length/2, thickness, color.RGBA{255, 255, 255, 255}, true)
	}
}

func (g *Game) magnifyRect() image.Rectangle {
	size := g.iconSize()
	x := g.width - size - HelpMargin
	y := g.height - size - HelpMargin
	return image.Rect(x, y, x+size, y+size)
}

func (g *Game) helpRect() image.Rectangle {
	size := g.iconSize()
	x := g.width - size*2 - HelpMargin*2
	y := g.height - size - HelpMargin
	return image.Rect(x, y, x+size, y+size)
}

func (g *Game) geyserRect() image.Rectangle {
	size := g.iconSize()
	x := g.width - size*4 - HelpMargin*4
	y := g.height - size - HelpMargin
	return image.Rect(x, y, x+size, y+size)
}

func (g *Game) geyserCloseRect() image.Rectangle {
	size := g.iconSize()
	x := g.width - size - HelpMargin
	y := HelpMargin
	return image.Rect(x, y, x+size, y+size)
}

func (g *Game) clampCamera() {
	if g.astWidth == 0 || g.astHeight == 0 {
		return
	}
	cx := float64(g.width) / 2
	cy := float64(g.height) / 2
	maxX := cx
	minX := cx - float64(g.astWidth)*2*g.zoom
	if g.camX < minX {
		g.camX = minX
	}
	if g.camX > maxX {
		g.camX = maxX
	}
	maxY := cy
	minY := cy - float64(g.astHeight)*2*g.zoom
	if g.camY < minY {
		g.camY = minY
	}
	if g.camY > maxY {
		g.camY = maxY
	}
}

func (g *Game) itemAt(mx, my int) (string, int, int, *ebiten.Image, bool) {
	const hitRadius = 10
	for _, gy := range g.geysers {
		x := float64(gy.X)*2*g.zoom + g.camX
		y := float64(gy.Y)*2*g.zoom + g.camY
		left, top, right, bottom := x-hitRadius, y-hitRadius, x+hitRadius, y+hitRadius
		if iconName := iconForGeyser(gy.ID); iconName != "" {
			if img, ok := g.icons[iconName]; ok && img != nil {
				maxDim := math.Max(float64(img.Bounds().Dx()), float64(img.Bounds().Dy()))
				scale := g.zoom * IconScale * float64(BaseIconPixels) / maxDim
				w := float64(img.Bounds().Dx()) * scale
				h := float64(img.Bounds().Dy()) * scale
				left = x - w/2
				top = y - h/2
				right = x + w/2
				bottom = y + h/2
			}
		}
		if float64(mx) >= left && float64(mx) <= right && float64(my) >= top && float64(my) <= bottom {
			info := displayGeyser(gy.ID) + "\n" + formatGeyserInfo(gy)
			var icon *ebiten.Image
			if n := iconForGeyser(gy.ID); n != "" {
				icon = g.icons[n]
			}
			return info, int(math.Round(x)), int(math.Round(y)), icon, true
		}
	}
	for _, poi := range g.pois {
		x := float64(poi.X)*2*g.zoom + g.camX
		y := float64(poi.Y)*2*g.zoom + g.camY
		left, top, right, bottom := x-hitRadius, y-hitRadius, x+hitRadius, y+hitRadius
		if iconName := iconForPOI(poi.ID); iconName != "" {
			if img, ok := g.icons[iconName]; ok && img != nil {
				maxDim := math.Max(float64(img.Bounds().Dx()), float64(img.Bounds().Dy()))
				scale := g.zoom * IconScale * float64(BaseIconPixels) / maxDim
				w := float64(img.Bounds().Dx()) * scale
				h := float64(img.Bounds().Dy()) * scale
				left = x - w/2
				top = y - h/2
				right = x + w/2
				bottom = y + h/2
			}
		}
		if float64(mx) >= left && float64(mx) <= right && float64(my) >= top && float64(my) <= bottom {
			info := displayPOI(poi.ID) + "\n" + formatPOIInfo(poi)
			var icon *ebiten.Image
			if n := iconForPOI(poi.ID); n != "" {
				icon = g.icons[n]
			}
			return info, int(math.Round(x)), int(math.Round(y)), icon, true
		}
	}
	return "", 0, 0, nil, false
}

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

func (g *Game) startIconLoader(names []string) {
	g.iconResults = make(chan loadedIcon, len(names))
	go func() {
		for _, n := range names {
			img, _ := loadImageFile(n)
			g.iconResults <- loadedIcon{name: n, img: img}
		}
		close(g.iconResults)
	}()
}

func (g *Game) Update() error {
	const panSpeed = PanSpeed

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
			g.geyserScroll -= float64(wheelY) * 10
			if g.geyserScroll < 0 {
				g.geyserScroll = 0
			}
			max := g.maxGeyserScroll()
			if g.geyserScroll > max {
				g.geyserScroll = max
			}
			g.needsRedraw = true
		}
		mx, my := ebiten.CursorPosition()
		if mx >= 0 && mx < g.width && my >= 0 && my < g.height {
			if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
				if g.geyserCloseRect().Overlaps(image.Rect(mx, my, mx+1, my+1)) {
					g.showGeyserList = false
					g.needsRedraw = true
				}
				g.mousePrev = true
			} else {
				g.mousePrev = false
			}
		} else {
			g.mousePrev = false
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
	if mxTmp < 0 || mxTmp >= g.width || myTmp < 0 || myTmp >= g.height {
		mousePressed = false
	}
	if g.ssPending > 0 || g.skipClickTicks > 0 {
		mousePressed = false
	}

	justPressed := mousePressed && !g.mousePrev

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
	if g.ssPending > 0 || g.skipClickTicks > 0 {
		touchIDs = nil
	}
	// Filter out touches outside the window
	valid := make([]ebiten.TouchID, 0, len(touchIDs))
	for _, id := range touchIDs {
		x, y := ebiten.TouchPosition(id)
		if x >= 0 && x < g.width && y >= 0 && y < g.height {
			valid = append(valid, id)
		}
	}
	touchIDs = valid
	if len(touchIDs) > 0 {
		g.touchUsed = true
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
					g.geyserScroll -= float64(dy)
					if g.geyserScroll < 0 {
						g.geyserScroll = 0
					}
					if max := g.maxGeyserScroll(); g.geyserScroll > max {
						g.geyserScroll = max
					}
					g.needsRedraw = true
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
					useNumbers := !g.mobile && g.zoom < LegendZoomThreshold && !g.screenshotMode
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
			if g.showGeyserList || g.showShotMenu {
				g.touchUI = true
			} else {
				pt := image.Rect(x, y, x+1, y+1)
				if g.helpRect().Overlaps(pt) || g.screenshotRect().Overlaps(pt) ||
					g.magnifyRect().Overlaps(pt) || g.geyserRect().Overlaps(pt) {
					g.touchUI = true
				} else {
					scale := g.uiScale()
					if g.legend != nil {
						lw := int(float64(g.legend.Bounds().Dx()) * scale)
						if x >= 0 && x < lw {
							g.touchUI = true
						}
					}
					useNumbers := !g.mobile && g.zoom < LegendZoomThreshold && !g.screenshotMode
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
			} else if g.screenshotRect().Overlaps(pt) {
				g.showShotMenu = true
				g.needsRedraw = true
			} else if g.magnifyRect().Overlaps(pt) {
				g.magnify = !g.magnify
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
			} else if g.helpRect().Overlaps(pt) {
				g.showHelp = !g.showHelp
				g.needsRedraw = true
			} else if g.geyserRect().Overlaps(pt) {
				g.camX = oldX
				g.camY = oldY
				g.dragging = false
				g.showGeyserList = true
				g.needsRedraw = true
			} else if g.touchUI {
				g.updateHover(mx, my)
			} else if g.mobile {
				if _, ix, iy, _, found := g.itemAt(mx, my); found {
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
			g.geyserScroll -= float64(wheelY) * 10
			if g.geyserScroll < 0 {
				g.geyserScroll = 0
			}
			max := g.maxGeyserScroll()
			if g.geyserScroll > max {
				g.geyserScroll = max
			}
		} else {
			handled := false
			scale := g.uiScale()
			useNumbers := !g.mobile && g.zoom < LegendZoomThreshold && !g.screenshotMode
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
	if !g.mobile && !g.screenshotMode {
		mx, my = mxTmp, myTmp
		if mx < 0 || mx >= g.width || my < 0 || my >= g.height {
			mx, my = -1, -1
		} else if mx != g.lastMouseX || my != g.lastMouseY {
			g.touchUsed = false
		}
		g.lastMouseX, g.lastMouseY = mx, my
		if g.showHelp {
			if !g.helpRect().Overlaps(image.Rect(mx, my, mx+1, my+1)) {
				g.showHelp = false
				g.needsRedraw = true
			}
		} else if justPressed && g.helpRect().Overlaps(image.Rect(mx, my, mx+1, my+1)) {
			g.showHelp = true
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
		} else if justPressed && g.screenshotRect().Overlaps(image.Rect(mx, my, mx+1, my+1)) {
			g.showShotMenu = true
			g.needsRedraw = true
		} else if justPressed && g.magnifyRect().Overlaps(image.Rect(mx, my, mx+1, my+1)) {
			g.magnify = !g.magnify
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
		} else if justPressed && g.geyserRect().Overlaps(image.Rect(mx, my, mx+1, my+1)) {
			g.camX = oldX
			g.camY = oldY
			g.dragging = false
			g.showGeyserList = true
			g.needsRedraw = true
		}
	}

	if !g.mobile && !g.screenshotMode && !g.touchUsed {
		g.updateHover(mx, my)
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

	g.clampCamera()

	if g.camX != oldX || g.camY != oldY || g.zoom != oldZoom {
		g.needsRedraw = true
	}

	g.mousePrev = mousePressed
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
		if msg == "" {
			msg = "Fetching..."
		}
		scale := 1.0
		if msg == "Fetching..." {
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
			drawBiomeTextured(screen, rect, tex, clr, g.camX, g.camY, g.zoom)
		} else if clr, ok := biomeColors["Space"]; ok {
			vector.DrawFilledRect(screen, float32(g.camX), float32(g.camY),
				float32(float64(g.astWidth)*2*g.zoom),
				float32(float64(g.astHeight)*2*g.zoom), clr, false)
		}
		labels := []label{}
		var highlightGeysers []Geyser
		var highlightPOIs []PointOfInterest
		useNumbers := !g.mobile && g.zoom < LegendZoomThreshold && !g.screenshotMode
		if g.legendMap == nil {
			g.initObjectLegend()
		}

		for _, bp := range g.biomes {
			clr, ok := biomeColors[bp.Name]
			if !ok {
				clr = color.RGBA{60, 60, 60, 255}
			}
			highlight := g.hoverBiome >= 0 && g.hoverBiome < len(g.legendBiomes) && g.legendBiomes[g.hoverBiome] == bp.Name
			texClr := clr
			if g.hoverBiome >= 0 && !highlight {
				texClr = color.RGBA{100, 100, 100, texClr.A}
			}
			if tex := g.biomeTextures[bp.Name]; tex != nil {
				drawBiomeTextured(screen, bp.Polygons, tex, texClr, g.camX, g.camY, g.zoom)
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
			if idx, ok := g.legendMap["g"+name]; ok && g.hoverItem == idx-1 {
				hover = true
				highlightGeysers = append(highlightGeysers, gy)
				continue
			}

			if iconName := iconForGeyser(gy.ID); iconName != "" {
				if img, ok := g.icons[iconName]; ok && img != nil {
					op := &ebiten.DrawImageOptions{Filter: ebiten.FilterLinear}
					maxDim := math.Max(float64(img.Bounds().Dx()), float64(img.Bounds().Dy()))
					scale := g.zoom * IconScale * float64(BaseIconPixels) / maxDim
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
					labels = append(labels, label{formatted, int(x) - (width*LabelCharWidth)/2, int(y+h/2) + 2, width, labelClr})
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
			labels = append(labels, label{formatted, int(x) - (width*LabelCharWidth)/2, int(y) + 4, width, labelClr})
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
			if idx, ok := g.legendMap["p"+name]; ok && g.hoverItem == idx-1 {
				hover = true
				highlightPOIs = append(highlightPOIs, poi)
				continue
			}

			if iconName := iconForPOI(poi.ID); iconName != "" {
				if img, ok := g.icons[iconName]; ok && img != nil {
					op := &ebiten.DrawImageOptions{Filter: ebiten.FilterLinear}
					maxDim := math.Max(float64(img.Bounds().Dx()), float64(img.Bounds().Dy()))
					scale := g.zoom * IconScale * float64(BaseIconPixels) / maxDim
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
					labels = append(labels, label{formatted, int(x) - (width*LabelCharWidth)/2, int(y+h/2) + 2, width, labelClr})
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
			labels = append(labels, label{formatted, int(x) - (width*LabelCharWidth)/2, int(y) + 4, width, labelClr})
		}

		labelScale := g.uiScale()
		for _, l := range labels {
			x := l.x
			if labelScale != 1.0 {
				x -= int(float64(l.width*LabelCharWidth) * (labelScale - 1) / 2)
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
			if g.asteroidSpecified {
				label += fmt.Sprintf("  ast: %d", g.asteroidID)
			}
			scale := g.uiScale()
			x := g.width/2 - int(float64(len(label)*LabelCharWidth)*scale/2)
			drawTextWithBGScale(screen, label, x, 10, scale)
		}

		if g.legend == nil {
			g.legend, g.legendBiomes = buildLegendImage(g.biomes)
		}
		opLegend := &ebiten.DrawImageOptions{}
		scale := g.uiScale()
		opLegend.GeoM.Scale(scale, scale)
		opLegend.GeoM.Translate(0, -g.biomeScroll)
		screen.DrawImage(g.legend, opLegend)
		if g.hoverBiome >= 0 {
			y0 := math.Round(float64(10+LegendRowSpacing+g.hoverBiome*LegendRowSpacing)*scale - g.biomeScroll)
			h := math.Round(float64(LegendRowSpacing) * scale)
			w := math.Round(float64(g.legend.Bounds().Dx()) * scale)
			vector.StrokeRect(screen, 0.5, float32(y0)-4, float32(w)-1, float32(h)-1, 2, color.RGBA{255, 0, 0, 255}, false)
		}
		if useNumbers && !g.screenshotMode {
			g.drawNumberLegend(screen)
		}

		if !g.mobile && !g.screenshotMode {
			size := g.iconSize()
			sr := g.screenshotRect()
			scx := float32(sr.Min.X + size/2)
			scy := float32(sr.Min.Y + size/2)
			vector.DrawFilledCircle(screen, scx, scy, float32(size)/2, color.RGBA{0, 0, 0, 180}, true)
			if cam, ok := g.icons["../icons/camera.png"]; ok && cam != nil {
				op := &ebiten.DrawImageOptions{Filter: ebiten.FilterLinear}
				scale := float64(size) / math.Max(float64(cam.Bounds().Dx()), float64(cam.Bounds().Dy()))
				op.GeoM.Scale(scale, scale)
				w := float64(cam.Bounds().Dx()) * scale
				h := float64(cam.Bounds().Dy()) * scale
				op.GeoM.Translate(float64(sr.Min.X)+(float64(size)-w)/2, float64(sr.Min.Y)+(float64(size)-h)/2)
				screen.DrawImage(cam, op)
			}
			if g.showShotMenu {
				g.drawScreenshotMenu(screen)
			}

			mr := g.magnifyRect()
			mcx := float32(mr.Min.X + size/2)
			mcy := float32(mr.Min.Y + size/2)
			vector.DrawFilledCircle(screen, mcx, mcy, float32(size)/2, color.RGBA{0, 0, 0, 180}, true)
			drawPlusMinus(screen, mr, g.magnify)

			hr := g.helpRect()
			cx := float32(hr.Min.X + size/2)
			cy := float32(hr.Min.Y + size/2)
			vector.DrawFilledCircle(screen, cx, cy, float32(size)/2, color.RGBA{0, 0, 0, 180}, true)
			if helpImg, ok := g.icons["../icons/help.png"]; ok && helpImg != nil {
				op := &ebiten.DrawImageOptions{Filter: ebiten.FilterLinear}
				sc := float64(size) / math.Max(float64(helpImg.Bounds().Dx()), float64(helpImg.Bounds().Dy()))
				op.GeoM.Scale(sc, sc)
				w := float64(helpImg.Bounds().Dx()) * sc
				h := float64(helpImg.Bounds().Dy()) * sc
				op.GeoM.Translate(float64(hr.Min.X)+(float64(size)-w)/2, float64(hr.Min.Y)+(float64(size)-h)/2)
				screen.DrawImage(helpImg, op)
			}

			gr := g.geyserRect()
			gcx := float32(gr.Min.X + size/2)
			gcy := float32(gr.Min.Y + size/2)
			vector.DrawFilledCircle(screen, gcx, gcy, float32(size)/2, color.RGBA{0, 0, 0, 180}, true)
			if icon, ok := g.icons["geyser_water.png"]; ok && icon != nil {
				op := &ebiten.DrawImageOptions{Filter: ebiten.FilterLinear}
				sc := float64(size) / math.Max(float64(icon.Bounds().Dx()), float64(icon.Bounds().Dy()))
				op.GeoM.Scale(sc, sc)
				w := float64(icon.Bounds().Dx()) * sc
				h := float64(icon.Bounds().Dy()) * sc
				op.GeoM.Translate(float64(gr.Min.X)+(float64(size)-w)/2, float64(gr.Min.Y)+(float64(size)-h)/2)
				screen.DrawImage(icon, op)
			}
		}

		if g.showHelp && !g.screenshotMode {
			scale := g.uiScale()
			hr := g.helpRect()
			tx := hr.Min.X - int(170*scale)
			ty := hr.Min.Y - int(70*scale)
			drawTextWithBGScale(screen, helpMessage, tx, ty, scale)
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

		highlightLabels := []label{}
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
					op := &ebiten.DrawImageOptions{Filter: ebiten.FilterLinear}
					maxDim := math.Max(float64(img.Bounds().Dx()), float64(img.Bounds().Dy()))
					scale := g.zoom * IconScale * float64(BaseIconPixels) / maxDim
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
					highlightLabels = append(highlightLabels, label{formatted, int(x) - (width*LabelCharWidth)/2, int(y+h/2) + 2, width, labelClr})
					continue
				}
			}

			vector.DrawFilledRect(screen, float32(x-2), float32(y-2), 4, 4, dotClr, true)
			vector.StrokeRect(screen, float32(x-3), float32(y-3), 6, 6, 2, dotClr, false)
			if useNumbers {
				formatted = strconv.Itoa(g.legendMap["g"+name])
				width = len(formatted)
			}
			highlightLabels = append(highlightLabels, label{formatted, int(x) - (width*LabelCharWidth)/2, int(y) + 4, width, labelClr})
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
					op := &ebiten.DrawImageOptions{Filter: ebiten.FilterLinear}
					maxDim := math.Max(float64(img.Bounds().Dx()), float64(img.Bounds().Dy()))
					scale := g.zoom * IconScale * float64(BaseIconPixels) / maxDim
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
					highlightLabels = append(highlightLabels, label{formatted, int(x) - (width*LabelCharWidth)/2, int(y+h/2) + 2, width, labelClr})
					continue
				}
			}

			vector.DrawFilledRect(screen, float32(x-2), float32(y-2), 4, 4, dotClr, true)
			vector.StrokeRect(screen, float32(x-3), float32(y-1), 6, 6, 2, dotClr, false)
			if useNumbers {
				formatted = strconv.Itoa(g.legendMap["p"+name])
				width = len(formatted)
			}
			highlightLabels = append(highlightLabels, label{formatted, int(x) - (width*LabelCharWidth)/2, int(y) + 4, width, labelClr})
		}

		if len(highlightLabels) > 0 {
			labelScale := g.uiScale()
			for _, l := range highlightLabels {
				x := l.x
				if labelScale != 1.0 {
					x -= int(float64(l.width*LabelCharWidth) * (labelScale - 1) / 2)
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
			drawInfoPanel(screen, g.infoText, g.infoIcon, tx, ty, scale)
		}

		g.needsRedraw = false
		g.lastDraw = time.Now()
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
	if g.width != outsideWidth || g.height != outsideHeight {
		// Keep the world position at the center of the screen fixed so
		// resizing doesn't shift the view.
		cxOld, cyOld := float64(g.width)/2, float64(g.height)/2
		worldX := (cxOld - g.camX) / g.zoom
		worldY := (cyOld - g.camY) / g.zoom

		g.width = outsideWidth
		g.height = outsideHeight

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
	return outsideWidth, outsideHeight
}

func main() {
	coord := flag.String("coord", "SNDST-A-7-0-0-0", "seed coordinate")
	out := flag.String("out", "", "optional path to save JSON")
	screenshot := flag.String("screenshot", "", "path to save a PNG screenshot and exit")
	flag.Parse()
	asteroidIdx := 0
	asteroidSpecified := false
	invalidAsteroid := false
	if runtime.GOARCH == "wasm" {
		if c := coordFromURL(); c != "" {
			*coord = c
		}
		if a, ok := asteroidFromURL(); ok {
			asteroidIdx = a
			asteroidSpecified = true
			if a < 0 {
				invalidAsteroid = true
			}
		}
	}

	game := &Game{
		icons:             make(map[string]*ebiten.Image),
		width:             DefaultWidth,
		height:            DefaultHeight,
		zoom:              1.0,
		minZoom:           MinZoom,
		loading:           true,
		status:            "Fetching...",
		statusError:       false,
		coord:             *coord,
		asteroidID:        asteroidIdx,
		asteroidSpecified: asteroidSpecified,
		mobile:            isMobile(),
		ssQuality:         1,
		hoverBiome:        -1,
		hoverItem:         -1,
		mousePrev:         false,
	}
	if invalidAsteroid {
		game.status = fmt.Sprintf("%s\nAsteroid ID: %d\nThis location does not contain Asteroid ID: %d", *coord, asteroidIdx, asteroidIdx)
		game.statusError = true
		game.loading = false
	} else {
		go func(idx int) {
			fmt.Println("Fetching:", *coord)
			cborData, err := fetchSeedCBOR(*coord)
			if err != nil {
				game.status = "Error: " + err.Error()
				game.statusError = false
				game.needsRedraw = true
				game.loading = false
				return
			}
			seed, err := decodeSeed(cborData)
			if err != nil {
				game.status = "Error: " + err.Error()
				game.statusError = false
				game.needsRedraw = true
				game.loading = false
				return
			}
			if game.asteroidSpecified && (idx < 0 || idx >= len(seed.Asteroids)) {
				game.status = fmt.Sprintf("%s\nAsteroid ID: %d\nThis location does not contain Asteroid ID: %d", game.coord, idx, idx)
				game.statusError = true
				game.needsRedraw = true
				game.loading = false
				return
			}
			if *out != "" {
				jsonData, _ := json.MarshalIndent(seed, "", "  ")
				_ = saveToFile(*out, jsonData)
			}
			astIdxSel := 0
			if idx >= 0 && idx < len(seed.Asteroids) {
				astIdxSel = idx
			}
			ast := seed.Asteroids[astIdxSel]
			bps := parseBiomePaths(ast.BiomePaths)
			game.geysers = ast.Geysers
			game.pois = ast.POIs
			game.biomes = bps
			game.astWidth = ast.SizeX
			game.astHeight = ast.SizeY
			game.legend, game.legendBiomes = buildLegendImage(bps)
			zoomX := float64(game.width) / (float64(game.astWidth) * 2)
			zoomY := float64(game.height) / (float64(game.astHeight) * 2)
			game.zoom = math.Min(zoomX, zoomY)
			game.minZoom = game.zoom * 0.25
			game.camX = (float64(game.width) - float64(game.astWidth)*2*game.zoom) / 2
			game.camY = (float64(game.height) - float64(game.astHeight)*2*game.zoom) / 2
			game.clampCamera()
			game.biomeTextures = loadBiomeTextures()
			names := []string{"../icons/camera.png", "../icons/help.png", "geyser_water.png"}
			set := make(map[string]struct{})
			for _, gy := range ast.Geysers {
				if n := iconForGeyser(gy.ID); n != "" {
					if _, ok := set[n]; !ok {
						set[n] = struct{}{}
						names = append(names, n)
					}
				}
			}
			for _, poi := range ast.POIs {
				if n := iconForPOI(poi.ID); n != "" {
					if _, ok := set[n]; !ok {
						set[n] = struct{}{}
						names = append(names, n)
					}
				}
			}
			game.startIconLoader(names)
			game.loading = false
			game.needsRedraw = true
		}(asteroidIdx)
	}
	if *screenshot != "" {
		game.screenshotPath = *screenshot
		game.screenshotMode = true
		game.showHelp = false
		game.showInfo = false
		game.showShotMenu = false
	}
	ebiten.SetWindowSize(game.width, game.height)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetWindowTitle("Geysers - " + *coord)
	ebiten.SetScreenClearedEveryFrame(false)
	ebiten.SetVsyncEnabled(true)
	if err := ebiten.RunGame(game); err != nil {
		fmt.Println("Error running game:", err)
	}
}
