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

var notoFont font.Face

func init() {
	tt, err := opentype.Parse(notoTTF)
	if err != nil {
		log.Fatalf("failed to parse font: %v", err)
	}
	const dpi = 72
	face, err := opentype.NewFace(tt, &opentype.FaceOptions{Size: 14, DPI: dpi, Hinting: font.HintingFull})
	if err != nil {
		log.Fatalf("failed to create font face: %v", err)
	}
	notoFont = face
}
