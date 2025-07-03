package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/fxamacker/cbor/v2"
)

func fetchSeedCBOR(coordinate string) ([]byte, error) {
	url := "https://ingest.mapsnotincluded.org/coordinate/" + coordinate
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Accept", "application/cbor")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected status %d: %s", resp.StatusCode, body)
	}

	return io.ReadAll(resp.Body)
}

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

func main() {
	coordinate := "SNDST-A-7-0-0-0"
	fmt.Println("Fetching:", coordinate)

	cborData, err := fetchSeedCBOR(coordinate)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	jsonData, err := decodeCBORToJSON(cborData)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	filename := coordinate + ".json"
	if err := saveToFile(filename, jsonData); err != nil {
		fmt.Println("Error writing file:", err)
		return
	}

	fmt.Println("Saved to:", filename)
}
