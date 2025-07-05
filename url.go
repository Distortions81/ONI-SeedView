//go:build !js

package main

// coordFromURL returns an empty string on non-web builds.
func coordFromURL() string {
	return ""
}

// asteroidFromURL returns "" on non-web builds.
func asteroidFromURL() (string, bool) {
	return "", false
}
