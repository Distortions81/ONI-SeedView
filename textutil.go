//go:build !test

package main

import (
	"strings"

	"github.com/hajimehoshi/ebiten/v2/text"
)

func textDimensions(str string) (int, int) {
	lines := strings.Split(str, "\n")
	if notoFont == nil {
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

	width := 0
	for _, l := range lines {
		b := text.BoundString(notoFont, l)
		if b.Dx() > width {
			width = b.Dx()
		}
	}
	lineHeight := notoFont.Metrics().Height.Ceil()
	return width, lineHeight * len(lines)
}
