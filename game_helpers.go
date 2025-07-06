package main

import (
	"image"
	"image/color"
	"math"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Game struct {
	geysers           []Geyser
	pois              []PointOfInterest
	biomes            []BiomePath
	asteroids         []Asteroid
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
	selectedBiome     int
	selectedItem      int
	showGeyserList    bool
	geyserScroll      float64
	biomeScroll       float64
	itemScroll        float64
	showHelp          bool
	lastWheel         time.Time
	loading           bool
	status            string
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
	hoverIcon         hoverIcon
	touchUsed         bool
	touchActive       bool
	touchStartX       int
	touchStartY       int
	touchMoved        bool
	touchUI           bool
	showShotMenu      bool
	showAstMenu       bool
	showOptions       bool
	asteroidScroll    float64
	screenshotMode    bool
	ssQuality         int
	ssSaved           time.Time
	ssPending         int
	skipClickTicks    int
	lastDraw          time.Time
	wasMinimized      bool
	asteroidID        string
	asteroidSpecified bool
	statusError       bool
	fitOnLoad         bool

	textures      bool
	vsync         bool
	showItemNames bool
	showLegend    bool
	useNumbers    bool
	iconScale     float64
	smartRender   bool
	linearFilter  bool

	noColor   bool
	ssNoColor bool
	grayImage *ebiten.Image
}

type label struct {
	text string
	x    int
	y    int
	clr  color.RGBA
}

type touchPoint struct {
	x int
	y int
}

type loadedIcon struct {
	name string
	img  *ebiten.Image
}

func (g *Game) uiScale() float64 { return uiScale }

func (g *Game) iconSize() int { return uiScaled(HelpIconSize) }

func (g *Game) filterMode() ebiten.Filter {
	if g.linearFilter {
		return ebiten.FilterLinear
	}
	return ebiten.FilterNearest
}

func (g *Game) inBounds(x, y int) bool {
	return x >= 0 && x < g.width && y >= 0 && y < g.height
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

func (g *Game) drawTooltip(dst *ebiten.Image, text string, rect image.Rectangle, scale float64) {
	w, h := textDimensions(text)
	x := rect.Min.X + rect.Dx()/2 - int(float64(w)*scale/2)
	y := rect.Min.Y - int(float64(h)*scale) - 4
	if x < 0 {
		x = 0
	}
	if x+int(float64(w)*scale) > g.width {
		x = g.width - int(float64(w)*scale)
	}
	if y < 0 {
		y = rect.Max.Y + 4
	}
	drawTextWithBGScale(dst, text, x, y, scale, false)
}

func (g *Game) helpRect() image.Rectangle {
	size := g.iconSize()
	x := g.width - size - uiScaled(HelpMargin)
	y := g.height - size - uiScaled(HelpMargin)
	return image.Rect(x, y, x+size, y+size)
}

func helpMenuSize() (int, int) {
	w, h := textDimensions(helpMessage)
	return w + uiScaled(4), h + uiScaled(4)
}

func (g *Game) helpMenuRect() image.Rectangle {
	w, h := helpMenuSize()
	r := g.helpRect()
	x := r.Max.X - w
	if x < 0 {
		x = 0
	}
	if x+w > g.width {
		x = g.width - w
	}
	y := r.Min.Y - h
	if y < 0 {
		y = 0
	}
	return image.Rect(x, y, x+w, y+h)
}

func (g *Game) helpCloseRect() image.Rectangle {
	size := g.iconSize()
	r := g.helpMenuRect()
	return image.Rect(r.Max.X-size-uiScaled(2), r.Min.Y+uiScaled(2), r.Max.X-uiScaled(2), r.Min.Y+size+uiScaled(2))
}

func (g *Game) geyserRect() image.Rectangle {
	size := g.iconSize()
	x := g.width - size*3 - uiScaled(HelpMargin*3)
	y := g.height - size - uiScaled(HelpMargin)
	return image.Rect(x, y, x+size, y+size)
}

func (g *Game) geyserCloseRect() image.Rectangle {
	size := g.iconSize()
	x := g.width - size - uiScaled(HelpMargin)
	y := uiScaled(HelpMargin)
	return image.Rect(x, y, x+size, y+size)
}

func (g *Game) bottomTrayRect() image.Rectangle {
	r := g.optionsRect()
	r = r.Union(g.geyserRect())
	r = r.Union(g.screenshotRect())
	r = r.Union(g.helpRect())
	m := uiScaled(4)
	return image.Rect(r.Min.X-m, r.Min.Y-m, r.Max.X+m, r.Max.Y+m)
}

func (g *Game) biomeLegendRect() image.Rectangle {
	if g.legend == nil {
		return image.Rectangle{}
	}
	w := g.legend.Bounds().Dx()
	h := g.legend.Bounds().Dy()
	y := -int(g.biomeScroll)
	return image.Rect(0, y, w, y+h)
}

func (g *Game) itemLegendRect() image.Rectangle {
	if g.legendImage == nil {
		return image.Rectangle{}
	}
	w := g.legendImage.Bounds().Dx()
	h := g.legendImage.Bounds().Dy()
	x := g.width - w - uiScaled(12)
	y := uiScaled(10) - int(g.itemScroll)
	return image.Rect(x, y, x+w, y+h)
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
				scale := g.zoom * IconScale * g.iconScale * float64(BaseIconPixels) / maxDim
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
				scale := g.zoom * IconScale * g.iconScale * float64(BaseIconPixels) / maxDim
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
	for _, name := range names {
		img, _ := loadImageFile(name)
		if g.icons != nil {
			g.icons[name] = img
		}
	}
	g.needsRedraw = true
}
