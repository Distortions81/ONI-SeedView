package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

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
