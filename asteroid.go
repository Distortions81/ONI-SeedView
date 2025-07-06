package main

import "strings"

// asteroidIndexByID searches the asteroid slice for the given ID and returns
// its index, or -1 if not found.
func asteroidIndexByID(asts []Asteroid, id string) int {
	for i, a := range asts {
		if strings.EqualFold(a.ID, id) {
			return i
		}
	}
	return -1
}
