package main

import (
	"image"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

func (g *Game) asteroidArrowRect() image.Rectangle {
	name := g.asteroidID
	if name == "" {
		name = "Unknown"
	}
	name = truncateString(name, 32)
	text := "Asteroid: " + name
	w, _ := textDimensions(text)
	x := g.width/2 + w/2 + uiScaled(4)
	size := uiScaled(12)
	baseline := seedBaseline() + notoFont.Metrics().Height.Ceil() + uiScaled(4)
	y := baseline + notoFont.Metrics().Descent.Ceil() + uiScaled(4) - size/2
	return image.Rect(x, y, x+size, y+size)
}

// asteroidInfoRect returns the bounding rectangle surrounding the seed
// coordinate, asteroid name, and arrow. This is used as the clickable area for
// opening the asteroid menu.
func (g *Game) asteroidInfoRect() image.Rectangle {
	if g.coord == "" {
		return image.Rectangle{}
	}
	sw, sh := textDimensions(g.coord)
	sx := g.width/2 - sw/2
	seedRect := image.Rect(sx-uiScaled(2), seedBaseline()-uiScaled(2), sx+sw+uiScaled(2), seedBaseline()-uiScaled(2)+sh+uiScaled(4))

	name := g.asteroidID
	if name == "" {
		name = "Unknown"
	}
	name = truncateString(name, 32)
	astText := "Asteroid: " + name
	aw, ah := textDimensions(astText)
	ax := g.width/2 - aw/2
	astBase := seedBaseline() + notoFont.Metrics().Height.Ceil() + uiScaled(4)
	astRect := image.Rect(ax-uiScaled(2), astBase-uiScaled(2), ax+aw+uiScaled(2), astBase-uiScaled(2)+ah+uiScaled(4))

	rect := seedRect.Union(astRect)
	rect = rect.Union(g.asteroidArrowRect())
	return image.Rect(rect.Min.X-uiScaled(2), rect.Min.Y-uiScaled(2), rect.Max.X+uiScaled(2), rect.Max.Y+uiScaled(2))
}

func drawDownArrow(dst *ebiten.Image, rect image.Rectangle, up bool) {

	cx := float32(rect.Min.X + rect.Dx()/2)
	cy := float32(rect.Min.Y + rect.Dy()/2)
	half := float32(rect.Dx()) / 2
	var p vector.Path
	if up {
		p.MoveTo(cx-half, cy+half/2)
		p.LineTo(cx+half, cy+half/2)
		p.LineTo(cx, cy-half/2)
	} else {
		p.MoveTo(cx-half, cy-half/2)
		p.LineTo(cx+half, cy-half/2)
		p.LineTo(cx, cy+half/2)
	}
	p.Close()
	vs, is := p.AppendVerticesAndIndicesForFilling(nil, nil)
	for i := range vs {
		vs[i].ColorR = 1
		vs[i].ColorG = 1
		vs[i].ColorB = 1
		vs[i].ColorA = 1
	}
	op := &ebiten.DrawTrianglesOptions{AntiAlias: true, ColorScaleMode: ebiten.ColorScaleModePremultipliedAlpha}
	dst.DrawTriangles(vs, is, whitePixel, op)
}

func drawCheck(dst *ebiten.Image, rect image.Rectangle) {
	thickness := float32(rect.Dx()) / 6
	x0 := float32(rect.Min.X) + float32(rect.Dx())*0.2
	y0 := float32(rect.Min.Y) + float32(rect.Dy())*0.5
	x1 := float32(rect.Min.X) + float32(rect.Dx())*0.45
	y1 := float32(rect.Min.Y) + float32(rect.Dy())*0.8
	x2 := float32(rect.Min.X) + float32(rect.Dx())*0.8
	y2 := float32(rect.Min.Y) + float32(rect.Dy())*0.2
	vector.StrokeLine(dst, x0, y0, x1, y1, thickness, color.White, true)
	vector.StrokeLine(dst, x1, y1, x2, y2, thickness, color.White, true)
}

func (g *Game) asteroidMenuSize() (int, int) {
	maxW, _ := textDimensions(AsteroidMenuTitle)
	for _, a := range g.asteroids {
		name := truncateString(a.ID, 64)
		w, _ := textDimensions(name)
		if w > maxW {
			maxW = w
		}
	}
	// Add extra space for the checkmark and some padding so
	// longer names don't butt up against the right edge of the menu.
	// Include an extra character width of padding for clarity.
	w := maxW + uiScaled(28) + LabelCharWidth
	h := (len(g.asteroids)+1)*menuSpacing() + uiScaled(4)
	return w, h
}

func (g *Game) asteroidMenuRect() image.Rectangle {
	w, h := g.asteroidMenuSize()
	ar := g.asteroidArrowRect()
	x := ar.Min.X + ar.Dx()/2 - w/2
	if x < 0 {
		x = 0
	}
	if x+w > g.width {
		x = g.width - w
	}
	y := ar.Max.Y + uiScaled(4)
	if y+h > g.height {
		y = g.height - h
	}
	return image.Rect(x, y, x+w, y+h)
}

func (g *Game) maxAsteroidScroll() float64 {
	_, h := g.asteroidMenuSize()
	max := float64(h) - float64(g.height)
	if max < 0 {
		return 0
	}
	return max
}

func (g *Game) adjustAsteroidScroll(d float64) {
	g.asteroidScroll += d
	if g.asteroidScroll < 0 {
		g.asteroidScroll = 0
	}
	if m := g.maxAsteroidScroll(); g.asteroidScroll > m {
		g.asteroidScroll = m
	}
	g.needsRedraw = true
}

func (g *Game) drawAsteroidMenu(dst *ebiten.Image) {
	rect := g.asteroidMenuRect()
	w, h := g.asteroidMenuSize()
	img := ebiten.NewImage(w, h)
	drawFrame(img, image.Rect(0, 0, w, h))
	pad := uiScaled(6)
	drawText(img, AsteroidMenuTitle, pad, pad, false)
	y := pad + menuSpacing() - int(g.asteroidScroll)
	for _, a := range g.asteroids {
		btn := image.Rect(uiScaled(4), y-uiScaled(4), w-uiScaled(4), y-uiScaled(4)+menuButtonHeight())
		drawButton(img, btn, a.ID == g.asteroidID)
		if a.ID == g.asteroidID {
			ck := image.Rect(btn.Min.X+uiScaled(4), btn.Min.Y+uiScaled(4), btn.Min.X+uiScaled(16), btn.Min.Y+uiScaled(16))
			drawCheck(img, ck)
		}
		name := truncateString(a.ID, 64)
		drawText(img, name, btn.Min.X+uiScaled(20), btn.Min.Y+uiScaled(4), false)
		y += menuSpacing()
	}
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(rect.Min.X), float64(rect.Min.Y))
	dst.DrawImage(img, op)
}

func (g *Game) clickAsteroidMenu(mx, my int) bool {
	rect := g.asteroidMenuRect()
	if !rect.Overlaps(image.Rect(mx, my, mx+1, my+1)) {
		return false
	}
	x := mx - rect.Min.X
	y := my - rect.Min.Y + int(g.asteroidScroll)
	mx = x
	my = y
	w, _ := g.asteroidMenuSize()
	yPos := uiScaled(6) + menuSpacing()
	for i, a := range g.asteroids {
		r := image.Rect(uiScaled(4), yPos-uiScaled(4), w-uiScaled(4), yPos-uiScaled(4)+menuButtonHeight())
		if r.Overlaps(image.Rect(mx, my, mx+1, my+1)) {
			g.showAstMenu = false
			g.asteroidScroll = 0
			if a.ID != g.asteroidID {
				g.loadAsteroid(a)
			}
			g.needsRedraw = true
			_ = i
			return true
		}
		yPos += menuSpacing()
	}
	return true
}

func (g *Game) loadAsteroid(ast Asteroid) {
	g.invalidateLegends()
	g.asteroidID = ast.ID
	// Reset cached legends so item numbering matches the newly loaded
	// asteroid contents. These are rebuilt lazily when needed.
	g.legendMap = nil
	g.legendEntries = nil
	g.legendColors = nil
	g.selectedItem = -1
	g.itemScroll = 0
	bps := parseBiomePaths(ast.BiomePaths)
	g.geysers = ast.Geysers
	g.pois = ast.POIs
	g.biomes = bps
	g.astWidth = ast.SizeX
	g.astHeight = ast.SizeY
	g.legend, g.legendBiomes = buildLegendImage(bps)
	g.fitOnLoad = true
	g.biomeTextures = loadBiomeTextures()
	names := []string{"../icons/camera.png", "../icons/help.png", "../icons/gear.png", "geyser_water.png"}
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
	g.startIconLoader(names)
}

func (g *Game) centerAndFit() {
	if g.astWidth == 0 || g.astHeight == 0 {
		return
	}
	zoomX := float64(g.width) / (float64(g.astWidth) * 2)
	zoomY := float64(g.height) / (float64(g.astHeight) * 2)
	g.zoom = math.Min(zoomX, zoomY)
	g.minZoom = g.zoom * 0.25
	g.camX = (float64(g.width) - float64(g.astWidth)*2*g.zoom) / 2
	g.camY = (float64(g.height) - float64(g.astHeight)*2*g.zoom) / 2
	g.clampCamera()
}
