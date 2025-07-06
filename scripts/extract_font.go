package main

import (
	"bytes"
	"encoding/binary"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"golang.org/x/image/font/opentype"
)

func main() {
	outPath := flag.String("o", "noto_first_font.ttf", "output file")
	inPath := flag.String("i", "../NotoSansMono.ttf", "input font file")
	flag.Parse()

	data, err := os.ReadFile(filepath.Clean(*inPath))
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to read %s: %v\n", *inPath, err)
		os.Exit(1)
	}

	start := 0
	end := len(data)
	if len(data) >= 12 && string(data[:4]) == "ttcf" {
		numFonts := binary.BigEndian.Uint32(data[8:12])
		if numFonts == 0 {
			fmt.Fprintln(os.Stderr, "no fonts in collection")
			os.Exit(1)
		}
		first := binary.BigEndian.Uint32(data[12:16])
		start = int(first)
		if numFonts > 1 {
			second := binary.BigEndian.Uint32(data[16:20])
			end = int(second)
		}
	}

	if start < 0 || end > len(data) || start >= end {
		fmt.Fprintln(os.Stderr, "invalid font offsets")
		os.Exit(1)
	}

	fontData := data[start:end]
	font, err := opentype.Parse(fontData)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse font: %v\n", err)
		os.Exit(1)
	}
	buf := &bytes.Buffer{}
	if _, err := font.WriteSourceTo(nil, buf); err != nil {
		fmt.Fprintf(os.Stderr, "failed to copy font data: %v\n", err)
		os.Exit(1)
	}
	if err := os.WriteFile(filepath.Clean(*outPath), buf.Bytes(), 0o644); err != nil {
		fmt.Fprintf(os.Stderr, "failed to write %s: %v\n", *outPath, err)
		os.Exit(1)
	}
}
