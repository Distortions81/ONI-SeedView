package main

import "math"

// uiScale represents the user configured scaling factor. The actual scale used
// for drawing also includes the HiDPI device scale factor.
var uiScale = 1.0

func setUIScale(scale float64) {
	if scale < 0.5 {
		scale = 0.5
	}
	uiScale = scale
	setFontSize(baseFontSize * uiScale)
}

func uiScaled(v int) int {
	return int(math.Round(float64(v) * uiScale * getHiDPIScale()))
}

func uiScaledF(v float64) float64 {
	return v * uiScale * getHiDPIScale()
}

// uiScaleFactor returns the combined UI scale including the HiDPI factor.
func (g *Game) uiScaleFactor() float64 { return uiScale * getHiDPIScale() }
