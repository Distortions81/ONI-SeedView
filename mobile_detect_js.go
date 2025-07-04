//go:build js && wasm

package main

import (
	"strings"
	"syscall/js"
)

func isMobile() bool {
	nav := js.Global().Get("navigator")
	if !nav.Truthy() {
		return false
	}
	ua := strings.ToLower(nav.Get("userAgent").String())
	if strings.Contains(ua, "android") || strings.Contains(ua, "iphone") || strings.Contains(ua, "ipad") {
		return true
	}
	return strings.Contains(ua, "mobile")
}
