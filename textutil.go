package main

import (
	"strings"

	"golang.org/x/image/font"
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
		b, _ := font.BoundString(notoFont, l)
		w := (b.Max.X - b.Min.X).Ceil()
		if w > width {
			width = w
		}
	}
	lineHeight := notoFont.Metrics().Height.Ceil()
	return width, lineHeight * len(lines)
}

func truncateString(s string, max int) string {
	if len([]rune(s)) <= max {
		return s
	}
	runes := []rune(s)
	if max <= 3 {
		return string(runes[:max])
	}
	return string(runes[:max-3]) + "..."
}
