//go:build !js

package main

// coordFromURL returns an empty string on non-web builds.
func coordFromURL() string {
	return ""
}
