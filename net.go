package main

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"google.golang.org/protobuf/proto"
	seedpb "oni-view/data/pb"
)

var seedProtoBaseURL = ProtoBaseURL

// fetchSeedJSON retrieves the seed data in JSON format for a given coordinate.
// It first tries the primary URL and falls back to the secondary if needed.
func fetchSeedJSON(coordinate string) ([]byte, error) {
	urls := []string{
		BaseURL + coordinate,
		FallbackBaseURL + coordinate,
	}

	var lastErr error
	for _, url := range urls {
		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Set("Accept", AcceptJSONHeader)

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			lastErr = fmt.Errorf("request failed: %v", err)
			continue
		}
		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			lastErr = fmt.Errorf("read failed: %v", err)
			continue
		}
		if resp.StatusCode != http.StatusOK {
			lastErr = fmt.Errorf("unexpected status %d: %s", resp.StatusCode, body)
			continue
		}
		return body, nil
	}
	return nil, lastErr
}

// fetchSeedProto retrieves the seed data in protobuf format for a given coordinate.
// It requests the protobuf endpoint and transparently decompresses gzip-encoded responses.
func fetchSeedProto(coordinate string) ([]byte, error) {
	url := seedProtoBaseURL + coordinate
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

// decodeSeed parses the JSON seed data into SeedData.
func decodeSeed(jsonData []byte) (*SeedData, error) {
	var seed SeedData
	if err := json.Unmarshal(jsonData, &seed); err != nil {
		return nil, fmt.Errorf("JSON decode failed: %v", err)
	}
	return &seed, nil
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
			ID:    a.Id,
			SizeX: int(a.SizeX),
			SizeY: int(a.SizeY),
		}
		for _, g := range a.Geysers {
			ast.Geysers = append(ast.Geysers, Geyser{
				ID:             g.Id,
				X:              int(g.X),
				Y:              int(g.Y),
				ActiveCycles:   g.ActiveCycles,
				AvgEmitRate:    g.AvgEmitRate,
				DormancyCycles: g.DormancyCycles,
				EmitRate:       g.EmitRate,
				EruptionTime:   g.EruptionTime,
				IdleTime:       g.IdleTime,
			})
		}
		for _, p := range a.PointsOfInterest {
			ast.POIs = append(ast.POIs, PointOfInterest{
				ID: p.Id,
				X:  int(p.X),
				Y:  int(p.Y),
			})
		}
		ast.BiomePaths = convertBiomePaths(a.BiomePaths)
		seed.Asteroids = append(seed.Asteroids, ast)
	}
	return seed, nil
}

func convertBiomePaths(bp *seedpb.BiomePathsCompact) BiomePathsCompact {
	if bp == nil {
		return BiomePathsCompact{}
	}
	var paths []BiomePath
	for _, p := range bp.Paths {
		var polys [][]Point
		for _, poly := range p.Polygons {
			var pts []Point
			for _, pt := range poly.Points {
				pts = append(pts, Point{X: int(pt.X), Y: int(pt.Y)})
			}
			polys = append(polys, pts)
		}
		paths = append(paths, BiomePath{Name: p.Name, Polygons: polys})
	}
	return BiomePathsCompact{Paths: paths}
}
