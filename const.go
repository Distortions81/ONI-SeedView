// Package constants for configuration.
package main

import (
	"image/color"
	"math"
	"time"
)

const (
	ClientVersion    = "v0.0.5-2507060137"
	BaseURL          = "https://ingest.mapsnotincluded.org/coordinate/"
	AcceptCBORHeader = "application/cbor"
	PanSpeed         = 15
	// CameraMargin controls how far the world can be panned
	// beyond the visible screen in pixels.
	CameraMargin    = -64
	KeyZoomFactor   = 1.05
	WheelZoomFactor = 1.1
	MinZoom         = 0.25
	MaxZoom         = 8.0
	IconScale       = 0.2
	// BaseIconPixels is the target pixel size used when scaling
	// geyser and POI icons. Larger icons are scaled down so the
	// on-screen size is consistent across all items.
	BaseIconPixels      = 96
	LegendZoomExponent  = 10
	DefaultWidth        = 1200
	DefaultHeight       = 1200
	InitialZoom         = 1.0
	LabelCharWidth      = 6
	NumberLegendXOffset = 150
	LegendRowSpacing    = 25
	// LegendScrollExtraRows controls how many additional rows the legend
	// panels can scroll beyond the bottom of the screen.
	LegendScrollExtraRows = 3
	ScreenshotFile        = "screenshot.png"
	HelpIconSize          = 24
	HelpMargin            = 10
	CrosshairSize         = 10
	InfoIconScale         = 0.3
	InfoIconSize          = 32
	InfoPanelAlpha        = 200
	TouchDragThreshold    = 20
	ScreenshotMenuSpacing = 26
	ScreenshotMenuTitle   = "Image quality:"
	ScreenshotSaveLabel   = "Save Screenshot"
	ScreenshotTakingLabel = "Taking Screenshot..."
	ScreenshotSavedLabel  = "Saved!"
	ScreenshotBWLabel     = "Black and White"
	ScreenshotCloseLabel  = "Close"
	GeyserRowSpacing      = 60
	OptionsMenuSpacing    = 26
	OptionsMenuTitle      = "Options:"
	AsteroidMenuSpacing   = 26
	AsteroidMenuTitle     = "Asteroids:"

	// BiomeTextureScale controls the repetition of biome textures.
	// Smaller values result in more repetitions.
	BiomeTextureScale = 0.1

	// WheelThrottle controls how often mouse wheel zoom is applied
	// in WASM to account for faster scroll events.
	WheelThrottle = 75 * time.Millisecond

	// WebAssetBase specifies the path used to fetch image assets
	// when running in WebAssembly. The assets should be served
	// relative to the page URL.
	WebAssetBase = "assets/"
)

var LegendZoomThreshold = math.Pow(WheelZoomFactor, LegendZoomExponent)

var ScreenshotQualities = []string{"Low (3K)", "Medium (6K)", "High (8K)"}

var ScreenshotScales = []float64{4.0, 8.0, 10.0}

var biomeOrder = []string{
	"Sandstone",
	"Barren",
	"Space",
	"FrozenWastes",
	"BoggyMarsh",
	"ToxicJungle",
	"Ocean",
	"Rust",
	"Forest",
	"Radioactive",
	"Swamp",
	"Wasteland",
	"Metallic",
	"Moo",
	"IceCaves",
	"CarrotQuarry",
	"SugarWoods",
	"PrehistoricGarden",
	"PrehistoricRaptor",
	"PrehistoricWetlands",
	"OilField",
	"MagmaCore",
}

var (
	buttonActiveColor   = color.RGBA{0, 96, 96, 255}
	buttonInactiveColor = color.RGBA{40, 40, 40, 255}
	buttonBorderColor   = color.RGBA{255, 255, 255, 255}
	frameColor          = color.RGBA{0, 0, 0, 255}
	bottomTrayColor     = color.RGBA{0, 0, 0, 180}
)
