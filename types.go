package main

import (
	"encoding/json"

	_ "embed"
)

// Data structures used to decode the geyser information.
type Geyser struct {
	ID             string  `json:"id"`
	X              int     `json:"x"`
	Y              int     `json:"y"`
	ActiveCycles   float64 `json:"activeCycles"`
	AvgEmitRate    float64 `json:"avgEmitRate"`
	DormancyCycles float64 `json:"dormancyCycles"`
	EmitRate       float64 `json:"emitRate"`
	EruptionTime   float64 `json:"eruptionTime"`
	IdleTime       float64 `json:"idleTime"`
}

type PointOfInterest struct {
	ID string `json:"id"`
	X  int    `json:"x"`
	Y  int    `json:"y"`
}

type Point struct {
	X int
	Y int
}

type BiomePath struct {
	Name     string
	Polygons [][]Point
}

type Asteroid struct {
	ID         string            `json:"id"`
	SizeX      int               `json:"sizeX"`
	SizeY      int               `json:"sizeY"`
	Geysers    []Geyser          `json:"geysers"`
	POIs       []PointOfInterest `json:"pointsOfInterest"`
	BiomePaths string            `json:"biomePaths"`
}

type SeedData struct {
	Asteroids []Asteroid `json:"asteroids"`
}

type nameTables struct {
	Biomes  map[string]string `json:"biomes"`
	Geysers map[string]string `json:"geysers"`
	POIs    map[string]string `json:"pois"`
}

var names nameTables

//go:embed data/names.json
var namesData []byte

func init() {
	_ = json.Unmarshal(namesData, &names)
}
