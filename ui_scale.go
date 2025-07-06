package main

import "math"

var (
	uiBaseScale = 1.0
	uiUserScale = 1.0
	uiScale     = 1.0
)

func updateUIScale() {
	total := uiBaseScale * uiUserScale
	if total < 0.5 {
		total = 0.5
	}
	if math.Abs(total-uiScale) < 0.01 {
		return
	}
	uiScale = total
	setFontSize(baseFontSize * uiScale)
}

func setUIScale(scale float64) {
	if scale < 0.5 {
		scale = 0.5
	}
	uiUserScale = scale
	updateUIScale()
}

func setUIBaseScale(scale float64) {
	uiBaseScale = scale
	updateUIScale()
}

func uiScaled(v int) int {
	return int(math.Round(float64(v) * uiScale))
}

func uiScaledF(v float64) float64 {
	return v * uiScale
}

func (g *Game) uiScaleFactor() float64 { return uiScale }
