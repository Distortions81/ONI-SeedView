package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"image/color"
	"io"
	"net/http"
	"os"

	"github.com/fxamacker/cbor/v2"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// fetchSeedCBOR retrieves the seed data in CBOR format for a given coordinate.
func fetchSeedCBOR(coordinate string) ([]byte, error) {
	url := "https://ingest.mapsnotincluded.org/coordinate/" + coordinate
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Accept", "application/cbor")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected status %d: %s", resp.StatusCode, body)
	}
	return io.ReadAll(resp.Body)
}

// normalize converts maps with interface{} keys so JSON encoding works.
func normalize(value interface{}) interface{} {
	switch v := value.(type) {
	case map[interface{}]interface{}:
		m := make(map[string]interface{})
		for key, val := range v {
			m[fmt.Sprintf("%v", key)] = normalize(val)
		}
		return m
	case map[string]interface{}:
		for key, val := range v {
			v[key] = normalize(val)
		}
		return v
	case []interface{}:
		for i, val := range v {
			v[i] = normalize(val)
		}
		return v
	default:
		return v
	}
}

// decodeCBORToJSON converts CBOR bytes into pretty JSON.
func decodeCBORToJSON(cborData []byte) ([]byte, error) {
	var decoded interface{}
	if err := cbor.Unmarshal(cborData, &decoded); err != nil {
		return nil, fmt.Errorf("CBOR decode failed: %v", err)
	}

	jsonData, err := json.MarshalIndent(normalize(decoded), "", "  ")
	if err != nil {
		return nil, fmt.Errorf("JSON encode failed: %v", err)
	}
	return jsonData, nil
}

func saveToFile(filename string, data []byte) error {
	return os.WriteFile(filename, data, 0644)
}

// Data structures used to decode the geyser information.
type Geyser struct {
	ID string `json:"id"`
	X  int    `json:"x"`
	Y  int    `json:"y"`
}

type Asteroid struct {
	SizeX   int      `json:"sizeX"`
	SizeY   int      `json:"sizeY"`
	Geysers []Geyser `json:"geysers"`
}

type SeedData struct {
	Asteroids []Asteroid `json:"asteroids"`
}

// decodeSeed parses the CBOR seed data into SeedData.
func decodeSeed(cborData []byte) (*SeedData, error) {
	jsonData, err := decodeCBORToJSON(cborData)
	if err != nil {
		return nil, err
	}
	var seed SeedData
	if err := json.Unmarshal(jsonData, &seed); err != nil {
		return nil, fmt.Errorf("JSON decode failed: %v", err)
	}
	return &seed, nil
}

// Game implements ebiten.Game and displays geysers with their names.
type Game struct {
	geysers []Geyser
	width   int
	height  int
}

func (g *Game) Update() error { return nil }

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{30, 30, 30, 255})
	for _, gy := range g.geysers {
		// Invert Y to place origin at bottom-left like the game world.
		y := g.height - gy.Y
		ebitenutil.DrawRect(screen, float64(gy.X-2), float64(y-2), 4, 4, color.RGBA{255, 0, 0, 255})
		ebitenutil.DebugPrintAt(screen, gy.ID, gy.X+4, y-4)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return g.width, g.height
}

func main() {
	coord := flag.String("coord", "SNDST-A-7-0-0-0", "seed coordinate")
	out := flag.String("out", "", "optional path to save JSON")
	flag.Parse()

	fmt.Println("Fetching:", *coord)
	cborData, err := fetchSeedCBOR(*coord)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	seed, err := decodeSeed(cborData)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	if *out != "" {
		jsonData, _ := json.MarshalIndent(seed, "", "  ")
		if err := saveToFile(*out, jsonData); err != nil {
			fmt.Println("Error writing file:", err)
		}
	}

	ast := seed.Asteroids[0]
	game := &Game{geysers: ast.Geysers, width: ast.SizeX, height: ast.SizeY}
	ebiten.SetWindowSize(game.width*2, game.height*2)
	ebiten.SetWindowTitle("Geysers - " + *coord)
	if err := ebiten.RunGame(game); err != nil {
		fmt.Println("Error running game:", err)
	}
}
