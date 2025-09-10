package main

// parseBiomePaths returns the structured biome paths as-is.
func parseBiomePaths(data BiomePathsCompact) []BiomePath {
	return data.Paths
}
