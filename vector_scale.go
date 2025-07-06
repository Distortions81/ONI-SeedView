package main

import "github.com/hajimehoshi/ebiten/v2"

func vectorScale() float32 {
	return float32(ebiten.Monitor().DeviceScaleFactor())
}
