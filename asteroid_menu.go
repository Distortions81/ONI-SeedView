//go:build !test

package main

import (
	"image"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

func (g *Game) asteroidArrowRect() image.Rectangle {
	scale := g.uiScale()
	name := g.asteroidID
	if name == "" {
		name = "Unknown"
	}
	text := "Asteroid: " + name
	w, _ := textDimensions(text)
	x := g.width/2 + int(float64(w)*scale/2) + int(4*scale)
	size := int(12 * scale)
	y := int(30*scale + 8*scale - float64(size)/2)
	return image.Rect(x, y, x+size, y+size)
}

func drawDownArrow(dst *ebiten.Image, rect image.Rectangle, up bool) {
	drawFrame(dst, rect)

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
	vector.StrokeLine(dst, x0, y0, x1, y1, thickness, colorWhite, true)
	vector.StrokeLine(dst, x1, y1, x2, y2, thickness, colorWhite, true)
}

func (g *Game) asteroidMenuSize() (int, int) {
	maxW, _ := textDimensions(AsteroidMenuTitle)
	for _, a := range g.asteroids {
		w, _ := textDimensions(a.ID)
		if w > maxW {
			maxW = w
		}
	}
	w := maxW + 24
	h := (len(g.asteroids)+1)*AsteroidMenuSpacing + 4
	return w, h
}

func (g *Game) asteroidMenuRect() image.Rectangle {
	w, h := g.asteroidMenuSize()
	scale := g.uiScale()
	w = int(float64(w) * scale)
	h = int(float64(h) * scale)
	ar := g.asteroidArrowRect()
	x := ar.Min.X + ar.Dx()/2 - w/2
	if x < 0 {
		x = 0
	}
	if x+w > g.width {
		x = g.width - w
	}
	y := ar.Max.Y + int(4*scale)
	if y+h > g.height {
		y = g.height - h
	}
	return image.Rect(x, y, x+w, y+h)
}

func (g *Game) maxAsteroidScroll() float64 {
	_, h := g.asteroidMenuSize()
	scale := g.uiScale()
	max := float64(h)*scale - float64(g.height)
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
	scale := g.uiScale()
	rect := g.asteroidMenuRect()
	w, h := g.asteroidMenuSize()
	img := ebiten.NewImage(w, h)
	drawFrame(img, image.Rect(0, 0, w, h))
	drawText(img, AsteroidMenuTitle, 6, 6)
	y := 6 + AsteroidMenuSpacing - int(g.asteroidScroll)
	for _, a := range g.asteroids {
		btn := image.Rect(4, y-4, w-4, y-4+22)
		drawButton(img, btn, a.ID == g.asteroidID)
		if a.ID == g.asteroidID {
			ck := image.Rect(btn.Min.X+4, btn.Min.Y+4, btn.Min.X+16, btn.Min.Y+16)
			drawCheck(img, ck)
		}
		drawText(img, a.ID, btn.Min.X+20, btn.Min.Y+4)
		y += AsteroidMenuSpacing
	}
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(scale, scale)
	op.GeoM.Translate(float64(rect.Min.X), float64(rect.Min.Y))
	dst.DrawImage(img, op)
}

func (g *Game) clickAsteroidMenu(mx, my int) bool {
	rect := g.asteroidMenuRect()
	if !rect.Overlaps(image.Rect(mx, my, mx+1, my+1)) {
		return false
	}
	scale := g.uiScale()
	x := int(float64(mx-rect.Min.X) / scale)
	y := int(float64(my-rect.Min.Y)/scale) + int(g.asteroidScroll)
	mx = x
	my = y
	w, _ := g.asteroidMenuSize()
	yPos := 6 + AsteroidMenuSpacing
	for i, a := range g.asteroids {
		r := image.Rect(4, yPos-4, w-4, yPos-4+22)
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
		yPos += AsteroidMenuSpacing
	}
	return true
}

func (g *Game) loadAsteroid(ast Asteroid) {
	g.asteroidID = ast.ID
	bps := parseBiomePaths(ast.BiomePaths)
	g.geysers = ast.Geysers
	g.pois = ast.POIs
	g.biomes = bps
	g.astWidth = ast.SizeX
	g.astHeight = ast.SizeY
	g.legend, g.legendBiomes = buildLegendImage(bps)
	g.centerAndFit()
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

var colorWhite = color.RGBA{255, 255, 255, 255}
