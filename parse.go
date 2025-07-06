package main

import (
	"strconv"
	"strings"
)

// parseBiomePaths converts the raw biome path string into structured paths.
func parseBiomePaths(data string) []BiomePath {
	data = strings.ReplaceAll(strings.ReplaceAll(strings.TrimSpace(data), "\r", ""), "\n", " ")
	// Locate biome names followed by a colon without using regular expressions.
	var locs [][2]int
	for i := 0; i < len(data); i++ {
		if data[i] == ':' {
			start := i - 1
			for start >= 0 {
				c := data[start]
				if !((c >= 'A' && c <= 'Z') || (c >= 'a' && c <= 'z') || (c >= '0' && c <= '9') || c == '_') {
					break
				}
				start--
			}
			start++
			if start < i {
				locs = append(locs, [2]int{start, i + 1})
			}
		}
	}

	var paths []BiomePath
	for i, loc := range locs {
		name := strings.TrimSuffix(data[loc[0]:loc[1]], ":")
		start := loc[1]
		end := len(data)
		if i+1 < len(locs) {
			end = locs[i+1][0]
		}
		segment := strings.TrimSpace(data[start:end])
		if segment == "" {
			continue
		}
		bp := parseBiomeLine(name + ":" + segment)
		if bp.Name != "" {
			paths = append(paths, bp)
		}
	}
	return paths
}

func parseBiomeLine(line string) BiomePath {
	parts := strings.SplitN(strings.TrimSpace(line), ":", 2)
	if len(parts) != 2 {
		return BiomePath{}
	}
	bp := BiomePath{Name: parts[0]}
	segs := strings.Split(parts[1], ";")
	for _, seg := range segs {
		seg = strings.TrimSpace(seg)
		if seg == "" {
			continue
		}
		var poly []Point
		for _, pair := range strings.Fields(seg) {
			xy := strings.Split(pair, ",")
			if len(xy) != 2 {
				continue
			}
			x, err1 := strconv.Atoi(xy[0])
			y, err2 := strconv.Atoi(xy[1])
			if err1 != nil || err2 != nil {
				continue
			}
			poly = append(poly, Point{X: x, Y: y})
		}
		if len(poly) > 0 {
			bp.Polygons = append(bp.Polygons, poly)
		}
	}
	return bp
}
