//go:build !test

package main

import (
	"embed"
	"encoding/base64"
	"log"

	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

//go:embed NotoSansMono.ttf.base64
var notoSansMonoBase64 string

var notoFont font.Face

func init() {
	data, err := base64.StdEncoding.DecodeString(notoSansMonoBase64)
	if err != nil {
		log.Fatalf("failed to decode font: %v", err)
	}
	tt, err := opentype.Parse(data)
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
