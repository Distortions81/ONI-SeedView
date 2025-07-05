//go:build !test

package main

import (
	"embed"
	"fmt"
	"image"
	"image/draw"
	"image/jpeg"
	"image/png"
	"io"
	"path"
	"path/filepath"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
)

// embed all PNG assets so no runtime fetching is needed
//
//go:embed assets/*.png icons/*.png
var assetFS embed.FS

func openAsset(name string) (io.ReadCloser, error) {
	n := path.Clean(name)
	n = strings.TrimPrefix(n, "../")
	if !strings.HasPrefix(n, "assets/") && !strings.HasPrefix(n, "icons/") {
		n = path.Join("assets", n)
	}
	return assetFS.Open(n)
}

func toCamel(s string) string {
	parts := strings.Split(s, "_")
	for i, p := range parts {
		if len(p) == 0 {
			continue
		}
		parts[i] = strings.ToUpper(p[:1]) + p[1:]
	}
	return strings.Join(parts, "")
}

func assetExists(name string) bool {
	f, err := openAsset(name)
	if err != nil {
		return false
	}
	f.Close()
	return true
}

func resolveAssetName(name string) string {
	if assetExists(name) {
		return name
	}
	ext := filepath.Ext(name)
	base := strings.TrimSuffix(name, ext)
	camel := toCamel(base) + ext
	if assetExists(camel) {
		return camel
	}
	for _, prefix := range []string{"geyser_", "building_", "poi_"} {
		if strings.HasPrefix(base, prefix) {
			trimmed := base[len(prefix):]
			camel = toCamel(trimmed) + ext
			if assetExists(camel) {
				return camel
			}
		}
	}
	return name
}

var missingImage = ebiten.NewImage(1, 1)

func loadImageFile(name string) (*ebiten.Image, error) {
	resolved := resolveAssetName(name)
	f, err := openAsset(resolved)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var src image.Image
	ext := strings.ToLower(filepath.Ext(name))
	switch ext {
	case ".png":
		src, err = png.Decode(f)
	case ".jpg", ".jpeg":
		src, err = jpeg.Decode(f)
	default:
		return nil, fmt.Errorf("unsupported image format: %s", ext)
	}
	if err != nil {
		return nil, err
	}
	if nrgba, ok := src.(*image.NRGBA); ok {
		for i := 3; i < len(nrgba.Pix); i += 4 {
			if nrgba.Pix[i] < 64 {
				nrgba.Pix[i] = 0
			}
		}
		return ebiten.NewImageFromImage(nrgba), nil
	}
	bounds := src.Bounds()
	dst := image.NewNRGBA(bounds)
	draw.Draw(dst, bounds, src, bounds.Min, draw.Src)
	for i := 3; i < len(dst.Pix); i += 4 {
		if dst.Pix[i] < 64 {
			dst.Pix[i] = 0
		}
	}
	return ebiten.NewImageFromImage(dst), nil
}

func loadImage(cache map[string]*ebiten.Image, name string) (*ebiten.Image, error) {
	if img, ok := cache[name]; ok {
		if img == missingImage {
			return nil, fmt.Errorf("missing")
		}
		return img, nil
	}
	img, err := loadImageFile(name)
	if err != nil {
		cache[name] = missingImage
		return nil, err
	}
	cache[name] = img
	return img, nil
}

func loadBiomeTextures() map[string]*ebiten.Image {
	f, err := openAsset("../icons/textures.png")
	if err != nil {
		fmt.Println("load atlas:", err)
		return nil
	}
	defer f.Close()
	src, err := png.Decode(f)
	if err != nil {
		fmt.Println("decode atlas:", err)
		return nil
	}
	bounds := src.Bounds()
	atlas := image.NewNRGBA(bounds)
	draw.Draw(atlas, bounds, src, bounds.Min, draw.Src)
	w, h := atlas.Bounds().Dx(), atlas.Bounds().Dy()
	cols, rows := 4, 6
	tw, th := w/cols, h/rows
	textures := make(map[string]*ebiten.Image)
	for i, name := range biomeOrder {
		if i >= cols*rows {
			break
		}
		col := i % cols
		row := i / cols
		r := image.Rect(col*tw+1, row*th+1, (col+1)*tw-1, (row+1)*th-1)
		sub := atlas.SubImage(r).(*image.NRGBA)
		textures[name] = ebiten.NewImageFromImage(sub)
	}
	return textures
}
