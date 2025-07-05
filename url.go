//go:build !js

package main

// coordFromURL returns an empty string on non-web builds.
func coordFromURL() string {
	return ""
}

// asteroidFromURL returns -1 on non-web builds.
func asteroidFromURL() int {
	return -1
}
