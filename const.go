// Package constants for configuration.
package main

import "math"

const (
	BaseURL             = "https://ingest.mapsnotincluded.org/coordinate/"
	AcceptCBORHeader    = "application/cbor"
	PanSpeed            = 5
	KeyZoomFactor       = 1.05
	WheelZoomFactor     = 1.1
	MinZoom             = 0.1
	IconScale           = 0.25
	LegendZoomExponent  = 10
	DefaultWidth        = 600
	DefaultHeight       = 800
	InitialZoom         = 1.0
	LabelCharWidth      = 6
	NumberLegendXOffset = 150
	LegendRowSpacing    = 18
	ScreenshotFile      = "screenshot.png"
)

var LegendZoomThreshold = math.Pow(WheelZoomFactor, LegendZoomExponent)
