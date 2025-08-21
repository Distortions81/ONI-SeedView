package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

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

// decodeSeed parses the JSON seed data into SeedData.
func decodeSeed(jsonData []byte) (*SeedData, error) {
	var seed SeedData
	if err := json.Unmarshal(jsonData, &seed); err != nil {
		return nil, fmt.Errorf("JSON decode failed: %v", err)
	}
	return &seed, nil
}
