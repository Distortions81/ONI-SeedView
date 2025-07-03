package main_test

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"testing"
)

type Point struct {
	X int
	Y int
}

type BiomePath struct {
	Name     string
	Polygons [][]Point
}

func colorFromARGB(hex uint32) (r uint8, g uint8, b uint8, a uint8) {
	return uint8(hex >> 16), uint8(hex >> 8), uint8(hex), uint8(hex >> 24)
}

func parseBiomeLine(line string) BiomePath {
	parts := strings.SplitN(strings.TrimSpace(line), ":", 2)
	if len(parts) != 2 {
		return BiomePath{}
	}
	bp := BiomePath{Name: parts[0]}
	segs := strings.Split(parts[1], ";")
	for _, seg := range segs {
		seg = strings.TrimSpace(seg)
		if seg == "" {
			continue
		}
		var poly []Point
		for _, pair := range strings.Fields(seg) {
			xy := strings.Split(pair, ",")
			if len(xy) != 2 {
				continue
			}
			x, err1 := strconv.Atoi(xy[0])
			y, err2 := strconv.Atoi(xy[1])
			if err1 != nil || err2 != nil {
				continue
			}
			poly = append(poly, Point{X: x, Y: y})
		}
		if len(poly) > 0 {
			bp.Polygons = append(bp.Polygons, poly)
		}
	}
	return bp
}

func parseBiomePaths(data string) []BiomePath {
	var paths []BiomePath
	lines := strings.Split(strings.ReplaceAll(strings.TrimSpace(data), "\r", ""), "\n")
	var current string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if strings.Contains(line, ":") {
			if current != "" {
				bp := parseBiomeLine(current)
				if bp.Name != "" {
					paths = append(paths, bp)
				}
			}
			current = line
		} else if current != "" {
			current += " " + line
		}
	}
	if current != "" {
		bp := parseBiomeLine(current)
		if bp.Name != "" {
			paths = append(paths, bp)
		}
	}
	return paths
}

func loadLocalSample(t *testing.T) struct {
	Asteroids []struct {
		BiomePaths string `json:"biomePaths"`
	}
} {
	b, err := os.ReadFile("render-issue/SNDST-A-7-0-0-0.json")
	if err != nil {
		t.Fatalf("read sample: %v", err)
	}
	var seed struct {
		Asteroids []struct {
			BiomePaths string `json:"biomePaths"`
		}
	}
	if err := json.Unmarshal(b, &seed); err != nil {
		t.Fatalf("decode sample: %v", err)
	}
	return seed
}

func TestColorFromARGB(t *testing.T) {
	r, g, b, a := colorFromARGB(0x11223344)
	if r != 0x22 || g != 0x33 || b != 0x44 || a != 0x11 {
		t.Errorf("got %02x %02x %02x %02x", r, g, b, a)
	}
}

func TestParseBiomePathsLocal(t *testing.T) {
	seed := loadLocalSample(t)
	if len(seed.Asteroids) == 0 {
		t.Fatal("no asteroids in sample")
	}
	paths := parseBiomePaths(seed.Asteroids[0].BiomePaths)
	if got := len(paths); got != 7 {
		t.Fatalf("expected 7 biomes, got %d", got)
	}
	counts := map[string]int{}
	for _, bp := range paths {
		counts[bp.Name] = len(bp.Polygons)
	}
	expected := map[string]int{
		"Sandstone":    1,
		"MagmaCore":    1,
		"OilField":     1,
		"ToxicJungle":  10,
		"Ocean":        6,
		"BoggyMarsh":   9,
		"FrozenWastes": 9,
	}
	for k, v := range expected {
		if counts[k] != v {
			t.Errorf("%s polygons: want %d got %d", k, v, counts[k])
		}
	}
}

func TestFetchBiomePaths(t *testing.T) {
	req, err := http.NewRequest("GET", "https://ingest.mapsnotincluded.org/coordinate/SNDST-A-7-0-0-0", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Accept", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Skipf("network unavailable: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Skipf("unexpected status: %d", resp.StatusCode)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("read: %v", err)
	}
	var seed struct {
		Asteroids []struct {
			BiomePaths string `json:"biomePaths"`
		}
	}
	if err := json.Unmarshal(body, &seed); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(seed.Asteroids) == 0 {
		t.Fatal("no asteroids")
	}
	paths := parseBiomePaths(seed.Asteroids[0].BiomePaths)
	if len(paths) != 7 {
		t.Fatalf("expected 7 biomes, got %d", len(paths))
	}
}
