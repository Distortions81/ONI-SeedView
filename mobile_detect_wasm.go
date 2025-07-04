//go:build js && wasm

package main

import (
	"strings"
	"syscall/js"
)

// isMobile attempts to detect mobile browsers when running under WebAssembly.
func isMobile() bool {
	nav := js.Global().Get("navigator")
	if !nav.Truthy() {
		return false
	}
	ua := strings.ToLower(nav.Get("userAgent").String())
	return strings.Contains(ua, "mobile") ||
		strings.Contains(ua, "android") ||
		strings.Contains(ua, "iphone") ||
		strings.Contains(ua, "ipad")
}
