//go:build !test

package main

import (
	_ "embed"
	"log"

	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

//go:embed data/NotoSansMono.ttf
var notoTTF []byte

const (
	baseFontSize       = 12.0
	screenshotFontSize = 14.0
)

var (
	notoFont   font.Face
	fontSize   = baseFontSize
	fontParsed *opentype.Font
	fontChange []func()
)

func loadFont(size float64) font.Face {
	if fontParsed == nil {
		var err error
		fontParsed, err = opentype.Parse(notoTTF)
		if err != nil {
			log.Fatalf("failed to parse font: %v", err)
		}
	}
	const dpi = 72
	face, err := opentype.NewFace(fontParsed, &opentype.FaceOptions{Size: size, DPI: dpi, Hinting: font.HintingFull})
	if err != nil {
		log.Fatalf("failed to create font face: %v", err)
	}
	return face
}

func setFontSize(size float64) {
	notoFont = loadFont(size)
	fontSize = size
	for _, cb := range fontChange {
		cb()
	}
}

func increaseFontSize() { setFontSize(fontSize + 2) }

func decreaseFontSize() {
	if fontSize > 6 {
		setFontSize(fontSize - 2)
	}
}

func registerFontChange(fn func()) {
	fontChange = append(fontChange, fn)
}

func fontScale() float64 { return fontSize / baseFontSize }

func rowSpacing() int {
	if notoFont != nil {
		return notoFont.Metrics().Height.Ceil() + 8
	}
	return int(float64(LegendRowSpacing) * fontScale())
}

func menuButtonHeight() int {
	if notoFont != nil {
		return notoFont.Metrics().Height.Ceil() + 5
	}
	return int(float64(22) * fontScale())
}

func menuSpacing() int {
	if notoFont != nil {
		return menuButtonHeight() + 4
	}
	return int(float64(26) * fontScale())
}

func init() {
	setFontSize(fontSize)
}
