package main

import (
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"google.golang.org/protobuf/proto"
	seedpb "oni-view/data/pb"
)

var seedProtoBaseURL = ProtoBaseURL

// geyserTypeFromID maps numeric geyser IDs to their string descriptors.
func geyserTypeFromID(id int32) string {
	switch id {
	case 0:
		return "steam"
	case 1:
		return "hot_hydrogen"
	case 2:
		return "methane"
	case 3:
		return "chlorine_gas"
	case 4:
		return "chlorine_gas_cool"
	case 5:
		return "hot_steam"
	case 6:
		return "hot_co2"
	case 7:
		return "hot_po2"
	case 8:
		return "slimy_po2"
	case 9:
		return "hot_water"
	case 10:
		return "slush_water"
	case 11:
		return "filthy_water"
	case 12:
		return "slush_salt_water"
	case 13:
		return "salt_water"
	case 14:
		return "liquid_co2"
	case 15:
		return "oil_drip"
	case 16:
		return "liquid_sulfur"
	case 17:
		return "molten_iron"
	case 18:
		return "molten_copper"
	case 19:
		return "molten_gold"
	case 20:
		return "molten_aluminum"
	case 21:
		return "molten_cobalt"
	case 22:
		return "molten_tungsten"
	case 23:
		return "molten_niobium"
	case 24:
		return "big_volcano"
	case 25:
		return "small_volcano"
	case 26:
		return "OilWell"
	default:
		return ""
	}
}

// fetchSeedProto retrieves the seed data in protobuf format for a given coordinate.
// It requests the protobuf endpoint and transparently decompresses gzip-encoded responses.
func fetchSeedProto(coordinate string) ([]byte, error) {
	base := strings.TrimSuffix(seedProtoBaseURL, "/")
	url := base + "/" + coordinate
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Accept", AcceptProtoHeader)
	req.Header.Set("Accept-Encoding", GzipEncoding)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %v", err)
	}
	defer resp.Body.Close()

	var reader io.Reader = resp.Body
	if resp.Header.Get("Content-Encoding") == GzipEncoding {
		gz, err := gzip.NewReader(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("gzip init failed: %v", err)
		}
		defer gz.Close()
		reader = gz
	}
	body, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("read failed: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status %d: %s", resp.StatusCode, body)
	}
	return body, nil
}

// decodeSeedProto parses the protobuf seed data into SeedData.
func decodeSeedProto(protoData []byte) (*SeedData, error) {
	var pb seedpb.Cluster
	if err := proto.Unmarshal(protoData, &pb); err != nil {
		return nil, fmt.Errorf("protobuf decode failed: %v", err)
	}
	seed := &SeedData{}
	for _, a := range pb.Asteroids {
		ast := Asteroid{
			ID:    asteroidNameFromID(a.Id),
			SizeX: int(a.SizeX),
			SizeY: int(a.SizeY),
		}
		for _, g := range a.Geysers {
			ast.Geysers = append(ast.Geysers, Geyser{
				ID:             geyserTypeFromID(g.Id),
				X:              int(g.X),
				Y:              int(g.Y),
				EmitRate:       float64(g.EmitRate),
				AvgEmitRate:    float64(g.AvgEmitRate),
				IdleTime:       float64(g.IdleTime),
				EruptionTime:   float64(g.EruptionTime),
				DormancyCycles: float64(g.DormancyCyclesRounded),
				ActiveCycles:   float64(g.ActiveCyclesRounded),
			})
		}
		for _, p := range a.PointsOfInterest {
			ast.POIs = append(ast.POIs, PointOfInterest{
				ID: p.Id.String(),
				X:  int(p.X),
				Y:  int(p.Y),
			})
		}
		ast.BiomePaths = parseBiomePaths(a.BiomePaths)
		seed.Asteroids = append(seed.Asteroids, ast)
	}
	return seed, nil
}

// parseBiomePaths decodes the compact biome path string format into a
// BiomePathsCompact structure. The format uses newline-separated zone entries
// where each line contains a numeric zone ID followed by colon-separated delta
// encoded coordinate pairs for each polygon.
func parseBiomePaths(s string) BiomePathsCompact {
	if s == "" {
		return BiomePathsCompact{}
	}
	zoneNames := map[int]string{
		0:  "FrozenWastes",
		1:  "CrystalCaverns",
		2:  "BoggyMarsh",
		3:  "Sandstone",
		4:  "ToxicJungle",
		5:  "MagmaCore",
		6:  "OilField",
		7:  "Space",
		8:  "Ocean",
		9:  "Rust",
		10: "Forest",
		11: "Radioactive",
		12: "Swamp",
		13: "Wasteland",
		15: "Metallic",
		16: "Barren",
		17: "Moo",
		18: "IceCaves",
		19: "CarrotQuarry",
		20: "SugarWoods",
		21: "PrehistoricGarden",
		22: "PrehistoricRaptor",
		23: "PrehistoricWetlands",
	}
	s = strings.ReplaceAll(s, "\\n", "\n")
	lines := strings.Split(s, "\n")
	var paths []BiomePath
	for _, line := range lines {
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			continue
		}
		zoneID, err := strconv.Atoi(parts[0])
		if err != nil {
			continue
		}
		zoneName, ok := zoneNames[zoneID]
		if !ok {
			zoneName = parts[0]
		}
		var polys [][]Point
		for _, polyStr := range strings.Split(parts[1], "|") {
			nums := strings.Fields(polyStr)
			if len(nums)%2 != 0 {
				continue
			}
			var pts []Point
			prevX, prevY := 0, 0
			for i := 0; i < len(nums); i += 2 {
				dx, err1 := strconv.Atoi(nums[i])
				dy, err2 := strconv.Atoi(nums[i+1])
				if err1 != nil || err2 != nil {
					continue
				}
				x := prevX + dx
				y := prevY + dy
				pts = append(pts, Point{X: x, Y: y})
				prevX = x
				prevY = y
			}
			polys = append(polys, pts)
		}
		paths = append(paths, BiomePath{Name: zoneName, Polygons: polys})
	}
	return BiomePathsCompact{Paths: paths}
}
