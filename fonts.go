package main

import (
	_ "embed"

	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

//go:embed data/NotoSansMono.ttf
var notoTTF []byte

const (
	baseFontSize = 11.0
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
			panic("failed to parse font: " + err.Error())
		}
	}
	dpi := 72.0 * getHiDPIScale()
	face, err := opentype.NewFace(fontParsed, &opentype.FaceOptions{Size: size, DPI: dpi, Hinting: font.HintingFull})
	if err != nil {
		panic("failed to create font face: " + err.Error())
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

func registerFontChange(fn func()) {
	fontChange = append(fontChange, fn)
}

func fontScale() float64 { return fontSize / baseFontSize }

func rowSpacing() int {
	if notoFont != nil {
		return notoFont.Metrics().Height.Ceil() + uiScaled(8)
	}
	return int(float64(LegendRowSpacing) * fontScale())
}

func menuButtonHeight() int {
	if notoFont != nil {
		return notoFont.Metrics().Height.Ceil() + uiScaled(5)
	}
	return int(float64(22) * fontScale())
}

func menuSpacing() int {
	if notoFont != nil {
		return menuButtonHeight() + uiScaled(4)
	}
	return int(float64(26) * fontScale())
}

func seedBaseline() int {
	return uiScaled(10)
}

func init() {
	setUIScale(1.0)
}
