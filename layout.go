//go:build !test

package main

import (
	"math"
)

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	screenW := outsideWidth
	screenH := outsideHeight

	/*
		if !g.screenshotMode {
				minDim := screenW
				if screenH < minDim {
					minDim = screenH
				}

					newSize := baseFontSize
					if minDim < 800 {
						newSize = baseFontSize - 2
					} else if minDim > 1200 {
						newSize = baseFontSize + 2
					}
					if fontSize != newSize {
						setFontSize(newSize)
					}
		}
	*/

	if g.width != screenW || g.height != screenH {
		// Keep the world position at the center of the screen fixed so
		// resizing doesn't shift the view.
		cxOld, cyOld := float64(g.width)/2, float64(g.height)/2
		worldX := (cxOld - g.camX) / g.zoom
		worldY := (cyOld - g.camY) / g.zoom

		g.width = screenW
		g.height = screenH

		cxNew, cyNew := float64(g.width)/2, float64(g.height)/2
		g.camX = cxNew - worldX*g.zoom
		g.camY = cyNew - worldY*g.zoom
		g.clampCamera()
		if g.astWidth > 0 && g.astHeight > 0 {
			zx := float64(g.width) / (float64(g.astWidth) * 2)
			zy := float64(g.height) / (float64(g.astHeight) * 2)
			g.minZoom = math.Min(zx, zy) * 0.25
		}

		if max := g.maxBiomeScroll(); max == 0 {
			g.biomeScroll = 0
		} else if g.biomeScroll > max {
			g.biomeScroll = max
		}
		if max := g.maxItemScroll(); max == 0 {
			g.itemScroll = 0
		} else if g.itemScroll > max {
			g.itemScroll = max
		}

		g.needsRedraw = true
	}
	return screenW, screenH
}
