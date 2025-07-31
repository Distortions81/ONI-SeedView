package main

import (
	"fmt"
	"io"
	"net/http"

	"github.com/fxamacker/cbor/v2"
)

// fetchSeedCBOR retrieves the seed data in CBOR format for a given coordinate.
func fetchSeedCBOR(coordinate string) ([]byte, error) {
	url := BaseURL + coordinate
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Accept", AcceptCBORHeader)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %v", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected status %d: %s", resp.StatusCode, body)
	}
	return io.ReadAll(resp.Body)
}

// decodeSeed parses the CBOR seed data into SeedData.
func decodeSeed(cborData []byte) (*SeedData, error) {
	var seed SeedData
	if err := cbor.Unmarshal(cborData, &seed); err != nil {
		return nil, fmt.Errorf("CBOR decode failed: %v", err)
	}
	return &seed, nil
}
