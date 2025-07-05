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
	"Sandstone":           darkenColor(colorFromARGB(0xFFF2BB47), 0.6),
	"Barren":              darkenColor(colorFromARGB(0xFF97752C), 0.6),
	"Space":               darkenColor(colorFromARGB(0xFFD0D0D0), 0.6),
	"FrozenWastes":        darkenColor(colorFromARGB(0xFF83cdde), 0.6),
	"BoggyMarsh":          darkenColor(colorFromARGB(0xFF7B974B), 0.6),
	"ToxicJungle":         darkenColor(colorFromARGB(0xFFCB95A3), 0.6),
	"Ocean":               darkenColor(colorFromARGB(0xFF4C4CFF), 0.6),
	"Rust":                darkenColor(colorFromARGB(0xFFFFA007), 0.6),
	"Forest":              darkenColor(colorFromARGB(0xFF8EC039), 0.6),
	"Radioactive":         darkenColor(colorFromARGB(0xFF4AE458), 0.6),
	"Swamp":               darkenColor(colorFromARGB(0xFFEB9B3F), 0.6),
	"Wasteland":           darkenColor(colorFromARGB(0xFFCC3636), 0.6),
	"Metallic":            darkenColor(colorFromARGB(0xFFFFA007), 0.6),
	"Moo":                 darkenColor(colorFromARGB(0xFF8EC039), 0.6),
	"IceCaves":            darkenColor(colorFromARGB(0xFF6C9BD3), 0.6),
	"CarrotQuarry":        darkenColor(colorFromARGB(0xFFCDA2C7), 0.6),
	"SugarWoods":          darkenColor(colorFromARGB(0xFFA2CDA4), 0.6),
	"PrehistoricGarden":   darkenColor(colorFromARGB(0xFF006127), 0.6),
	"PrehistoricRaptor":   darkenColor(colorFromARGB(0xFF352F8C), 0.6),
	"PrehistoricWetlands": darkenColor(colorFromARGB(0xFF645906), 0.6),
	"OilField":            darkenColor(colorFromARGB(0xFF52321D), 0.6),
	"MagmaCore":           darkenColor(colorFromARGB(0xFFDE5A3B), 0.6),
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
