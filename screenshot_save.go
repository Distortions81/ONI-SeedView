//go:build !js

package main

import "os"

func saveImageData(filename string, data []byte) error {
	return os.WriteFile(filename, data, 0644)
}
