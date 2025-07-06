package main

// shortNameTables holds mappings of item IDs to shortened labels.
type shortNameTables struct {
	Geysers map[string]string `json:"geysers"`
	POIs    map[string]string `json:"pois"`
}

var shortNames = shortNameTables{
	Geysers: map[string]string{
		"steam":             "Cool Steam",
		"hot_steam":         "Steam",
		"hot_water":         "Water",
		"slush_water":       "Cool Slush",
		"filthy_water":      "Polluted Water",
		"slush_salt_water":  "Salt Slush",
		"salt_water":        "Salt Water",
		"small_volcano":     "Minor Volcano",
		"big_volcano":       "Volcano",
		"liquid_co2":        "CO2 Vent",
		"hot_co2":           "CO2 Geyser",
		"hot_hydrogen":      "Hydrogen",
		"hot_po2":           "Hot PO2",
		"slimy_po2":         "Infected PO2",
		"chlorine_gas":      "Chlorine",
		"chlorine_gas_cool": "Cool Chlorine",
		"methane":           "Gas",
		"molten_copper":     "Copper",
		"molten_iron":       "Iron",
		"molten_gold":       "Gold",
		"oil_drip":          "Leaky Oil",
		"molten_aluminum":   "Aluminum",
		"molten_cobalt":     "Cobalt",
		"liquid_sulfur":     "Sulfur Vent",
		"molten_tungsten":   "Tungsten",
		"molten_niobium":    "Niobium",
		"OilWell":           "Oil Well",
	},
	POIs: map[string]string{
		"Headquarters":               "Print Pod",
		"WarpConduitSender":          "Tele In",
		"WarpConduitReceiver":        "Tele Out",
		"WarpPortal":                 "Tele Send",
		"WarpReceiver":               "Tele Recv",
		"GeneShuffler":               "Vacillator",
		"MassiveHeatSink":            "Thermo-Null",
		"SapTree":                    "Sap Tree",
		"GravitasPedestal":           "Artifact",
		"PropSurfaceSatellite1":      "Crashed Sat",
		"PropSurfaceSatellite2":      "Wrecked Sat",
		"PropSurfaceSatellite3":      "Crushed Sat",
		"TemporalTearOpener":         "Tear Opener",
		"CryoTank":                   "Cryotank",
		"PropFacilityStatue":         "Vending",
		"GeothermalVentEntity":       "Geo Vent",
		"GeothermalControllerEntity": "Geo Controller",
		"POICeresTechUnlock":         "Ceres Unlock",
	},
}
