//go:build test

package main

import "strings"

func textDimensions(str string) (int, int) {
	lines := strings.Split(str, "\n")
	width := 0
	for _, l := range lines {
		if len(l) > width {
			width = len(l)
		}
	}
	scale := fontScale()
	w := int(float64(width*LabelCharWidth) * scale)
	h := int(float64(len(lines)*20) * scale)
	return w, h
}
