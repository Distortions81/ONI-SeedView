//go:build test

package main

import "golang.org/x/image/font"

const baseFontSize = 13.0

var (
	notoFont font.Face
	fontSize = baseFontSize
)

func setFontSize(size float64) { fontSize = size }
func increaseFontSize()        { setFontSize(fontSize + 2) }
func decreaseFontSize() {
	if fontSize > 6 {
		setFontSize(fontSize - 2)
	}
}

func fontScale() float64 { return fontSize / baseFontSize }
