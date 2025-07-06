//go:build test

package main

import "golang.org/x/image/font"

const baseFontSize = 12.0

var (
	notoFont   font.Face
	fontSize   = baseFontSize
	fontChange []func()
)

func setFontSize(size float64) {
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

func fontScale() float64 { return fontSize / baseFontSize }

func rowSpacing() int { return int(float64(LegendRowSpacing) * fontScale()) }

func registerFontChange(fn func()) { fontChange = append(fontChange, fn) }

func menuButtonHeight() int { return int(float64(22) * fontScale()) }

func menuSpacing() int { return int(float64(26) * fontScale()) }
