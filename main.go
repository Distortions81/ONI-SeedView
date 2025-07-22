package main

import (
	"flag"
	"fmt"
	"runtime"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	coord := flag.String("coord", "V-FRST-C-1331877-0-0-0", "seed coordinate")
	screenshot := flag.String("screenshot", "", "path to save a PNG screenshot and exit")
	flag.Parse()
	asteroidIDVal := ""
	asteroidSpecified := false
	if runtime.GOARCH == "wasm" {
		if c := coordFromURL(); c != "" {
			c = strings.TrimSpace(c)
			*coord = c
		}
		if a, ok := asteroidFromURL(); ok {
			asteroidIDVal = a
			asteroidSpecified = true
		}
	}

	game := &Game{
		icons:             make(map[string]*ebiten.Image),
		width:             DefaultWidth,
		height:            DefaultHeight,
		zoom:              1.0,
		minZoom:           MinZoom,
		loading:           true,
		status:            "Fetching...",
		statusError:       false,
		coord:             *coord,
		asteroidID:        asteroidIDVal,
		asteroidSpecified: asteroidSpecified,
		textures:          true,
		vsync:             true,
		showItemNames:     true,
		showLegend:        true,
		mobile:            isMobile(),
		useNumbers:        !isMobile(),
		iconScale:         1.0,
		smartRender:       true,
		linearFilter:      true,
		hidpi:             true,
		ssQuality:         1,
		hoverBiome:        -1,
		hoverItem:         -1,
		selectedBiome:     -1,
		selectedItem:      -1,
	}
	game.initEUI()
	setHiDPI(game.hidpi)
	registerFontChange(game.invalidateLegends)
	loadGameData(game, *coord, asteroidIDVal)
	if *screenshot != "" {

		game.screenshotPath = *screenshot
		//game.screenshotMode = true
		//game.showHelp = false
		//game.showInfo = false
		game.showShotMenu = false

	}
	ebiten.SetWindowSize(game.width, game.height)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetWindowTitle("Geysers - " + *coord)
	ebiten.SetScreenClearedEveryFrame(false)
	ebiten.SetVsyncEnabled(game.vsync)
	if err := ebiten.RunGame(game); err != nil {
		fmt.Println("Error running game:", err)
	}
}

func loadGameData(game *Game, coord, asteroidID string) {
	cborData, err := fetchSeedCBOR(coord)
	if err != nil {
		game.status = "Error: " + err.Error()
		game.statusError = false
		game.needsRedraw = true
		game.loading = false
		return
	}
	seed, err := decodeSeed(cborData)
	game.asteroids = seed.Asteroids
	if err != nil {
		game.status = "Error: " + err.Error()
		game.statusError = false
		game.needsRedraw = true
		game.loading = false
		return
	}
	astIdxSel := 0
	if game.asteroidSpecified {
		astIdxSel = asteroidIndexByID(seed.Asteroids, asteroidID)
		if astIdxSel < 0 {
			game.status = fmt.Sprintf("%s\nAsteroid ID: %s\nThis location does not contain Asteroid ID: %s", game.coord, asteroidID, asteroidID)
			if len(seed.Asteroids) > 0 {
				valid := make([]string, 0, len(seed.Asteroids))
				for _, a := range seed.Asteroids {
					valid = append(valid, a.ID)
				}
				lines := make([]string, 0, (len(valid)+2)/3)
				for i, v := range valid {
					if i%3 == 0 {
						lines = append(lines, v)
					} else {
						lines[len(lines)-1] = lines[len(lines)-1] + ", " + v
					}
				}
				game.status += fmt.Sprintf("\nValid IDs: %s", strings.Join(lines, "\n"))
			}
			game.statusError = true
			game.needsRedraw = true
			game.loading = false
			return
		}
	}
	ast := seed.Asteroids[astIdxSel]
	game.invalidateLegends()
	game.legendMap = nil
	game.legendEntries = nil
	game.legendColors = nil
	game.selectedItem = -1
	game.itemScroll = 0
	game.asteroidID = ast.ID
	bps := parseBiomePaths(ast.BiomePaths)
	game.geysers = ast.Geysers
	game.pois = ast.POIs
	game.biomes = bps
	game.astWidth = ast.SizeX
	game.astHeight = ast.SizeY
	game.legend, game.legendBiomes = buildLegendImage(bps)
	game.fitOnLoad = true
	game.biomeTextures = loadBiomeTextures()
	names := []string{"../icons/camera.png", "../icons/help.png", "../icons/gear.png", "geyser_water.png"}
	set := make(map[string]struct{})
	for _, gy := range ast.Geysers {
		if n := iconForGeyser(gy.ID); n != "" {
			if _, ok := set[n]; !ok {
				set[n] = struct{}{}
				names = append(names, n)
			}
		}
	}
	for _, poi := range ast.POIs {
		if n := iconForPOI(poi.ID); n != "" {
			if _, ok := set[n]; !ok {
				set[n] = struct{}{}
				names = append(names, n)
			}
		}
	}
	game.startIconLoader(names)
	game.loading = false
	game.needsRedraw = true
}
