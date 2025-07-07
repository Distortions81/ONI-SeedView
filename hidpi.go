package main

import "github.com/hajimehoshi/ebiten/v2"

var hiDPIEnabled = true
var hiDPIScale = 1.0

func getHiDPIScale() float64 {
	if !hiDPIEnabled {
		return 1.0
	}
	if hiDPIScale == 1.0 {
		s := ebiten.Monitor().DeviceScaleFactor()
		if s > 0 {
			hiDPIScale = s
		}
	}
	return hiDPIScale
}

func setHiDPI(enabled bool) {
	hiDPIEnabled = enabled
	hiDPIScale = 1.0
	if enabled {
		if s := ebiten.Monitor().DeviceScaleFactor(); s > 0 {
			hiDPIScale = s
		}
	}
	// Reload fonts with new DPI
	setFontSize(fontSize)
}
