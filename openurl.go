//go:build !js

package main

// openExternalURL is a no-op on non-web builds.
func openExternalURL(_ string) {}

