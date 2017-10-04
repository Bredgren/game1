package sprite

import (
	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/Bredgren/geo"
	"github.com/hajimehoshi/ebiten"
	"github.com/solovev/gopsd"
)

var (
	spriteRe = regexp.MustCompile(`Sprite +(\w+)`)
	frameRe  = regexp.MustCompile(`Frame +(\d+) +([\d.\w]+)`)
	rectRe   = regexp.MustCompile(`Rect +(\w+)`)
	pointRe  = regexp.MustCompile(`Point +(\w+)`)
	imgRe    = regexp.MustCompile(`Img`)
)

// Psd takes data in psd format and extracts Descs from it.
//
// Layers are parsed with the following format.
//  Sprite <Name>
//    Frame <#> <duration>
//      Rect <name>
//      Point <name>
//      Img
//        <If Img is a folder then the layers it contains are joined, their names don't matter>
func Psd(data []byte) ([]Desc, error) {
	doc, err := gopsd.ParseFromBuffer(data)
	if err != nil {
		return nil, err
	}
	var descs []Desc
	rootLayer := doc.GetTreeRepresentation()
	for _, layer := range rootLayer.Children {
		m := spriteRe.FindStringSubmatch(layer.Name)
		if len(m) == 0 {
			continue
		}
		spriteName := m[1]
		d := Desc{
			Name:   spriteName,
			Frames: make([]FrameDesc, len(layer.Children)),
		}
		for i := range d.Frames {
			d.Frames[i].Points = map[string][]geo.Vec{}
			d.Frames[i].Rects = map[string][]geo.Rect{}
		}
		for _, spriteLayer := range layer.Children {
			m := frameRe.FindStringSubmatch(spriteLayer.Name)
			if len(m) == 0 {
				continue
			}
			frameNum, err := strconv.Atoi(m[1])
			if err != nil {
				return nil, fmt.Errorf("sprite '%s': invalid frame number '%s': %v", spriteName, m[1], err)
			}
			if frameNum >= len(d.Frames) {
				return nil, fmt.Errorf("sprite '%s': frame number '%d' out of bounds", spriteName, frameNum)
			}
			if d.Frames[frameNum].Img != nil {
				return nil, fmt.Errorf("sprite '%s': duplicate frame number '%d'", spriteName, frameNum)
			}
			d.Frames[frameNum].Duration, err = time.ParseDuration(m[2])
			if err != nil {
				return nil, fmt.Errorf("sprite %s, frame %d: %v", spriteName, frameNum, err)
			}
			for _, frameLayer := range spriteLayer.Children {
				r := frameLayer.Rectangle
				switch {
				case rectRe.MatchString(frameLayer.Name):
					m := rectRe.FindStringSubmatch(frameLayer.Name)
					name := m[1]
					rect := geo.RectXYWH(float64(r.X), float64(r.Y), float64(r.Width), float64(r.Height))
					d.Frames[frameNum].Rects[name] = append(d.Frames[frameNum].Rects[name], rect)
				case pointRe.MatchString(frameLayer.Name):
					m := pointRe.FindStringSubmatch(frameLayer.Name)
					name := m[1]
					p := geo.VecXY(float64(r.X), float64(r.Y))
					d.Frames[frameNum].Points[name] = append(d.Frames[frameNum].Points[name], p)
				case imgRe.MatchString(frameLayer.Name):
					img, _ := ebiten.NewImage(int(doc.Width), int(doc.Height), ebiten.FilterNearest)
					if frameLayer.IsFolder {
						// Backwards to make sure that layers are drawn in the correct order
						for i := len(frameLayer.Children) - 1; i >= 0; i-- {
							child := frameLayer.Children[i]
							if err := drawImg(img, child); err != nil {
								return nil, fmt.Errorf("draw img layer for sprite '%s', frame '%d': %v", spriteName, frameNum, err)
							}
						}
					} else {
						if err := drawImg(img, frameLayer); err != nil {
							return nil, fmt.Errorf("draw img for sprite '%s', frame '%d': %v", spriteName, frameNum, err)
						}
					}
					d.Frames[frameNum].Img = img
				}
			}
		}
		descs = append(descs, d)
	}
	return descs, nil
}

func drawImg(img *ebiten.Image, layer *gopsd.Layer) error {
	rawImg, err := layer.GetImage()
	if err != nil {
		return fmt.Errorf("get image from layer: %v", err)
	}
	i, err := ebiten.NewImageFromImage(rawImg, ebiten.FilterNearest)
	if err != nil {
		return fmt.Errorf("create ebiten image from layer: %v", err)
	}
	opts := ebiten.DrawImageOptions{}
	opts.GeoM.Translate(float64(layer.Rectangle.X), float64(layer.Rectangle.Y))
	img.DrawImage(i, &opts)
	return nil
}
