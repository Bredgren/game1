package asset

import (
	"bytes"
	"image"
	"log"
	"path/filepath"

	"github.com/hajimehoshi/ebiten"

	// For decoding png assets
	_ "image/png"
)

var img = map[string]*ebiten.Image{}

const root = "assets"

func init() {
	names, err := AssetDir(filepath.Join(root, "img"))
	if err != nil {
		log.Fatalf("Reading image assets: %s", err)
	}
	for _, imgName := range names {
		r := bytes.NewReader(MustAsset(filepath.Join(root, "img", imgName)))
		i, _, err := image.Decode(r)
		if err != nil {
			log.Fatalf("Loading image %s: %s", imgName, err)
		}

		extension := filepath.Ext(imgName)
		name := imgName[0 : len(imgName)-len(extension)]
		img[name], err = ebiten.NewImageFromImage(i, ebiten.FilterNearest)
		if err != nil {
			log.Fatalf("Converting image %s: %s", imgName, err)
		}
	}
}

// Img retrieves an image by filename (minus extension).
func Img(name string) *ebiten.Image {
	i, ok := img[name]
	if !ok {
		log.Fatalf("No image %s", name)
	}
	return i
}
