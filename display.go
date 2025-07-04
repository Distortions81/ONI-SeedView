package main

import (
	"strings"
)

// simplifyID strips namespace prefixes from an ID string.
func simplifyID(id string) string {
	if idx := strings.LastIndex(id, "/"); idx >= 0 {
		id = id[idx+1:]
	}
	if idx := strings.LastIndex(id, "::"); idx >= 0 {
		id = id[idx+2:]
	}
	return id
}

func iconForGeyser(id string) string {
	id = simplifyID(id)
	switch id {
	case "steam":
		return "geyser_cool_steam_vent.png"
	case "hot_steam":
		return "geyser_steam_vent.png"
	case "hot_water":
		return "geyser_water.png"
	case "slush_water":
		return "geyser_cool_slush_geyser.png"
	case "filthy_water":
		return "geyser_polluted_water_vent.png"
	case "slush_salt_water":
		return "geyser_cool_salt_slush_geyser.png"
	case "salt_water":
		return "geyser_salt_water.png"
	case "small_volcano":
		return "geyser_minor_volcano.png"
	case "big_volcano":
		return "geyser_volcano.png"
	case "liquid_co2":
		return "geyser_carbon_dioxide.png"
	case "hot_co2":
		return "geyser_carbon_dioxide_vent.png"
	case "hot_hydrogen":
		return "geyser_hydrogen_vent.png"
	case "hot_po2":
		return "geyser_hot_polluted_oxygen_vent.png"
	case "slimy_po2":
		return "geyser_infectious_polluted_oxygen_vent.png"
	case "chlorine_gas", "chlorine_gas_cool":
		return "geyser_chlorine_gas_vent.png"
	case "methane":
		return "geyser_natural_gas_geyser.png"
	case "molten_copper":
		return "geyser_copper_volcano.png"
	case "molten_iron":
		return "geyser_iron_volcano.png"
	case "molten_gold":
		return "geyser_gold_volcano.png"
	case "oil_drip":
		return "geyser_leaky_oil_fissure.png"
	case "molten_aluminum":
		return "geyser_aluminum_volcano.png"
	case "molten_cobalt":
		return "geyser_cobalt_volcano.png"
	case "liquid_sulfur":
		return "geyser_liquid_sulfur_geyser.png"
	case "molten_tungsten":
		return "geyser_tungsten_volcano.png"
	case "molten_niobium":
		return "geyser_niobium_volcano.png"
	case "OilWell":
		return "geyser_oil_reservoir.png"
	default:
		return ""
	}
}

func iconForPOI(id string) string {
	id = simplifyID(id)
	switch id {
	case "Headquarters":
		return "building_printing_pod.png"
	case "WarpConduitSender":
		return "building_supply_teleporter_input.png"
	case "WarpConduitReceiver":
		return "building_supply_teleporter_output.png"
	case "WarpPortal":
		return "building_teleporter_transmitter.png"
	case "WarpReceiver":
		return "building_teleporter_receiver.png"
	case "GeneShuffler":
		return "building_neural_vacillator.png"
	case "MassiveHeatSink":
		return "building_anti_entropy_thermo_nullifier.png"
	case "SapTree":
		return "building_sap_tree.png"
	case "GravitasPedestal":
		return "poi_artifact_outline.png"
	case "PropSurfaceSatellite1":
		return "poi_crashed_satellite.png"
	case "PropSurfaceSatellite2":
		return "poi_wrecked_satellite.png"
	case "PropSurfaceSatellite3":
		return "poi_crushed_satellite.png"
	case "TemporalTearOpener":
		return "building_temporal_tear_opener.png"
	case "CryoTank":
		return "building_cryotank.png"
	case "PropFacilityStatue":
		return "poi_prop_facility_statue.png"
	case "GeothermalVentEntity":
		return "poi_geothermal_vent_entity.png"
	case "GeothermalControllerEntity":
		return "poi_geothermal_controller_entity.png"
	case "POICeresTechUnlock":
		return "poi_ceres_tech_unlock.png"
	default:
		return ""
	}
}

func displayBiome(id string) string {
	if v, ok := names.Biomes[id]; ok {
		return v
	}
	return id
}

func displayGeyser(id string) string {
	id = simplifyID(id)
	if v, ok := names.Geysers[id]; ok {
		return v
	}
	return id
}

func displayPOI(id string) string {
	id = simplifyID(id)
	if v, ok := names.POIs[id]; ok {
		return v
	}
	return id
}

// formatLabel splits long names so they wrap every two words.
// It returns the formatted text and the length of the longest line.
func formatLabel(name string) (string, int) {
	words := strings.Fields(name)
	if len(words) <= 2 {
		return name, len(name)
	}
	var lines []string
	for i := 0; i < len(words); i += 2 {
		end := i + 2
		if end > len(words) {
			end = len(words)
		}
		line := strings.Join(words[i:end], " ")
		lines = append(lines, line)
	}
	width := 0
	for _, l := range lines {
		if len(l) > width {
			width = len(l)
		}
	}
	for i, l := range lines {
		if len(l) < width {
			pad := (width - len(l)) / 2
			lines[i] = strings.Repeat(" ", pad) + l
		}
	}
	formatted := strings.Join(lines, "\n")
	return formatted, width
}
