//go:build !js

package main

import "runtime"

// isMobile reports whether the program is running on a native mobile OS.
// When compiled for WebAssembly a separate implementation is used.
func isMobile() bool {
	return runtime.GOOS == "android" || runtime.GOOS == "ios"
}
