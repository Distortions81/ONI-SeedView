package main

// shortNameTables holds mappings of item IDs to shortened labels.
type shortNameTables struct {
	Geysers map[string]string `json:"geysers"`
	POIs    map[string]string `json:"pois"`
}
