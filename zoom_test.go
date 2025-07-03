package main_test

import (
	"image"
	"math"
	"testing"
)

// transform simulates the coordinate math used when drawing polygons.
func transform(pts []Point, camX, camY, zoom float64) []image.Point {
	out := make([]image.Point, len(pts))
	for i, p := range pts {
		x := float64(p.X*2)*zoom + camX
		y := float64(p.Y*2)*zoom + camY
		out[i] = image.Point{int(math.Round(x)), int(math.Round(y))}
	}
	return out
}

func TestZoomCoordinateTransform(t *testing.T) {
	pts := []Point{{0, 0}, {1, 0}, {1, 1}, {0, 1}}
	got := transform(pts, 10, 20, 1.5)
	want := []image.Point{{10, 20}, {13, 20}, {13, 23}, {10, 23}}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("vertex %d: want %v got %v", i, want[i], got[i])
		}
	}
}
