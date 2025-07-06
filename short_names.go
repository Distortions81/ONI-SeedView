package main

import (
	"encoding/json"

	_ "embed"
)

//go:embed data/short_names.json
var shortNamesData []byte

// shortNameTables holds mappings of item IDs to shortened labels.
type shortNameTables struct {
	Geysers map[string]string `json:"geysers"`
	POIs    map[string]string `json:"pois"`
}

var shortNames shortNameTables

func init() {
	_ = json.Unmarshal(shortNamesData, &shortNames)
}
