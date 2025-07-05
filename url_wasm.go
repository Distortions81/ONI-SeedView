//go:build js && wasm

package main

import (
	"strings"
	"syscall/js"
)

// coordFromURL retrieves the seed coordinate from the URL if present.
// It supports query parameters like ?coord=SEED or ?seed=SEED and
// also checks the hash fragment for the same patterns.
func coordFromURL() string {
	loc := js.Global().Get("location")
	if !loc.Truthy() {
		return ""
	}
	search := loc.Get("search").String()
	search = strings.TrimPrefix(search, "?")
	for _, part := range strings.Split(search, "&") {
		if strings.HasPrefix(part, "coord=") {
			return strings.TrimPrefix(part, "coord=")
		}
		if strings.HasPrefix(part, "seed=") {
			return strings.TrimPrefix(part, "seed=")
		}
	}
	hash := loc.Get("hash").String()
	hash = strings.TrimPrefix(hash, "#")
	if strings.HasPrefix(hash, "coord=") {
		return strings.TrimPrefix(hash, "coord=")
	}
	if strings.HasPrefix(hash, "seed=") {
		return strings.TrimPrefix(hash, "seed=")
	}
	if hash != "" && !strings.Contains(hash, "=") {
		return hash
	}
	return ""
}
