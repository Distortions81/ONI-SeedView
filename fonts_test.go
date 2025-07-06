//go:build test

package main

import "golang.org/x/image/font"

var (
	notoFont font.Face
	fontSize = 13.0
)

func setFontSize(size float64) { fontSize = size }
func increaseFontSize()        { setFontSize(fontSize + 2) }
func decreaseFontSize() {
	if fontSize > 6 {
		setFontSize(fontSize - 2)
	}
}
