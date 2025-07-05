//go:build !test

package main

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

func (g *Game) optionsRect() image.Rectangle {
	size := g.iconSize()
	x := g.width - size*5 - HelpMargin*5
	y := g.height - size - HelpMargin
	return image.Rect(x, y, x+size, y+size)
}

func (g *Game) optionsMenuRect() image.Rectangle {
	labels := []string{
		OptionsMenuTitle,
		"Textures",
		"Vsync",
		"Touch Controls",
		"Show Item Names",
		"Show Legends",
		"Use Item Numbers",
		"Icon Size +",
		"Icon Size -",
		"Smart Rendering",
		"Version: " + ClientVersion,
	}
	maxW := 0
	for _, s := range labels {
		if len(s) > maxW {
			maxW = len(s)
		}
	}
	w := maxW*LabelCharWidth + 4
	h := (len(labels)+1)*OptionsMenuSpacing + 4
	x := g.optionsRect().Min.X - w - 10
	if x < 0 {
		x = 0
	}
	y := g.optionsRect().Min.Y - h
	if y < 0 {
		y = 0
	}
	return image.Rect(x, y, x+w, y+h)
}

func (g *Game) drawOptionsMenu(dst *ebiten.Image) {
	rect := g.optionsMenuRect()
	vector.DrawFilledRect(dst, float32(rect.Min.X), float32(rect.Min.Y), float32(rect.Dx()), float32(rect.Dy()), colorRGBA(0, 0, 0, 200), false)
	drawTextWithBG(dst, OptionsMenuTitle, rect.Min.X+2, rect.Min.Y+2)
	y := rect.Min.Y + 2 + OptionsMenuSpacing
	items := []struct {
		label   string
		enabled bool
	}{
		{"Textures", g.textures},
		{"Vsync", g.vsync},
		{"Touch Controls", g.mobile},
		{"Show Item Names", g.showItemNames},
		{"Show Legends", g.showLegend},
		{"Use Item Numbers", g.useNumbers},
		{"Icon Size +", false},
		{"Icon Size -", false},
		{"Smart Rendering", g.smartRender},
	}
	for _, it := range items {
		if it.enabled {
			drawTextWithBGBorder(dst, it.label, rect.Min.X+2, y, color.RGBA{255, 255, 255, 255})
		} else {
			drawTextWithBG(dst, it.label, rect.Min.X+2, y)
		}
		y += OptionsMenuSpacing
	}
	drawTextWithBG(dst, "Version: "+ClientVersion, rect.Min.X+2, y)
}

func (g *Game) clickOptionsMenu(mx, my int) bool {
	rect := g.optionsMenuRect()
	if !rect.Overlaps(image.Rect(mx, my, mx+1, my+1)) {
		return false
	}
	y := rect.Min.Y + 2 + OptionsMenuSpacing
	for i := 0; i < 9; i++ {
		r := image.Rect(rect.Min.X, y-2, rect.Min.X+rect.Dx(), y-2+16+4)
		if r.Overlaps(image.Rect(mx, my, mx+1, my+1)) {
			switch i {
			case 0:
				g.textures = !g.textures
			case 1:
				g.vsync = !g.vsync
				ebiten.SetVsyncEnabled(g.vsync)
			case 2:
				g.mobile = !g.mobile
			case 3:
				g.showItemNames = !g.showItemNames
			case 4:
				g.showLegend = !g.showLegend
			case 5:
				g.useNumbers = !g.useNumbers
			case 6:
				g.iconScale += 0.25
			case 7:
				if g.iconScale > 0.25 {
					g.iconScale -= 0.25
				}
			case 8:
				g.smartRender = !g.smartRender
			}
			g.needsRedraw = true
			return true
		}
		y += OptionsMenuSpacing
	}
	return true
}
