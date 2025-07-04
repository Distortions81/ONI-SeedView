//go:build js && wasm

package main

import "syscall/js"

func saveImageData(filename string, data []byte) error {
	uint8Array := js.Global().Get("Uint8Array").New(len(data))
	js.CopyBytesToJS(uint8Array, data)
	blob := js.Global().Get("Blob").New([]any{uint8Array})
	url := js.Global().Get("URL").Call("createObjectURL", blob)
	a := js.Global().Get("document").Call("createElement", "a")
	a.Set("href", url)
	a.Set("download", filename)
	a.Call("click")
	js.Global().Get("URL").Call("revokeObjectURL", url)
	return nil
}
