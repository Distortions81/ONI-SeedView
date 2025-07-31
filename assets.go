package main

import (
	"embed"
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"io"
	"path"
	"path/filepath"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
)

// embed all image assets so no runtime fetching is needed
//
//go:embed objects/*.png icons/*.png biomes/*.png
var assetFS embed.FS

func openAsset(name string) (io.ReadCloser, error) {
	n := path.Clean(name)
	n = strings.TrimPrefix(n, "../")
	if !strings.HasPrefix(n, "objects/") && !strings.HasPrefix(n, "icons/") && !strings.HasPrefix(n, "biomes/") {
		n = path.Join("objects", n)
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
	defer func() { _ = f.Close() }()
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

func loadImageFile(name string) (*ebiten.Image, error) {
	resolved := resolveAssetName(name)
	f, err := openAsset(resolved)
	if err != nil {
		return nil, err
	}
	defer func() { _ = f.Close() }()
	var src image.Image
	ext := strings.ToLower(filepath.Ext(name))
	switch ext {
	case ".png":
		src, err = png.Decode(f)
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

func loadBiomeTextures() map[string]*ebiten.Image {
	textures := make(map[string]*ebiten.Image)
	for _, name := range biomeOrder {
		img, err := loadImageFile("../biomes/" + name + ".png")
		if err != nil {
			//fmt.Println("load biome", name, ":", err)
			continue
		}
		textures[name] = img
	}
	return textures
}
