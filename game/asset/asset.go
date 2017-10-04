package asset

import (
	"path/filepath"
	// For decoding png assets
	// _ "image/png"
)

// var img = map[string]*ebiten.Image{}

// var sheetDesc = map[string]sprite.SheetDesc{}

const root = "assets"

func init() {
	// initImg()
	// initSheetDesc()
}

// Psd returns the contents of the given PSD file.
func Psd(name string) []byte {
	return MustAsset(filepath.Join(root, "psd", name+".psd"))
}

// func initImg() {
// 	names, err := AssetDir(filepath.Join(root, "img"))
// 	if err != nil {
// 		log.Fatalf("Reading image assets: %s", err)
// 	}
// 	for _, imgName := range names {
// 		r := bytes.NewReader(MustAsset(filepath.Join(root, "img", imgName)))
// 		i, _, err := image.Decode(r)
// 		if err != nil {
// 			log.Fatalf("Loading image %s: %s", imgName, err)
// 		}
//
// 		extension := filepath.Ext(imgName)
// 		name := imgName[0 : len(imgName)-len(extension)]
// 		img[name], err = ebiten.NewImageFromImage(i, ebiten.FilterNearest)
// 		if err != nil {
// 			log.Fatalf("Converting image %s: %s", imgName, err)
// 		}
// 	}
// }
//
// // Img retrieves an image by filename (minus extension).
// func Img(name string) *ebiten.Image {
// 	i, ok := img[name]
// 	if !ok {
// 		log.Fatalf("No image %s", name)
// 	}
// 	return i
// }

// func initSheetDesc() {
// 	names, err := AssetDir(filepath.Join(root, "sheetDesc"))
// 	if err != nil {
// 		log.Fatalf("Reading sheet desc assets: %s", err)
// 	}
// 	for _, descName := range names {
// 		r := MustAsset(filepath.Join(root, "sheetDesc", descName))
// 		var s sprite.SheetDesc
// 		err := json.Unmarshal(r, &s)
// 		if err != nil {
// 			log.Fatalf("Loading sheet %s: %s", descName, err)
// 		}
//
// 		extension := filepath.Ext(descName)
// 		name := descName[0 : len(descName)-len(extension)]
// 		sheetDesc[name] = s
// 	}
// }

// // SheetDesc retrieves an a sprite sheet description by filename (minus extension).
// func SheetDesc(name string) sprite.SheetDesc {
// 	s, ok := sheetDesc[name]
// 	if !ok {
// 		log.Fatalf("No sheetDesc %s", name)
// 	}
// 	return s
// }
