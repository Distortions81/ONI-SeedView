//go:build !test

package main

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

func colorFromARGB(hex uint32) color.RGBA {
	return color.RGBA{
		R: uint8(hex >> 16),
		G: uint8(hex >> 8),
		B: uint8(hex),
		A: uint8(hex >> 24),
	}
}

func darkenColor(c color.RGBA, factor float64) color.RGBA {
	return color.RGBA{
		R: uint8(float64(c.R) * factor),
		G: uint8(float64(c.G) * factor),
		B: uint8(float64(c.B) * factor),
		A: c.A,
	}
}

var biomeColors = map[string]color.RGBA{
	"Sandstone":           {R: 204, G: 179, B: 61, A: 255},
	"Barren":              {R: 204, G: 154, B: 61, A: 255},
	"Space":               {R: 204, G: 204, B: 204, A: 255},
	"FrozenWastes":        {R: 61, G: 183, B: 204, A: 255},
	"BoggyMarsh":          {R: 138, G: 204, B: 61, A: 255},
	"ToxicJungle":         {R: 204, G: 92, B: 136, A: 255},
	"Ocean":               {R: 61, G: 61, B: 204, A: 255},
	"Rust":                {R: 204, G: 118, B: 61, A: 255},
	"Forest":              {R: 98, G: 204, B: 61, A: 255},
	"Radioactive":         {R: 61, G: 204, B: 90, A: 255},
	"Swamp":               {R: 204, G: 145, B: 61, A: 255},
	"Wasteland":           {R: 204, G: 61, B: 61, A: 255},
	"Metallic":            {R: 204, G: 118, B: 61, A: 255},
	"Moo":                 {R: 98, G: 204, B: 61, A: 255},
	"IceCaves":            {R: 61, G: 134, B: 204, A: 255},
	"CarrotQuarry":        {R: 204, G: 92, B: 147, A: 255},
	"SugarWoods":          {R: 92, G: 204, B: 106, A: 255},
	"PrehistoricGarden":   {R: 61, G: 204, B: 90, A: 255},
	"PrehistoricRaptor":   {R: 61, G: 61, B: 204, A: 255},
	"PrehistoricWetlands": {R: 154, G: 134, B: 61, A: 255},
	"OilField":            {R: 111, G: 78, B: 55, A: 255},
	"MagmaCore":           {R: 204, G: 90, B: 61, A: 255},
}

func init() {
	for c, _ := range biomeColors {
		biomeColors[c] = darkenColor(biomeColors[c], 0.8)
	}
}

var whitePixel = func() *ebiten.Image {
	img := ebiten.NewImage(1, 1)
	img.Fill(color.White)
	return img
}()

// uniqueColor generates a visually distinct color for the given index.
func uniqueColor(index int) color.RGBA {
	h := float64((index * 137) % 360)
	period := 30.0
	s := 0.5 + 0.4*math.Sin(float64(index)*2*math.Pi/period)
	l := 0.45 + 0.3*math.Cos(float64(index)*2*math.Pi/period)
	if s < 0.3 {
		s = 0.3
	} else if s > 0.9 {
		s = 0.9
	}
	if l < 0.3 {
		l = 0.3
	} else if l > 0.8 {
		l = 0.8
	}
	c := (1 - math.Abs(2*l-1)) * s
	x := c * (1 - math.Abs(math.Mod(h/60.0, 2)-1))
	m := l - c/2
	var r1, g1, b1 float64
	switch {
	case h < 60:
		r1, g1, b1 = c, x, 0
	case h < 120:
		r1, g1, b1 = x, c, 0
	case h < 180:
		r1, g1, b1 = 0, c, x
	case h < 240:
		r1, g1, b1 = 0, x, c
	case h < 300:
		r1, g1, b1 = x, 0, c
	default:
		r1, g1, b1 = c, 0, x
	}
	r := uint8(math.Round((r1 + m) * 255))
	g := uint8(math.Round((g1 + m) * 255))
	b := uint8(math.Round((b1 + m) * 255))
	return color.RGBA{r, g, b, 255}
}
