//go:build !test

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"
	"io"
	"math"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"time"

	_ "embed"
	"github.com/fxamacker/cbor/v2"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const helpMessage = "Controls:\n" +
	"Arrow keys/WASD or drag to pan\n" +
	"Mouse wheel or +/- to zoom\n" +
	"Pinch to zoom on touch"

// fetchSeedCBOR retrieves the seed data in CBOR format for a given coordinate.
func fetchSeedCBOR(coordinate string) ([]byte, error) {
	url := BaseURL + coordinate
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Accept", AcceptCBORHeader)

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

type nameTables struct {
	Biomes  map[string]string `json:"biomes"`
	Geysers map[string]string `json:"geysers"`
	POIs    map[string]string `json:"pois"`
}

var names nameTables

//go:embed names.json
var namesData []byte

func openFile(name string) (io.ReadCloser, error) {
	if runtime.GOOS == "js" {
		resp, err := http.Get(name)
		if err != nil {
			return nil, err
		}
		if resp.StatusCode != http.StatusOK {
			b, _ := io.ReadAll(resp.Body)
			resp.Body.Close()
			return nil, fmt.Errorf("failed to fetch %s: %s", name, b)
		}
		return resp.Body, nil
	}
	return os.Open(name)
}

func openAsset(name string) (io.ReadCloser, error) {
	if runtime.GOOS == "js" {
		return openFile(path.Join(WebAssetBase, name))
	}
	return openFile(filepath.Join("assets", name))
}

func init() {
	_ = json.Unmarshal(namesData, &names)
}

func colorFromARGB(hex uint32) color.RGBA {
	return color.RGBA{
		R: uint8(hex >> 16),
		G: uint8(hex >> 8),
		B: uint8(hex),
		A: uint8(hex >> 24),
	}
}

func darkenColor(c color.RGBA, factor float64) color.RGBA {
	return color.RGBA{
		R: uint8(float64(c.R) * factor),
		G: uint8(float64(c.G) * factor),
		B: uint8(float64(c.B) * factor),
		A: c.A,
	}
}

var biomeColors = map[string]color.RGBA{
	"Sandstone": darkenColor(colorFromARGB(0xFFF2BB47), 0.6),
	"Barren":    darkenColor(colorFromARGB(0xFF97752C), 0.6),
	// Use a lighter gray so the space region stands out from the
	// background color of the window.
	"Space":               darkenColor(colorFromARGB(0xFFD0D0D0), 0.6),
	"FrozenWastes":        darkenColor(colorFromARGB(0xFF4F80B5), 0.6),
	"BoggyMarsh":          darkenColor(colorFromARGB(0xFF7B974B), 0.6),
	"ToxicJungle":         darkenColor(colorFromARGB(0xFFCB95A3), 0.6),
	"Ocean":               darkenColor(colorFromARGB(0xFF4C4CFF), 0.6),
	"Rust":                darkenColor(colorFromARGB(0xFFFFA007), 0.6),
	"Forest":              darkenColor(colorFromARGB(0xFF8EC039), 0.6),
	"Radioactive":         darkenColor(colorFromARGB(0xFF4AE458), 0.6),
	"Swamp":               darkenColor(colorFromARGB(0xFFEB9B3F), 0.6),
	"Wasteland":           darkenColor(colorFromARGB(0xFFCC3636), 0.6),
	"Metallic":            darkenColor(colorFromARGB(0xFFFFA007), 0.6),
	"Moo":                 darkenColor(colorFromARGB(0xFF8EC039), 0.6),
	"IceCaves":            darkenColor(colorFromARGB(0xFF6C9BD3), 0.6),
	"CarrotQuarry":        darkenColor(colorFromARGB(0xFFCDA2C7), 0.6),
	"SugarWoods":          darkenColor(colorFromARGB(0xFFA2CDA4), 0.6),
	"PrehistoricGarden":   darkenColor(colorFromARGB(0xFF006127), 0.6),
	"PrehistoricRaptor":   darkenColor(colorFromARGB(0xFF352F8C), 0.6),
	"PrehistoricWetlands": darkenColor(colorFromARGB(0xFF645906), 0.6),
	"OilField":            darkenColor(colorFromARGB(0xFF52321D), 0.6),
	"MagmaCore":           darkenColor(colorFromARGB(0xFFDE5A3B), 0.6),
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
func toCamel(s string) string {
	parts := strings.Split(s, "_")
	for i, p := range parts {
		if len(p) == 0 {
			continue
		}
		parts[i] = strings.ToUpper(p[:1]) + p[1:]
	}
	return strings.Join(parts, "")
}

func assetExists(name string) bool {
	f, err := openAsset(name)
	if err != nil {
		return false
	}
	f.Close()
	return true
}

func resolveAssetName(name string) string {
	if assetExists(name) {
		return name
	}
	ext := filepath.Ext(name)
	base := strings.TrimSuffix(name, ext)
	camel := toCamel(base) + ext
	if assetExists(camel) {
		return camel
	}
	for _, prefix := range []string{"geyser_", "building_", "poi_"} {
		if strings.HasPrefix(base, prefix) {
			trimmed := base[len(prefix):]
			camel = toCamel(trimmed) + ext
			if assetExists(camel) {
				return camel
			}
		}
	}
	return name
}

var missingImage = ebiten.NewImage(1, 1)

func loadImageFile(name string) (*ebiten.Image, error) {
	resolved := resolveAssetName(name)
	f, err := openAsset(resolved)
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
	default:
		return nil, fmt.Errorf("unsupported image format: %s", ext)
	}
	if err != nil {
		return nil, err
	}
	if nrgba, ok := src.(*image.NRGBA); ok {
		for i := 3; i < len(nrgba.Pix); i += 4 {
			if nrgba.Pix[i] < 64 {
				nrgba.Pix[i] = 0
			}
		}
		return ebiten.NewImageFromImage(nrgba), nil
	}
	bounds := src.Bounds()
	dst := image.NewNRGBA(bounds)
	draw.Draw(dst, bounds, src, bounds.Min, draw.Src)
	for i := 3; i < len(dst.Pix); i += 4 {
		if dst.Pix[i] < 64 {
			dst.Pix[i] = 0
		}
	}
	return ebiten.NewImageFromImage(dst), nil
}

func loadImage(cache map[string]*ebiten.Image, name string) (*ebiten.Image, error) {
	if img, ok := cache[name]; ok {
		if img == missingImage {
			return nil, fmt.Errorf("missing")
		}
		return img, nil
	}
	img, err := loadImageFile(name)
	if err != nil {
		cache[name] = missingImage
		return nil, err
	}
	cache[name] = img
	return img, nil
}

// simplifyID strips namespace prefixes from an ID string. The seed data often
// includes paths like "expansion1::geysers/liquid_sulfur". This helper returns
// only the final identifier ("liquid_sulfur" in that example).
func simplifyID(id string) string {
	if idx := strings.LastIndex(id, "/"); idx >= 0 {
		id = id[idx+1:]
	}
	if idx := strings.LastIndex(id, "::"); idx >= 0 {
		id = id[idx+2:]
	}
	return id
}

func iconForGeyser(id string) string {
	id = simplifyID(id)
	switch id {
	case "steam":
		return "geyser_cool_steam_vent.png"
	case "hot_steam":
		return "geyser_steam_vent.png"
	case "hot_water":
		return "geyser_water.png"
	case "slush_water":
		return "geyser_cool_slush_geyser.png"
	case "filthy_water":
		return "geyser_polluted_water_vent.png"
	case "slush_salt_water":
		return "geyser_cool_salt_slush_geyser.png"
	case "salt_water":
		return "geyser_salt_water.png"
	case "small_volcano":
		return "geyser_minor_volcano.png"
	case "big_volcano":
		return "geyser_volcano.png"
	case "liquid_co2":
		return "geyser_carbon_dioxide.png"
	case "hot_co2":
		return "geyser_carbon_dioxide_vent.png"
	case "hot_hydrogen":
		return "geyser_hydrogen_vent.png"
	case "hot_po2":
		return "geyser_hot_polluted_oxygen_vent.png"
	case "slimy_po2":
		return "geyser_infectious_polluted_oxygen_vent.png"
	case "chlorine_gas", "chlorine_gas_cool":
		return "geyser_chlorine_gas_vent.png"
	case "methane":
		return "geyser_natural_gas_geyser.png"
	case "molten_copper":
		return "geyser_copper_volcano.png"
	case "molten_iron":
		return "geyser_iron_volcano.png"
	case "molten_gold":
		return "geyser_gold_volcano.png"
	case "oil_drip":
		return "geyser_leaky_oil_fissure.png"
	case "molten_aluminum":
		return "geyser_aluminum_volcano.png"
	case "molten_cobalt":
		return "geyser_cobalt_volcano.png"
	case "liquid_sulfur":
		return "geyser_liquid_sulfur_geyser.png"
	case "molten_tungsten":
		return "geyser_tungsten_volcano.png"
	case "molten_niobium":
		return "geyser_niobium_volcano.png"
	case "OilWell":
		return "geyser_oil_reservoir.png"
	default:
		return ""
	}
}

func iconForPOI(id string) string {
	id = simplifyID(id)
	switch id {
	case "Headquarters":
		return "building_printing_pod.png"
	case "WarpConduitSender":
		return "building_supply_teleporter_input.png"
	case "WarpConduitReceiver":
		return "building_supply_teleporter_output.png"
	case "WarpPortal":
		return "building_teleporter_transmitter.png"
	case "WarpReceiver":
		return "building_teleporter_receiver.png"
	case "GeneShuffler":
		return "building_neural_vacillator.png"
	case "MassiveHeatSink":
		return "building_anti_entropy_thermo_nullifier.png"
	case "SapTree":
		return "building_sap_tree.png"
	case "GravitasPedestal":
		return "poi_artifact_outline.png"
	case "PropSurfaceSatellite1":
		return "poi_crashed_satellite.png"
	case "PropSurfaceSatellite2":
		return "poi_wrecked_satellite.png"
	case "PropSurfaceSatellite3":
		return "poi_crushed_satellite.png"
	case "TemporalTearOpener":
		return "building_temporal_tear_opener.png"
	case "CryoTank":
		return "building_cryotank.png"
	case "PropFacilityStatue":
		return "poi_prop_facility_statue.png"
	case "GeothermalVentEntity":
		return "poi_geothermal_vent_entity.png"
	case "GeothermalControllerEntity":
		return "poi_geothermal_controller_entity.png"
	case "POICeresTechUnlock":
		return "poi_ceres_tech_unlock.png"
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
	id = simplifyID(id)
	if v, ok := names.Geysers[id]; ok {
		return v
	}
	return id
}

func displayPOI(id string) string {
	id = simplifyID(id)
	if v, ok := names.POIs[id]; ok {
		return v
	}
	return id
}

// formatLabel splits long names so they wrap every two words. It returns the
// formatted text and the length of the longest line for centering purposes.
func formatLabel(name string) (string, int) {
	words := strings.Fields(name)
	if len(words) <= 2 {
		return name, len(name)
	}
	var lines []string
	for i := 0; i < len(words); i += 2 {
		end := i + 2
		if end > len(words) {
			end = len(words)
		}
		line := strings.Join(words[i:end], " ")
		lines = append(lines, line)
	}
	width := 0
	for _, l := range lines {
		if len(l) > width {
			width = len(l)
		}
	}
	for i, l := range lines {
		if len(l) < width {
			pad := (width - len(l)) / 2
			lines[i] = strings.Repeat(" ", pad) + l
		}
	}
	formatted := strings.Join(lines, "\n")
	return formatted, width
}

var whitePixel = func() *ebiten.Image {
	img := ebiten.NewImage(1, 1)
	img.Fill(color.White)
	return img
}()

var tundraPattern = func() *ebiten.Image {
	const size = 32
	img := ebiten.NewImage(size, size)
	img.Fill(color.RGBA{0, 0, 0, 0})
	line := color.RGBA{255, 255, 255, 5}
	for i := 0; i < size; i++ {
		img.Set(i, size-1-i, line)
		if i+4 < size {
			img.Set(i, size-1-(i+4), line)
		}
	}
	return img
}()

var magmaPattern = func() *ebiten.Image {
	const size = 32
	img := ebiten.NewImage(size, size)
	img.Fill(color.RGBA{0, 0, 0, 0})
	lava := color.RGBA{255, 80, 0, 8}
	for x := 0; x < size; x++ {
		y := int((math.Sin(float64(x)/3) + 1) * float64(size) / 4)
		if y < size {
			img.Set(x, y, lava)
		}
		if y+8 < size {
			img.Set(x, y+8, lava)
		}
	}
	return img
}()

var oceanPattern = func() *ebiten.Image {
	const size = 32
	img := ebiten.NewImage(size, size)
	img.Fill(color.RGBA{0, 0, 0, 0})
	wave := color.RGBA{255, 255, 255, 4}
	for x := 0; x < size; x++ {
		y := int((math.Sin(float64(x)/2) + 1) * float64(size) / 4)
		if y < size {
			img.Set(x, y, wave)
		}
		if y+12 < size {
			img.Set(x, y+12, wave)
		}
	}
	return img
}()

var sandPattern = func() *ebiten.Image {
	const size = 32
	img := ebiten.NewImage(size, size)
	img.Fill(color.RGBA{0, 0, 0, 0})
	dot := color.RGBA{255, 255, 255, 3}
	for x := 0; x < size; x++ {
		for y := 0; y < size; y++ {
			if (x+y)%16 == 0 {
				img.Set(x, y, dot)
			}
		}
	}
	return img
}()

var toxicPattern = func() *ebiten.Image {
	const size = 32
	img := ebiten.NewImage(size, size)
	img.Fill(color.RGBA{0, 0, 0, 0})
	slime := color.RGBA{100, 255, 100, 4}
	cx, cy := float64(size)/2, float64(size)/2
	for x := 0; x < size; x++ {
		for y := 0; y < size; y++ {
			dx := float64(x) - cx
			dy := float64(y) - cy
			if dx*dx+dy*dy < float64(size) {
				img.Set(x, y, slime)
			}
		}
	}
	return img
}()

var oilPattern = func() *ebiten.Image {
	const size = 32
	img := ebiten.NewImage(size, size)
	img.Fill(color.RGBA{0, 0, 0, 0})
	line := color.RGBA{255, 255, 255, 4}
	for x := 0; x < size; x++ {
		y := int((math.Sin(float64(x)/2) + 1) * float64(size) / 4)
		if y < size {
			img.Set(x, y, line)
		}
	}
	return img
}()

var marshPattern = func() *ebiten.Image {
	const size = 32
	img := ebiten.NewImage(size, size)
	img.Fill(color.RGBA{0, 0, 0, 0})
	grass := color.RGBA{255, 255, 255, 3}
	for x := 0; x < size; x += 10 {
		for y := 0; y < size; y++ {
			img.Set(x, y, grass)
		}
	}
	return img
}()

func drawTextWithBG(dst *ebiten.Image, text string, x, y int) {
	lines := strings.Split(text, "\n")
	width := 0
	for _, l := range lines {
		if len(l) > width {
			width = len(l)
		}
	}
	height := len(lines) * 16
	vector.DrawFilledRect(dst, float32(x-2), float32(y-2), float32(width*6+4), float32(height+4), color.RGBA{0, 0, 0, 128}, false)
	ebitenutil.DebugPrintAt(dst, text, x, y)
}

func drawTextWithBGBorder(dst *ebiten.Image, text string, x, y int, border color.Color) {
	lines := strings.Split(text, "\n")
	width := 0
	for _, l := range lines {
		if len(l) > width {
			width = len(l)
		}
	}
	height := len(lines) * 16
	bx := x - 2
	by := y - 2
	bw := width*6 + 4
	bh := height + 4
	vector.DrawFilledRect(dst, float32(bx-1), float32(by-1), float32(bw+2), float32(bh+2), border, false)
	vector.DrawFilledRect(dst, float32(bx), float32(by), float32(bw), float32(bh), color.RGBA{0, 0, 0, 128}, false)
	ebitenutil.DebugPrintAt(dst, text, x, y)
}

// uniqueColor generates a visually distinct color for the given index.
func uniqueColor(index int) color.RGBA {
	h := float64((index * 137) % 360)
	// Increase the oscillation period so the palette offers more
	// distinct options before repeating. With a period of 30 we can
	// generate a wider range of saturation and lightness values.
	period := 30.0
	s := 0.5 + 0.4*math.Sin(float64(index)*2*math.Pi/period)
	l := 0.45 + 0.3*math.Cos(float64(index)*2*math.Pi/period)
	if s < 0.3 {
		s = 0.3
	} else if s > 0.9 {
		s = 0.9
	}
	if l < 0.3 {
		l = 0.3
	} else if l > 0.8 {
		l = 0.8
	}
	c := (1 - math.Abs(2*l-1)) * s
	x := c * (1 - math.Abs(math.Mod(h/60.0, 2)-1))
	m := l - c/2
	var r1, g1, b1 float64
	switch {
	case h < 60:
		r1, g1, b1 = c, x, 0
	case h < 120:
		r1, g1, b1 = x, c, 0
	case h < 180:
		r1, g1, b1 = 0, c, x
	case h < 240:
		r1, g1, b1 = 0, x, c
	case h < 300:
		r1, g1, b1 = x, 0, c
	default:
		r1, g1, b1 = c, 0, x
	}
	r := uint8(math.Round((r1 + m) * 255))
	g := uint8(math.Round((g1 + m) * 255))
	b := uint8(math.Round((b1 + m) * 255))
	return color.RGBA{r, g, b, 255}
}

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
		vs[i].DstX = float32(x)
		vs[i].DstY = float32(y)
		vs[i].SrcX = 0
		vs[i].SrcY = 0
		vs[i].ColorR = float32(r) / 0xffff
		vs[i].ColorG = float32(g) / 0xffff
		vs[i].ColorB = float32(b) / 0xffff
		vs[i].ColorA = float32(a) / 0xffff
	}
	op := &ebiten.DrawTrianglesOptions{
		AntiAlias:      true,
		ColorScaleMode: ebiten.ColorScaleModePremultipliedAlpha,
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
		vs[i].DstX = float32(x)
		vs[i].DstY = float32(y)
		vs[i].SrcX = 0
		vs[i].SrcY = 0
		vs[i].ColorR = float32(r) / 0xffff
		vs[i].ColorG = float32(g) / 0xffff
		vs[i].ColorB = float32(b) / 0xffff
		vs[i].ColorA = float32(a) / 0xffff
	}
	op := &ebiten.DrawTrianglesOptions{
		AntiAlias:      true,
		ColorScaleMode: ebiten.ColorScaleModePremultipliedAlpha,
		FillRule:       ebiten.FillRuleEvenOdd,
	}
	dst.DrawTriangles(vs, is, whitePixel, op)
}

func drawTundraGradient(dst *ebiten.Image, polys [][]Point, camX, camY, zoom float64) {
	if len(polys) == 0 {
		return
	}
	minY, maxY := math.MaxFloat64, -math.MaxFloat64
	for _, pts := range polys {
		for _, pt := range pts {
			y := float64(pt.Y*2)*zoom + camY
			if y < minY {
				minY = y
			}
			if y > maxY {
				maxY = y
			}
		}
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
	for i := range vs {
		x := float64(vs[i].DstX)*zoom + camX
		y := float64(vs[i].DstY)*zoom + camY
		vs[i].DstX = float32(math.Round(x))
		vs[i].DstY = float32(math.Round(y))
		vs[i].SrcX = 0
		vs[i].SrcY = 0
		alpha := float32(0.0)
		if maxY > minY {
			alpha = float32(0.15 * (1 - (y-minY)/(maxY-minY)))
		}
		vs[i].ColorR = 1
		vs[i].ColorG = 1
		vs[i].ColorB = 1
		vs[i].ColorA = alpha
	}
	op := &ebiten.DrawTrianglesOptions{
		AntiAlias:      true,
		ColorScaleMode: ebiten.ColorScaleModeStraightAlpha,
		FillRule:       ebiten.FillRuleEvenOdd,
	}
	dst.DrawTriangles(vs, is, whitePixel, op)
}

func drawPattern(dst *ebiten.Image, polys [][]Point, camX, camY, zoom float64, pattern *ebiten.Image) {
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
	for i := range vs {
		worldX := float64(vs[i].DstX)
		worldY := float64(vs[i].DstY)
		x := worldX*zoom + camX
		y := worldY*zoom + camY
		vs[i].DstX = float32(math.Round(x))
		vs[i].DstY = float32(math.Round(y))
		vs[i].SrcX = float32(worldX * zoom)
		vs[i].SrcY = float32(worldY * zoom)
		vs[i].ColorR = 1
		vs[i].ColorG = 1
		vs[i].ColorB = 1
		vs[i].ColorA = 1
	}
	op := &ebiten.DrawTrianglesOptions{
		AntiAlias:      true,
		ColorScaleMode: ebiten.ColorScaleModeStraightAlpha,
		FillRule:       ebiten.FillRuleEvenOdd,
		Address:        ebiten.AddressRepeat,
		Filter:         ebiten.FilterLinear,
	}
	dst.DrawTriangles(vs, is, pattern, op)
}

func drawBiomeOutline(dst *ebiten.Image, polys [][]Point, camX, camY, zoom float64) {
	for _, pts := range polys {
		if len(pts) < 2 {
			continue
		}
		for i := 0; i < len(pts); i++ {
			a := pts[i]
			b := pts[(i+1)%len(pts)]
			x0 := float32(math.Round(float64(a.X*2)*zoom + camX))
			y0 := float32(math.Round(float64(a.Y*2)*zoom + camY))
			x1 := float32(math.Round(float64(b.X*2)*zoom + camX))
			y1 := float32(math.Round(float64(b.Y*2)*zoom + camY))
			vector.StrokeLine(dst, x0, y0, x1, y1, 1, color.RGBA{255, 255, 255, 128}, true)
		}
	}
}

// buildLegendImage returns an image of biome colors with a single background.
func buildLegendImage(biomes []BiomePath) *ebiten.Image {
	set := make(map[string]struct{})
	for _, b := range biomes {
		set[b.Name] = struct{}{}
	}
	names := make([]string, 0, len(set))
	for n := range set {
		names = append(names, n)
	}
	sort.Strings(names)

	maxLen := 0
	for _, name := range names {
		if l := len(displayBiome(name)); l > maxLen {
			maxLen = l
		}
	}

	width := 30 + maxLen*6 + 5
	height := LegendRowSpacing*(len(names)+1) + 7

	img := ebiten.NewImage(width, height)
	img.Fill(color.RGBA{0, 0, 0, 77})

	y := 10
	drawTextWithBG(img, "Biomes", 5, y)
	y += LegendRowSpacing
	for _, name := range names {
		clr, ok := biomeColors[name]
		if !ok {
			clr = color.RGBA{60, 60, 60, 255}
		}
		vector.DrawFilledRect(img, 5, float32(y), 20, 10, clr, false)
		drawTextWithBGBorder(img, displayBiome(name), 30, y, clr)
		y += LegendRowSpacing
	}

	return img
}

// initObjectLegend prepares the mapping of object names to numeric labels and
// caches the legend text.
func (g *Game) initObjectLegend() {
	if g.legendMap != nil {
		return
	}
	g.legendMap = make(map[string]int)
	g.legendColors = nil
	counter := 1
	for _, gy := range g.geysers {
		name := displayGeyser(gy.ID)
		if _, ok := g.legendMap["g"+name]; !ok {
			g.legendMap["g"+name] = counter
			g.legendEntries = append(g.legendEntries, fmt.Sprintf("%d: %s", counter, name))
			g.legendColors = append(g.legendColors, uniqueColor(counter))
			counter++
		}
	}
	for _, poi := range g.pois {
		name := displayPOI(poi.ID)
		if _, ok := g.legendMap["p"+name]; !ok {
			g.legendMap["p"+name] = counter
			g.legendEntries = append(g.legendEntries, fmt.Sprintf("%d: %s", counter, name))
			g.legendColors = append(g.legendColors, uniqueColor(counter))
			counter++
		}
	}
}

// drawNumberLegend draws the cached legend image on the top right corner with
// a single semi-transparent background.
func (g *Game) drawNumberLegend(dst *ebiten.Image) {
	if len(g.legendEntries) == 0 {
		return
	}
	if g.legendImage == nil {
		maxLen := 0
		for _, e := range g.legendEntries {
			if len(e) > maxLen {
				maxLen = len(e)
			}
		}
		width := maxLen*6 + 10
		height := LegendRowSpacing*(len(g.legendEntries)+1) + 7
		img := ebiten.NewImage(width, height)
		img.Fill(color.RGBA{0, 0, 0, 77})
		y := 10
		drawTextWithBG(img, "Items", 5, y)
		y += LegendRowSpacing
		for i, e := range g.legendEntries {
			clr := color.RGBA{}
			if i < len(g.legendColors) {
				clr = g.legendColors[i]
			}
			drawTextWithBGBorder(img, e, 5, y, clr)
			y += LegendRowSpacing
		}
		g.legendImage = img
	}
	scale := 1.0
	if g.height > 850 {
		scale = 2.0
	}
	w := float64(g.legendImage.Bounds().Dx()) * scale
	x := float64(dst.Bounds().Dx()) - w - 12
	y := 10.0
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(scale, scale)
	op.GeoM.Translate(x, y)
	dst.DrawImage(g.legendImage, op)
}

// Game implements ebiten.Game and displays geysers with their names.
type Game struct {
	geysers        []Geyser
	pois           []PointOfInterest
	biomes         []BiomePath
	icons          map[string]*ebiten.Image
	width          int
	height         int
	astWidth       int
	astHeight      int
	camX           float64
	camY           float64
	zoom           float64
	dragging       bool
	lastX          int
	lastY          int
	touches        map[ebiten.TouchID]touchPoint
	pinchDist      float64
	needsRedraw    bool
	screenshotPath string
	captured       bool
	legend         *ebiten.Image
	legendMap      map[string]int
	legendEntries  []string
	legendColors   []color.RGBA
	legendImage    *ebiten.Image
	showHelp       bool
	lastWheel      time.Time
	loading        bool
	status         string
	iconResults    chan loadedIcon
	coord          string
	mobile         bool
}

type label struct {
	text string
	x    int
	y    int
	clr  color.RGBA
}

type touchPoint struct {
	x int
	y int
}

type loadedIcon struct {
	name string
	img  *ebiten.Image
}

func (g *Game) helpRect() image.Rectangle {
	size := HelpIconSize
	x := g.width - size - HelpMargin
	y := g.height - size - HelpMargin
	return image.Rect(x, y, x+size, y+size)
}

func (g *Game) clampCamera() {
	w := float64(g.astWidth) * 2 * g.zoom
	h := float64(g.astHeight) * 2 * g.zoom
	margin := float64(CameraMargin)
	if g.camX < -w-margin {
		g.camX = -w - margin
	}
	if g.camX > float64(g.width)+margin {
		g.camX = float64(g.width) + margin
	}
	if g.camY < -h-margin {
		g.camY = -h - margin
	}
	if g.camY > float64(g.height)+margin {
		g.camY = float64(g.height) + margin
	}
}

func (g *Game) startIconLoader(names []string) {
	g.iconResults = make(chan loadedIcon, len(names))
	go func() {
		for _, n := range names {
			img, _ := loadImageFile(n)
			g.iconResults <- loadedIcon{name: n, img: img}
		}
		close(g.iconResults)
	}()
}

func (g *Game) Update() error {
	const panSpeed = PanSpeed

iconsLoop:
	for g.iconResults != nil {
		select {
		case li, ok := <-g.iconResults:
			if !ok {
				g.iconResults = nil
				continue
			}
			if g.icons != nil {
				g.icons[li.name] = li.img
			}
			g.needsRedraw = true
		default:
			break iconsLoop
		}
	}

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

	// Touch gestures
	touchIDs := ebiten.TouchIDs()
	switch len(touchIDs) {
	case 1:
		id := touchIDs[0]
		x, y := ebiten.TouchPosition(id)
		if last, ok := g.touches[id]; ok {
			g.camX += float64(x - last.x)
			g.camY += float64(y - last.y)
		}
		g.touches = map[ebiten.TouchID]touchPoint{id: {x: x, y: y}}
		g.pinchDist = 0
	case 2:
		id1, id2 := touchIDs[0], touchIDs[1]
		x1, y1 := ebiten.TouchPosition(id1)
		x2, y2 := ebiten.TouchPosition(id2)
		dx := float64(x2 - x1)
		dy := float64(y2 - y1)
		dist := math.Hypot(dx, dy)
		midX := float64(x1+x2) / 2
		midY := float64(y1+y2) / 2
		if g.pinchDist != 0 {
			factor := dist / g.pinchDist
			oldZoom := g.zoom
			g.zoom *= factor
			if g.zoom < MinZoom {
				g.zoom = MinZoom
			}
			worldX := (midX - g.camX) / oldZoom
			worldY := (midY - g.camY) / oldZoom
			g.camX = midX - worldX*g.zoom
			g.camY = midY - worldY*g.zoom
		}
		g.pinchDist = dist
		g.touches = map[ebiten.TouchID]touchPoint{
			id1: {x: x1, y: y1},
			id2: {x: x2, y: y2},
		}
	default:
		g.touches = nil
		g.pinchDist = 0
	}

	// Zoom with keyboard
	zoomFactor := 1.0
	if ebiten.IsKeyPressed(ebiten.KeyEqual) || ebiten.IsKeyPressed(ebiten.KeyKPAdd) {
		zoomFactor *= KeyZoomFactor
	}
	if ebiten.IsKeyPressed(ebiten.KeyMinus) || ebiten.IsKeyPressed(ebiten.KeyKPSubtract) {
		zoomFactor /= KeyZoomFactor
	}

	// Zoom with mouse wheel
	_, wheelY := ebiten.Wheel()
	if runtime.GOARCH == "wasm" && wheelY != 0 {
		now := time.Now()
		if now.Sub(g.lastWheel) < WheelThrottle {
			wheelY = 0
		} else {
			g.lastWheel = now
		}
	}
	if wheelY != 0 {
		if wheelY > 0 {
			zoomFactor *= WheelZoomFactor
		} else {
			zoomFactor /= WheelZoomFactor
		}
	}

	if zoomFactor != 1.0 {
		oldZoom := g.zoom
		g.zoom *= zoomFactor
		if g.zoom < MinZoom {
			g.zoom = MinZoom
		}
		cx, cy := float64(g.width)/2, float64(g.height)/2
		worldX := (cx - g.camX) / oldZoom
		worldY := (cy - g.camY) / oldZoom
		g.camX = cx - worldX*g.zoom
		g.camY = cy - worldY*g.zoom
	}

	mx, my := ebiten.CursorPosition()
	if g.showHelp {
		if !g.helpRect().Overlaps(image.Rect(mx, my, mx+1, my+1)) {
			g.showHelp = false
			g.needsRedraw = true
		}
	} else if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) && g.helpRect().Overlaps(image.Rect(mx, my, mx+1, my+1)) {
		g.showHelp = true
		g.needsRedraw = true
	}

	g.clampCamera()

	if g.camX != oldX || g.camY != oldY || g.zoom != oldZoom {
		g.needsRedraw = true
	}

	if g.captured {
		return ebiten.Termination
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	if g.loading || (len(g.biomes) == 0 && g.status != "") {
		screen.Fill(color.RGBA{30, 30, 30, 255})
		msg := g.status
		if msg == "" {
			msg = "Fetching..."
		}
		x := g.width/2 - len(msg)*LabelCharWidth/2
		y := g.height / 2
		drawTextWithBG(screen, msg, x, y)
		return
	}
	if g.needsRedraw {
		screen.Fill(color.RGBA{30, 30, 30, 255})
		if clr, ok := biomeColors["Space"]; ok {
			vector.DrawFilledRect(screen, float32(g.camX), float32(g.camY),
				float32(float64(g.astWidth)*2*g.zoom),
				float32(float64(g.astHeight)*2*g.zoom), clr, false)
		}
		labels := []label{}
		useNumbers := !g.mobile && g.zoom < LegendZoomThreshold
		if g.legendMap == nil {
			g.initObjectLegend()
		}

		for _, bp := range g.biomes {
			clr, ok := biomeColors[bp.Name]
			if !ok {
				clr = color.RGBA{60, 60, 60, 255}
			}
			drawBiome(screen, bp.Polygons, clr, g.camX, g.camY, g.zoom)
			/*
				switch bp.Name {
				case "FrozenWastes", "IceCaves":
					drawTundraGradient(screen, bp.Polygons, g.camX, g.camY, g.zoom)
					drawPattern(screen, bp.Polygons, g.camX, g.camY, g.zoom, tundraPattern)
				case "MagmaCore":
					drawPattern(screen, bp.Polygons, g.camX, g.camY, g.zoom, magmaPattern)
				case "Ocean":
					drawPattern(screen, bp.Polygons, g.camX, g.camY, g.zoom, oceanPattern)
				case "Sandstone":
					drawPattern(screen, bp.Polygons, g.camX, g.camY, g.zoom, sandPattern)
				case "ToxicJungle":
					drawPattern(screen, bp.Polygons, g.camX, g.camY, g.zoom, toxicPattern)
				case "OilField":
					drawPattern(screen, bp.Polygons, g.camX, g.camY, g.zoom, oilPattern)
				case "BoggyMarsh":
					drawPattern(screen, bp.Polygons, g.camX, g.camY, g.zoom, marshPattern)
				} */
			drawBiomeOutline(screen, bp.Polygons, g.camX, g.camY, g.zoom)
		}

		for _, gy := range g.geysers {
			x := math.Round((float64(gy.X) * 2 * g.zoom) + g.camX)
			y := math.Round((float64(gy.Y) * 2 * g.zoom) + g.camY)

			name := displayGeyser(gy.ID)
			formatted, width := formatLabel(name)
			clr := color.RGBA{}
			if idx, ok := g.legendMap["g"+name]; ok && idx-1 < len(g.legendColors) {
				clr = g.legendColors[idx-1]
			}

			if iconName := iconForGeyser(gy.ID); iconName != "" {
				if img, ok := g.icons[iconName]; ok && img != nil {
					op := &ebiten.DrawImageOptions{Filter: ebiten.FilterLinear}
					op.GeoM.Scale(g.zoom*IconScale, g.zoom*IconScale)
					w := float64(img.Bounds().Dx()) * g.zoom * IconScale
					h := float64(img.Bounds().Dy()) * g.zoom * IconScale
					op.GeoM.Translate(x-w/2, y-h/2)
					screen.DrawImage(img, op)
					if useNumbers {
						formatted = strconv.Itoa(g.legendMap["g"+name])
						width = len(formatted)
					}
					labels = append(labels, label{formatted, int(x) - (width*LabelCharWidth)/2, int(y+h/2) + 2, clr})
					continue
				}
			}

			vector.DrawFilledRect(screen, float32(x-2), float32(y-2), 4, 4, clr, true)
			if useNumbers {
				formatted = strconv.Itoa(g.legendMap["g"+name])
				width = len(formatted)
			}
			labels = append(labels, label{formatted, int(x) - (width*LabelCharWidth)/2, int(y) + 4, clr})
		}
		for _, poi := range g.pois {
			x := math.Round((float64(poi.X) * 2 * g.zoom) + g.camX)
			y := math.Round((float64(poi.Y) * 2 * g.zoom) + g.camY)

			name := displayPOI(poi.ID)
			formatted, width := formatLabel(name)
			clr := color.RGBA{}
			if idx, ok := g.legendMap["p"+name]; ok && idx-1 < len(g.legendColors) {
				clr = g.legendColors[idx-1]
			}

			if iconName := iconForPOI(poi.ID); iconName != "" {
				if img, ok := g.icons[iconName]; ok && img != nil {
					op := &ebiten.DrawImageOptions{Filter: ebiten.FilterLinear}
					op.GeoM.Scale(g.zoom*IconScale, g.zoom*IconScale)
					w := float64(img.Bounds().Dx()) * g.zoom * IconScale
					h := float64(img.Bounds().Dy()) * g.zoom * IconScale
					op.GeoM.Translate(x-w/2, y-h/2)
					screen.DrawImage(img, op)
					if useNumbers {
						formatted = strconv.Itoa(g.legendMap["p"+name])
						width = len(formatted)
					}
					labels = append(labels, label{formatted, int(x) - (width*LabelCharWidth)/2, int(y+h/2) + 2, clr})
					continue
				}
			}

			vector.DrawFilledRect(screen, float32(x-2), float32(y-2), 4, 4, clr, true)
			if useNumbers {
				formatted = strconv.Itoa(g.legendMap["p"+name])
				width = len(formatted)
			}
			labels = append(labels, label{formatted, int(x) - (width*LabelCharWidth)/2, int(y) + 4, clr})
		}

		if g.legend == nil {
			g.legend = buildLegendImage(g.biomes)
		}
		opLegend := &ebiten.DrawImageOptions{}
		if g.height > 850 {
			opLegend.GeoM.Scale(2, 2)
		}
		screen.DrawImage(g.legend, opLegend)
		for _, l := range labels {
			if l.clr.A != 0 {
				drawTextWithBGBorder(screen, l.text, l.x, l.y, l.clr)
			} else {
				drawTextWithBG(screen, l.text, l.x, l.y)
			}
		}
		if useNumbers {
			g.drawNumberLegend(screen)
		}

		if g.coord != "" {
			x := g.width/2 - len(g.coord)*LabelCharWidth/2
			drawTextWithBG(screen, g.coord, x, 10)
		}

		// Draw help icon
		hr := g.helpRect()
		cx := float32(hr.Min.X + HelpIconSize/2)
		cy := float32(hr.Min.Y + HelpIconSize/2)
		vector.DrawFilledCircle(screen, cx, cy, HelpIconSize/2, color.RGBA{0, 0, 0, 180}, true)
		ebitenutil.DebugPrintAt(screen, "?", hr.Min.X+7, hr.Min.Y+5)
		if g.showHelp {
			tx := hr.Min.X - 170
			ty := hr.Min.Y - 70
			drawTextWithBG(screen, helpMessage, tx, ty)
		}

		g.needsRedraw = false
	}
	if g.screenshotPath != "" && !g.captured {
		b := screen.Bounds()
		pixels := make([]byte, 4*b.Dx()*b.Dy())
		screen.ReadPixels(pixels)
		img := &image.RGBA{Pix: pixels, Stride: 4 * b.Dx(), Rect: b}
		if f, err := os.Create(g.screenshotPath); err == nil {
			_ = png.Encode(f, img)
			f.Close()
		}
		g.captured = true
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
		g.clampCamera()

		g.needsRedraw = true
	}
	return outsideWidth, outsideHeight
}

func main() {
	coord := flag.String("coord", "SNDST-A-7-0-0-0", "seed coordinate")
	out := flag.String("out", "", "optional path to save JSON")
	screenshot := flag.String("screenshot", "", "path to save a PNG screenshot and exit")
	flag.Parse()

	game := &Game{
		icons:   make(map[string]*ebiten.Image),
		width:   DefaultWidth,
		height:  DefaultHeight,
		zoom:    1.0,
		loading: true,
		status:  "Fetching...",
		coord:   *coord,
		mobile:  isMobile(),
	}
	go func() {
		fmt.Println("Fetching:", *coord)
		cborData, err := fetchSeedCBOR(*coord)
		if err != nil {
			game.status = "Error: " + err.Error()
			game.needsRedraw = true
			game.loading = false
			return
		}
		seed, err := decodeSeed(cborData)
		if err != nil {
			game.status = "Error: " + err.Error()
			game.needsRedraw = true
			game.loading = false
			return
		}
		if *out != "" {
			jsonData, _ := json.MarshalIndent(seed, "", "  ")
			_ = saveToFile(*out, jsonData)
		}
		ast := seed.Asteroids[0]
		bps := parseBiomePaths(ast.BiomePaths)
		game.geysers = ast.Geysers
		game.pois = ast.POIs
		game.biomes = bps
		game.astWidth = ast.SizeX
		game.astHeight = ast.SizeY
		game.legend = buildLegendImage(bps)
		game.camX = (float64(game.width) - float64(game.astWidth)*2*game.zoom) / 2
		game.camY = (float64(game.height) - float64(game.astHeight)*2*game.zoom) / 2
		game.clampCamera()
		names := []string{}
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
	}()
	if *screenshot != "" {
		game.screenshotPath = *screenshot
	}
	ebiten.SetWindowSize(game.width, game.height)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetWindowTitle("Geysers - " + *coord)
	ebiten.SetScreenClearedEveryFrame(false)
	ebiten.SetVsyncEnabled(true)
	if err := ebiten.RunGame(game); err != nil {
		fmt.Println("Error running game:", err)
	}
}
