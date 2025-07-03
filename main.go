//go:build !test

package main

import (
	"embed"
	"encoding/json"
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"io"
	"net/http"
	"os"
	"path/filepath"
       "regexp"
       "sort"
       "strconv"
       "strings"
       "math"

	"github.com/chai2010/webp"
	"github.com/fxamacker/cbor/v2"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
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
	SizeX      int               `json:"sizeX"`
	SizeY      int               `json:"sizeY"`
	Geysers    []Geyser          `json:"geysers"`
	POIs       []PointOfInterest `json:"pointsOfInterest"`
	BiomePaths string            `json:"biomePaths"`
}

type SeedData struct {
	Asteroids []Asteroid `json:"asteroids"`
}

//go:embed assets/* names.json
var embeddedFiles embed.FS

type nameTables struct {
	Biomes  map[string]string `json:"biomes"`
	Geysers map[string]string `json:"geysers"`
	POIs    map[string]string `json:"pois"`
}

var names nameTables

func init() {
	if b, err := embeddedFiles.ReadFile("names.json"); err == nil {
		_ = json.Unmarshal(b, &names)
	}
}

func colorFromARGB(hex uint32) color.RGBA {
	return color.RGBA{
		R: uint8(hex >> 16),
		G: uint8(hex >> 8),
		B: uint8(hex),
		A: uint8(hex >> 24),
	}
}

var biomeColors = map[string]color.RGBA{
	"Sandstone":           colorFromARGB(0xFFF2BB47),
	"Barren":              colorFromARGB(0xFF97752C),
	"Space":               colorFromARGB(0xFF242424),
	"FrozenWastes":        colorFromARGB(0xFF9DC9D6),
	"BoggyMarsh":          colorFromARGB(0xFF7B974B),
	"ToxicJungle":         colorFromARGB(0xFFCB95A3),
	"Ocean":               colorFromARGB(0xFF4C4CFF),
	"Rust":                colorFromARGB(0xFFFFA007),
	"Forest":              colorFromARGB(0xFF8EC039),
	"Radioactive":         colorFromARGB(0xFF4AE458),
	"Swamp":               colorFromARGB(0xFFEB9B3F),
	"Wasteland":           colorFromARGB(0xFFCC3636),
	"Metallic":            colorFromARGB(0xFFFFA007),
	"Moo":                 colorFromARGB(0xFF8EC039),
	"IceCaves":            colorFromARGB(0xFFABCFEA),
	"CarrotQuarry":        colorFromARGB(0xFFCDA2C7),
	"SugarWoods":          colorFromARGB(0xFFA2CDA4),
	"PrehistoricGarden":   colorFromARGB(0xFF006127),
	"PrehistoricRaptor":   colorFromARGB(0xFF352F8C),
	"PrehistoricWetlands": colorFromARGB(0xFF645906),
	"OilField":            colorFromARGB(0xFF52321D),
	"MagmaCore":           colorFromARGB(0xFFDE5A3B),
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

// parseBiomePaths converts the raw biome path string into structured paths.
func parseBiomePaths(data string) []BiomePath {
	data = strings.ReplaceAll(strings.ReplaceAll(strings.TrimSpace(data), "\r", ""), "\n", " ")
	re := regexp.MustCompile(`([A-Za-z0-9_]+:)`)
	locs := re.FindAllStringIndex(data, -1)
	var paths []BiomePath
	for i, loc := range locs {
		name := strings.TrimSuffix(data[loc[0]:loc[1]], ":")
		start := loc[1]
		end := len(data)
		if i+1 < len(locs) {
			end = locs[i+1][0]
		}
		segment := strings.TrimSpace(data[start:end])
		if segment == "" {
			continue
		}
		bp := parseBiomeLine(name + ":" + segment)
		if bp.Name != "" {
			paths = append(paths, bp)
		}
	}
	return paths
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

// loadImage loads an image from the assets directory and caches it.
func loadImage(cache map[string]*ebiten.Image, name string) (*ebiten.Image, error) {
	if img, ok := cache[name]; ok {
		return img, nil
	}
	f, err := embeddedFiles.Open(filepath.Join("assets", name))
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var src image.Image
	ext := strings.ToLower(filepath.Ext(name))
	switch ext {
	case ".png":
		src, err = png.Decode(f)
	case ".jpg", ".jpeg":
		src, err = jpeg.Decode(f)
	case ".webp":
		src, err = webp.Decode(f)
	default:
		return nil, fmt.Errorf("unsupported image format: %s", ext)
	}
	if err != nil {
		return nil, err
	}
	img := ebiten.NewImageFromImage(src)
	cache[name] = img
	return img, nil
}

func iconForGeyser(id string) string {
	switch id {
	case "steam":
		return "geyser_cool_steam_vent.webp"
	case "hot_steam":
		return "geyser_steam_vent.webp"
	case "hot_water":
		return "geyser_water.webp"
	case "slush_water":
		return "geyser_cool_slush_geyser.webp"
	case "filthy_water":
		return "geyser_polluted_water_vent.webp"
	case "slush_salt_water":
		return "geyser_cool_salt_slush_geyser.webp"
	case "salt_water":
		return "geyser_salt_water.webp"
	case "small_volcano":
		return "geyser_minor_volcano.webp"
	case "big_volcano":
		return "geyser_volcano.webp"
	case "liquid_co2":
		return "geyser_carbon_dioxide.webp"
	case "hot_co2":
		return "geyser_carbon_dioxide_vent.webp"
	case "hot_hydrogen":
		return "geyser_hydrogen_vent.webp"
	case "hot_po2":
		return "geyser_hot_polluted_oxygen_vent.webp"
	case "slimy_po2":
		return "geyser_infectious_polluted_oxygen_vent.webp"
	case "chlorine_gas", "chlorine_gas_cool":
		return "geyser_chlorine_gas_vent.webp"
	case "methane":
		return "geyser_natural_gas_geyser.webp"
	case "molten_copper":
		return "geyser_copper_volcano.webp"
	case "molten_iron":
		return "geyser_iron_volcano.webp"
	case "molten_gold":
		return "geyser_gold_volcano.webp"
	case "oil_drip":
		return "geyser_leaky_oil_fissure.webp"
	case "molten_aluminum":
		return "geyser_aluminum_volcano.webp"
	case "molten_cobalt":
		return "geyser_cobalt_volcano.webp"
	case "liquid_sulfur":
		return "geyser_liquid_sulfur_geyser.webp"
	case "molten_tungsten":
		return "geyser_tungsten_volcano.webp"
	case "molten_niobium":
		return "geyser_niobium_volcano.webp"
	case "OilWell":
		return "geyser_oil_reservoir.webp"
	default:
		return ""
	}
}

func iconForPOI(id string) string {
	switch id {
	case "Headquarters":
		return "building_printing_pod.webp"
	case "WarpConduitSender":
		return "building_supply_teleporter_input.webp"
	case "WarpConduitReceiver":
		return "building_supply_teleporter_output.webp"
	case "WarpPortal":
		return "building_teleporter_transmitter.webp"
	case "WarpReceiver":
		return "building_teleporter_receiver.webp"
	case "GeneShuffler":
		return "building_neural_vacillator.webp"
	case "MassiveHeatSink":
		return "building_anti_entropy_thermo_nullifier.webp"
	case "SapTree":
		return "building_sap_tree.webp"
	case "GravitasPedestal":
		return "poi_artifact_outline.webp"
	case "PropSurfaceSatellite1":
		return "poi_crashed_satellite.webp"
	case "PropSurfaceSatellite2":
		return "poi_wrecked_satellite.webp"
	case "PropSurfaceSatellite3":
		return "poi_crushed_satellite.webp"
	case "TemporalTearOpener":
		return "building_temporal_tear_opener.webp"
	case "CryoTank":
		return "building_cryotank.webp"
	case "PropFacilityStatue":
		return "poi_prop_facility_statue.webp"
	case "GeothermalVentEntity":
		return "poi_geothermal_vent_entity.webp"
	case "GeothermalControllerEntity":
		return "poi_geothermal_controller_entity.webp"
	case "POICeresTechUnlock":
		return "poi_ceres_tech_unlock.webp"
	default:
		return ""
	}
}

func displayBiome(id string) string {
	if v, ok := names.Biomes[id]; ok {
		return v
	}
	return id
}

func displayGeyser(id string) string {
	if v, ok := names.Geysers[id]; ok {
		return v
	}
	return id
}

func displayPOI(id string) string {
	if v, ok := names.POIs[id]; ok {
		return v
	}
	return id
}

var whitePixel = func() *ebiten.Image {
	img := ebiten.NewImage(1, 1)
	img.Fill(color.White)
	return img
}()

func drawPolygon(dst *ebiten.Image, pts []Point, clr color.Color, camX, camY, zoom float64) {
	if len(pts) == 0 {
		return
	}
	var p vector.Path
	p.MoveTo(float32(pts[0].X*2), float32(pts[0].Y*2))
	for _, pt := range pts[1:] {
		p.LineTo(float32(pt.X*2), float32(pt.Y*2))
	}
	p.Close()
	vs, is := p.AppendVerticesAndIndicesForFilling(nil, nil)
	r, g, b, a := clr.RGBA()
       for i := range vs {
               x := float64(vs[i].DstX)*zoom + camX
               y := float64(vs[i].DstY)*zoom + camY
               vs[i].DstX = float32(math.Round(x))
               vs[i].DstY = float32(math.Round(y))
               vs[i].SrcX = 0
               vs[i].SrcY = 0
               vs[i].ColorR = float32(r) / 0xffff
               vs[i].ColorG = float32(g) / 0xffff
               vs[i].ColorB = float32(b) / 0xffff
		vs[i].ColorA = float32(a) / 0xffff
	}
	op := &ebiten.DrawTrianglesOptions{
		AntiAlias:      true,
		ColorScaleMode: ebiten.ColorScaleModeStraightAlpha,
		FillRule:       ebiten.FillRuleEvenOdd,
	}
	dst.DrawTriangles(vs, is, whitePixel, op)
}

// drawBiome draws all polygons for a biome using a single path so holes render correctly.
func drawBiome(dst *ebiten.Image, polys [][]Point, clr color.Color, camX, camY, zoom float64) {
	if len(polys) == 0 {
		return
	}
	var p vector.Path
	for _, pts := range polys {
		if len(pts) == 0 {
			continue
		}
		p.MoveTo(float32(pts[0].X*2), float32(pts[0].Y*2))
		for _, pt := range pts[1:] {
			p.LineTo(float32(pt.X*2), float32(pt.Y*2))
		}
		p.Close()
	}
	vs, is := p.AppendVerticesAndIndicesForFilling(nil, nil)
	r, g, b, a := clr.RGBA()
       for i := range vs {
               x := float64(vs[i].DstX)*zoom + camX
               y := float64(vs[i].DstY)*zoom + camY
               vs[i].DstX = float32(math.Round(x))
               vs[i].DstY = float32(math.Round(y))
               vs[i].SrcX = 0
               vs[i].SrcY = 0
               vs[i].ColorR = float32(r) / 0xffff
               vs[i].ColorG = float32(g) / 0xffff
               vs[i].ColorB = float32(b) / 0xffff
		vs[i].ColorA = float32(a) / 0xffff
	}
	op := &ebiten.DrawTrianglesOptions{
		AntiAlias:      true,
		ColorScaleMode: ebiten.ColorScaleModeStraightAlpha,
		FillRule:       ebiten.FillRuleEvenOdd,
	}
	dst.DrawTriangles(vs, is, whitePixel, op)
}

// drawLegend renders a list of biome colors along the left side of the screen.
func drawLegend(dst *ebiten.Image, biomes []BiomePath) {
	set := make(map[string]struct{})
	for _, b := range biomes {
		set[b.Name] = struct{}{}
	}
	names := make([]string, 0, len(set))
	for n := range set {
		names = append(names, n)
	}
	sort.Strings(names)
	y := 10
	for _, name := range names {
		clr, ok := biomeColors[name]
		if !ok {
			clr = color.RGBA{60, 60, 60, 255}
		}
		vector.DrawFilledRect(dst, 5, float32(y), 20, 10, clr, false)
		ebitenutil.DebugPrintAt(dst, displayBiome(name), 30, y)
		y += 15
	}
}

// Game implements ebiten.Game and displays geysers with their names.
type Game struct {
	geysers     []Geyser
	pois        []PointOfInterest
	biomes      []BiomePath
	icons       map[string]*ebiten.Image
	width       int
	height      int
	astWidth    int
	astHeight   int
	camX        float64
	camY        float64
	zoom        float64
	dragging    bool
	lastX       int
	lastY       int
	needsRedraw bool
}

type label struct {
	text string
	x    int
	y    int
}

func (g *Game) Update() error {
	const panSpeed = 5

	oldX, oldY, oldZoom := g.camX, g.camY, g.zoom

	// Keyboard panning
	if ebiten.IsKeyPressed(ebiten.KeyLeft) || ebiten.IsKeyPressed(ebiten.KeyA) {
		g.camX += panSpeed
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) || ebiten.IsKeyPressed(ebiten.KeyD) {
		g.camX -= panSpeed
	}
	if ebiten.IsKeyPressed(ebiten.KeyUp) || ebiten.IsKeyPressed(ebiten.KeyW) {
		g.camY += panSpeed
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) || ebiten.IsKeyPressed(ebiten.KeyS) {
		g.camY -= panSpeed
	}

	// Mouse dragging
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		if g.dragging {
			g.camX += float64(x - g.lastX)
			g.camY += float64(y - g.lastY)
		}
		g.lastX = x
		g.lastY = y
		g.dragging = true
	} else {
		g.dragging = false
	}

	// Zoom with keyboard
	zoomFactor := 1.0
	if ebiten.IsKeyPressed(ebiten.KeyEqual) || ebiten.IsKeyPressed(ebiten.KeyKPAdd) {
		zoomFactor *= 1.05
	}
	if ebiten.IsKeyPressed(ebiten.KeyMinus) || ebiten.IsKeyPressed(ebiten.KeyKPSubtract) {
		zoomFactor /= 1.05
	}

	// Zoom with mouse wheel
	_, wheelY := ebiten.Wheel()
	if wheelY != 0 {
		if wheelY > 0 {
			zoomFactor *= 1.1
		} else {
			zoomFactor /= 1.1
		}
	}

	if zoomFactor != 1.0 {
		oldZoom := g.zoom
		g.zoom *= zoomFactor
		if g.zoom < 0.1 {
			g.zoom = 0.1
		}
		cx, cy := float64(g.width)/2, float64(g.height)/2
		worldX := (cx - g.camX) / oldZoom
		worldY := (cy - g.camY) / oldZoom
		g.camX = cx - worldX*g.zoom
		g.camY = cy - worldY*g.zoom
	}

	if g.camX != oldX || g.camY != oldY || g.zoom != oldZoom {
		g.needsRedraw = true
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	if g.needsRedraw {
		screen.Fill(color.RGBA{30, 30, 30, 255})
		labels := []label{}

		for _, bp := range g.biomes {
			clr, ok := biomeColors[bp.Name]
			if !ok {
				clr = color.RGBA{60, 60, 60, 255}
			}
			drawBiome(screen, bp.Polygons, clr, g.camX, g.camY, g.zoom)
		}

		for _, gy := range g.geysers {
			x := (float64(gy.X) * 2 * g.zoom) + g.camX
			y := (float64(gy.Y) * 2 * g.zoom) + g.camY
			if iconName := iconForGeyser(gy.ID); iconName != "" {
				if img, err := loadImage(g.icons, iconName); err == nil {
					op := &ebiten.DrawImageOptions{Filter: ebiten.FilterLinear}
					op.GeoM.Scale(g.zoom*0.25, g.zoom*0.25)
					w := float64(img.Bounds().Dx()) * g.zoom * 0.25
					h := float64(img.Bounds().Dy()) * g.zoom * 0.25
					op.GeoM.Translate(x-w/2, y-h/2)
					screen.DrawImage(img, op)
					d := displayGeyser(gy.ID)
					labels = append(labels, label{d, int(x) - (len(d)*6)/2, int(y+h/2) + 2})
				}
			} else {
				vector.DrawFilledRect(screen, float32(x-2), float32(y-2), 4, 4, color.RGBA{255, 0, 0, 255}, false)
				d := displayGeyser(gy.ID)
				labels = append(labels, label{d, int(x) - (len(d)*6)/2, int(y) + 4})
			}
		}

		for _, poi := range g.pois {
			x := (float64(poi.X) * 2 * g.zoom) + g.camX
			y := (float64(poi.Y) * 2 * g.zoom) + g.camY
			if iconName := iconForPOI(poi.ID); iconName != "" {
				if img, err := loadImage(g.icons, iconName); err == nil {
					op := &ebiten.DrawImageOptions{Filter: ebiten.FilterLinear}
					op.GeoM.Scale(g.zoom*0.25, g.zoom*0.25)
					w := float64(img.Bounds().Dx()) * g.zoom * 0.25
					h := float64(img.Bounds().Dy()) * g.zoom * 0.25
					op.GeoM.Translate(x-w/2, y-h/2)
					screen.DrawImage(img, op)
					d := displayPOI(poi.ID)
					labels = append(labels, label{d, int(x) - (len(d)*6)/2, int(y+h/2) + 2})
				}
			} else {
				d := displayPOI(poi.ID)
				labels = append(labels, label{d, int(x) - (len(d)*6)/2, int(y) + 4})
			}
		}

		drawLegend(screen, g.biomes)
		for _, l := range labels {
			ebitenutil.DebugPrintAt(screen, l.text, l.x, l.y)
		}

		g.needsRedraw = false
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	if g.width != outsideWidth || g.height != outsideHeight {
		// Keep the world position at the center of the screen fixed so
		// resizing doesn't shift the view.
		cxOld, cyOld := float64(g.width)/2, float64(g.height)/2
		worldX := (cxOld - g.camX) / g.zoom
		worldY := (cyOld - g.camY) / g.zoom

		g.width = outsideWidth
		g.height = outsideHeight

		cxNew, cyNew := float64(g.width)/2, float64(g.height)/2
		g.camX = cxNew - worldX*g.zoom
		g.camY = cyNew - worldY*g.zoom

		g.needsRedraw = true
	}
	return outsideWidth, outsideHeight
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
	game := &Game{
		geysers:   ast.Geysers,
		pois:      ast.POIs,
		biomes:    parseBiomePaths(ast.BiomePaths),
		icons:     make(map[string]*ebiten.Image),
		width:     600,
		height:    800,
		astWidth:  ast.SizeX,
		astHeight: ast.SizeY,
		zoom:      1.0,
	}
	game.camX = (float64(game.width) - float64(game.astWidth)*2*game.zoom) / 2
	game.camY = (float64(game.height) - float64(game.astHeight)*2*game.zoom) / 2
	game.needsRedraw = true
	ebiten.SetWindowSize(game.width, game.height)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetWindowTitle("Geysers - " + *coord)
	ebiten.SetScreenClearedEveryFrame(false)
	ebiten.SetVsyncEnabled(true)
	if err := ebiten.RunGame(game); err != nil {
		fmt.Println("Error running game:", err)
	}
}
