package main

// Data structures used to decode the geyser information.
type Geyser struct {
	ID             string
	X              int
	Y              int
	ActiveCycles   float64
	AvgEmitRate    float64
	DormancyCycles float64
	EmitRate       float64
	EruptionTime   float64
	IdleTime       float64
}

type PointOfInterest struct {
	ID string
	X  int
	Y  int
}

type Point struct {
	X int
	Y int
}

type BiomePath struct {
	Name     string
	Polygons [][]Point
}

type BiomePathsCompact struct {
	Paths []BiomePath
}

type Asteroid struct {
	ID         string
	SizeX      int
	SizeY      int
	Geysers    []Geyser
	POIs       []PointOfInterest
	BiomePaths BiomePathsCompact
}

type SeedData struct {
	Asteroids []Asteroid
}

type nameTables struct {
	Biomes  map[string]string
	Geysers map[string]string
	POIs    map[string]string
}

var names = nameTables{
	Biomes: map[string]string{
		"Sandstone":           "Sandstone",
		"Barren":              "Barren",
		"Space":               "Space",
		"FrozenWastes":        "Tundra",
		"CrystalCaverns":      "Crystal Caverns",
		"BoggyMarsh":          "Marsh",
		"ToxicJungle":         "Toxic Jungle",
		"Ocean":               "Ocean",
		"Rust":                "Rust",
		"Forest":              "Forest",
		"Radioactive":         "Radioactive",
		"Swamp":               "Swamp",
		"Wasteland":           "Wasteland",
		"Metallic":            "Metallic",
		"Moo":                 "Moo",
		"IceCaves":            "Ice Caves",
		"CarrotQuarry":        "Carrot Quarry",
		"SugarWoods":          "Sugar Woods",
		"PrehistoricGarden":   "Prehistoric Garden",
		"PrehistoricRaptor":   "Prehistoric Raptor",
		"PrehistoricWetlands": "Prehistoric Wetlands",
		"OilField":            "Oil Field",
		"MagmaCore":           "Magma",
	},
	Geysers: map[string]string{
		"steam":             "Cool Steam Vent",
		"hot_steam":         "Steam Vent",
		"hot_water":         "Water Geyser",
		"slush_water":       "Cool Slush Geyser",
		"filthy_water":      "Polluted Water Vent",
		"slush_salt_water":  "Cool Salt Slush Geyser",
		"salt_water":        "Salt Water Geyser",
		"small_volcano":     "Minor Volcano",
		"big_volcano":       "Volcano",
		"liquid_co2":        "Carbon Dioxide Vent",
		"hot_co2":           "Carbon Dioxide Geyser",
		"hot_hydrogen":      "Hydrogen Vent",
		"hot_po2":           "Hot Polluted Oxygen Vent",
		"slimy_po2":         "Infectious Polluted Oxygen Vent",
		"chlorine_gas":      "Chlorine Gas Vent",
		"chlorine_gas_cool": "Cool Chlorine Vent",
		"methane":           "Natural Gas Geyser",
		"molten_copper":     "Copper Volcano",
		"molten_iron":       "Iron Volcano",
		"molten_gold":       "Gold Volcano",
		"oil_drip":          "Leaky Oil Fissure",
		"molten_aluminum":   "Aluminum Volcano",
		"molten_cobalt":     "Cobalt Volcano",
		"liquid_sulfur":     "Liquid Sulfur Vent",
		"molten_tungsten":   "Tungsten Volcano",
		"molten_niobium":    "Niobium Volcano",
		"OilWell":           "Oil Reservoir",
	},
	POIs: map[string]string{
		"Headquarters":               "Printing Pod",
		"WarpConduitSender":          "Supply Teleporter Input",
		"WarpConduitReceiver":        "Supply Teleporter Output",
		"WarpPortal":                 "Teleporter Transmitter",
		"WarpReceiver":               "Teleporter Receiver",
		"GeneShuffler":               "Neural Vacillator",
		"MassiveHeatSink":            "Anti Entropy Thermo-Nullifier",
		"SapTree":                    "Juicy Sap Tree",
		"GravitasPedestal":           "Artifact Pedestal",
		"PropSurfaceSatellite1":      "Crashed Satellite",
		"PropSurfaceSatellite2":      "Wrecked Satellite",
		"PropSurfaceSatellite3":      "Crushed Satellite",
		"TemporalTearOpener":         "Temporal Tear Opener",
		"CryoTank":                   "Cryotank",
		"PropFacilityStatue":         "Vending Machine",
		"GeothermalVentEntity":       "Geothermal Vent",
		"GeothermalControllerEntity": "Geothermal Controller",
		"POICeresTechUnlock":         "Ceres Tech Unlock",
	},
}
