package main

import (
	"fmt"
	"image/color"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
)

var helpMessage string

type hoverIcon int

const (
	hoverNone hoverIcon = iota
	hoverOptions
	hoverGeysers
	hoverScreenshot
	hoverHelp
)

var errorBorderColor = color.RGBA{R: 244, G: 67, B: 54, A: 255}

var grayColorM = func() ebiten.ColorM {
	var m ebiten.ColorM
	const r = 0.299
	const g = 0.587
	const b = 0.114
	m.SetElement(0, 0, r)
	m.SetElement(0, 1, g)
	m.SetElement(0, 2, b)
	m.SetElement(1, 0, r)
	m.SetElement(1, 1, g)
	m.SetElement(1, 2, b)
	m.SetElement(2, 0, r)
	m.SetElement(2, 1, g)
	m.SetElement(2, 2, b)
	m.SetElement(3, 3, 1)
	return m
}()

func init() {
	lines := [][2]string{
		{"Arrow keys/WASD", "pan the camera"},
		{"Mouse wheel or +/-", "zoom in and out"},
		{"Drag with the mouse/touch", "pan"},
		{"Pinch with two fingers", "zoom on touch"},
		{"Click or tap geysers/POIs", "center and show details"},
		{"Tap legend entries", "highlight items"},
		{"Camera icon", "open screenshot menu"},
		{"Geyser-icon", "list all geysers"},
		{"Question mark", "toggle this help"},
		{"X button", "close this help"},
		{"Gear icon", "open options"},
	}
	width := 0
	for _, p := range lines {
		if len(p[0]) > width {
			width = len(p[0])
		}
	}
	var b strings.Builder
	b.WriteString("Controls:\n")
	for i, p := range lines {
		fmt.Fprintf(&b, "%-*s | %s", width, p[0], p[1])
		if i < len(lines)-1 {
			b.WriteByte('\n')
		}
	}
	helpMessage = b.String()
}
