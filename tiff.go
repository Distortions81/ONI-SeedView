package main

import (
	"bytes"
	"encoding/binary"
	"image"
	"image/draw"
)

// encodeTIFF writes the image to a minimal uncompressed TIFF.
func encodeTIFF(img image.Image) ([]byte, error) {
	b := img.Bounds()
	rgba, ok := img.(*image.RGBA)
	if !ok {
		tmp := image.NewRGBA(b)
		draw.Draw(tmp, b, img, b.Min, draw.Src)
		rgba = tmp
	}

	w := uint32(b.Dx())
	h := uint32(b.Dy())
	data := rgba.Pix

	const (
		typeByte  = 1
		typeASCII = 2
		typeShort = 3
		typeLong  = 4
	)

	nEntries := 10
	headerSize := 8
	ifdSize := 2 + nEntries*12 + 4
	bitsOffset := uint32(headerSize + ifdSize)
	extraOffset := bitsOffset + 8
	dataOffset := extraOffset + 2
	byteCounts := uint32(len(data))

	tags := []struct {
		tag uint16
		typ uint16
		cnt uint32
		val uint32
	}{
		{256, typeLong, 1, w},
		{257, typeLong, 1, h},
		{258, typeShort, 4, bitsOffset},
		{259, typeShort, 1, 1},
		{262, typeShort, 1, 2},
		{273, typeLong, 1, dataOffset},
		{277, typeShort, 1, 4},
		{278, typeLong, 1, h},
		{279, typeLong, 1, byteCounts},
		{338, typeShort, 1, extraOffset},
	}

	buf := new(bytes.Buffer)
	buf.Write([]byte{'I', 'I'})
	binary.Write(buf, binary.LittleEndian, uint16(42))
	binary.Write(buf, binary.LittleEndian, uint32(8))
	binary.Write(buf, binary.LittleEndian, uint16(len(tags)))
	for _, t := range tags {
		binary.Write(buf, binary.LittleEndian, t.tag)
		binary.Write(buf, binary.LittleEndian, t.typ)
		binary.Write(buf, binary.LittleEndian, t.cnt)
		if t.typ == typeShort && t.cnt == 1 {
			binary.Write(buf, binary.LittleEndian, uint16(t.val))
			binary.Write(buf, binary.LittleEndian, uint16(0))
		} else {
			binary.Write(buf, binary.LittleEndian, t.val)
		}
	}
	binary.Write(buf, binary.LittleEndian, uint32(0))
	for i := 0; i < 4; i++ {
		binary.Write(buf, binary.LittleEndian, uint16(8))
	}
	binary.Write(buf, binary.LittleEndian, uint16(2))
	buf.Write(data)
	return buf.Bytes(), nil
}
