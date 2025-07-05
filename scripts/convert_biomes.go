package main

import (
	"image/jpeg"
	"image/png"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	files, err := filepath.Glob("../biomes-png/*.png")
	if err != nil {
		log.Fatal(err)
	}
	for _, f := range files {
		in, err := os.Open(f)
		if err != nil {
			log.Fatal(err)
		}
		img, err := png.Decode(in)
		in.Close()
		if err != nil {
			log.Fatal(err)
		}
		base := filepath.Base(f)
		name := strings.TrimSuffix(base, ".png") + ".jpg"
		outName := filepath.Join("../biomes/", name)
		out, err := os.Create(outName)
		if err != nil {
			log.Fatal(err)
		}
		if err := jpeg.Encode(out, img, &jpeg.Options{Quality: 70}); err != nil {
			out.Close()
			log.Fatal(err)
		}
		log.Printf("Converted: %v\n", name)
		out.Close()
	}
}
