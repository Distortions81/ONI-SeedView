//go:build !test

package main

import (
	"fmt"
	"image"
	"image/draw"
	"image/jpeg"
	"image/png"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
)

func openFile(name string) (io.ReadCloser, error) {
	if runtime.GOOS == "js" {
		resp, err := http.Get(name)
		if err != nil {
			return nil, err
		}
		if resp.StatusCode != http.StatusOK {
			b, _ := io.ReadAll(resp.Body)
			resp.Body.Close()
			return nil, fmt.Errorf("failed to fetch %s: %s", name, b)
		}
		return resp.Body, nil
	}
	return os.Open(name)
}

func openAsset(name string) (io.ReadCloser, error) {
	if runtime.GOOS == "js" {
		return openFile(path.Join(WebAssetBase, name))
	}
	return openFile(filepath.Join("assets", name))
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
	atlas, err := loadImageFile("../icons/textures.png")
	if err != nil {
		fmt.Println("load atlas:", err)
		return nil
	}
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
		r := image.Rect(col*tw, row*th, (col+1)*tw, (row+1)*th)
		sub := atlas.SubImage(r).(*ebiten.Image)
		textures[name] = sub
	}
	return textures
}
