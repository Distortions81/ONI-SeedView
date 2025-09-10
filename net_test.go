package main

import (
	"compress/gzip"
	"net/http"
	"net/http/httptest"
	"testing"

	"google.golang.org/protobuf/proto"

	"oni-view/seedproto"
)

// TestFetchSeedProtoDecompressesGzip verifies that fetchSeedProto handles gzip responses.
func TestFetchSeedProtoDecompressesGzip(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		if got := r.Header.Get("Accept"); got != AcceptProtoHeader {
			t.Fatalf("unexpected Accept header: %s", got)
		}
		w.Header().Set("Content-Encoding", GzipEncoding)
		gz := gzip.NewWriter(w)
		gz.Write([]byte("hello"))
		gz.Close()
	}
	srv := httptest.NewServer(http.HandlerFunc(handler))
	defer srv.Close()

	old := seedProtoBaseURL
	seedProtoBaseURL = srv.URL + "/"
	defer func() { seedProtoBaseURL = old }()

	body, err := fetchSeedProto("test")
	if err != nil {
		t.Fatalf("fetchSeedProto error: %v", err)
	}
	if string(body) != "hello" {
		t.Fatalf("unexpected body: %s", body)
	}
}

// TestDecodeSeedProto verifies protobuf decoding into SeedData.
func TestDecodeSeedProto(t *testing.T) {
	pb := &seedproto.SeedDataProto{
		Asteroids: []*seedproto.AsteroidProto{
			{
				Id:    "A1",
				SizeX: 10,
				SizeY: 20,
				Geysers: []*seedproto.GeyserProto{
					{Id: "g1", X: 1, Y: 2},
				},
				PointsOfInterest: []*seedproto.PointOfInterestProto{
					{Id: "p1", X: 3, Y: 4},
				},
				BiomePaths: "bp",
			},
		},
	}
	data, err := proto.Marshal(pb)
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}
	seed, err := decodeSeedProto(data)
	if err != nil {
		t.Fatalf("decodeSeedProto error: %v", err)
	}
	if len(seed.Asteroids) != 1 {
		t.Fatalf("expected 1 asteroid, got %d", len(seed.Asteroids))
	}
	a := seed.Asteroids[0]
	if a.ID != "A1" || a.SizeX != 10 || a.SizeY != 20 {
		t.Fatalf("unexpected asteroid: %+v", a)
	}
	if len(a.Geysers) != 1 || a.Geysers[0].ID != "g1" {
		t.Fatalf("unexpected geysers: %+v", a.Geysers)
	}
	if len(a.POIs) != 1 || a.POIs[0].ID != "p1" {
		t.Fatalf("unexpected pois: %+v", a.POIs)
	}
	if a.BiomePaths != "bp" {
		t.Fatalf("unexpected biomePaths: %s", a.BiomePaths)
	}
}
