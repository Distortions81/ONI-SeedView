//go:build js && wasm

package main

import "syscall/js"

// openExternalURL opens the given URL in a new browser tab/window when running
// under WebAssembly. Falls back to setting location.href if window.open is not available.
func openExternalURL(u string) {
    w := js.Global()
    if !w.Truthy() {
        return
    }
    if open := w.Get("open"); open.Truthy() {
        // Try to open in a new tab
        open.Invoke(u, "_blank")
        return
    }
    if loc := w.Get("location"); loc.Truthy() {
        loc.Set("href", u)
    }
}

