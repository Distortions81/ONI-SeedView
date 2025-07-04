// Package constants for configuration.
package main

import (
	"math"
	"time"
)

const (
	BaseURL          = "https://ingest.mapsnotincluded.org/coordinate/"
	AcceptCBORHeader = "application/cbor"
	PanSpeed         = 15
	// CameraMargin controls how far the world can be panned
	// beyond the visible screen in pixels.
	CameraMargin        = -64
	KeyZoomFactor       = 1.05
	WheelZoomFactor     = 1.1
	MinZoom             = 0.5
	IconScale           = 0.2
	LegendZoomExponent  = 10
	DefaultWidth        = 1280
	DefaultHeight       = 720
	InitialZoom         = 1.0
	LabelCharWidth      = 6
	NumberLegendXOffset = 150
	LegendRowSpacing    = 25
	ScreenshotFile      = "screenshot.png"
	HelpIconSize        = 24
	HelpMargin          = 10
	CrosshairSize       = 10
	InfoIconScale       = 0.3
	InfoPanelAlpha      = 200
	TouchDragThreshold  = 10
	// WheelThrottle controls how often mouse wheel zoom is applied
	// in WASM to account for faster scroll events.
	WheelThrottle = 75 * time.Millisecond

	// WebAssetBase specifies the path used to fetch image assets
	// when running in WebAssembly. The assets should be served
	// relative to the page URL.
	WebAssetBase = "assets/"
)

var LegendZoomThreshold = math.Pow(WheelZoomFactor, LegendZoomExponent)
