//go:build !test

package main

import (
	_ "embed"
	"log"

	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

//go:embed NotoSansMono.ttf
var notoTTF []byte

var (
	notoFont   font.Face
	fontSize   = 13.0
	fontParsed *opentype.Font
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
}

func increaseFontSize() { setFontSize(fontSize + 2) }

func decreaseFontSize() {
	if fontSize > 6 {
		setFontSize(fontSize - 2)
	}
}

func init() {
	setFontSize(fontSize)
}
