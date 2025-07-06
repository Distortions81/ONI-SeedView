package main

import "math"

var uiScale = 1.0

func setUIScale(scale float64) {
	if scale < 0.5 {
		scale = 0.5
	}
	uiScale = scale
	setFontSize(baseFontSize * uiScale)
}

func uiScaled(v int) int {
	return int(math.Round(float64(v) * uiScale))
}

func uiScaledF(v float64) float64 {
	return v * uiScale
}

func (g *Game) uiScaleFactor() float64 { return uiScale }
