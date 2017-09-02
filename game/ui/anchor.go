package ui

import "github.com/Bredgren/geo"

// Anchor holds fields for positioning elements relative to each other.
// Center element:
//  Src: (0.5, 0.5)
//  Dst: (0.5, 0.5)
//  Offset: (0, 0)
// Left align and vertically center with some padding
//  Src: (0, 0.5)
//  Dst: (0, 0.5)
//  Offset: (10, 0)
type Anchor struct {
	// Src is a percentage of the source's bounding rectangle to use as an anchor point.
	// E.g (0, 0) uses is the top left as an anchor, and (1, 1) uses the bottom right.
	Src geo.Vec
	// Dst is a percentage of the bounds that the element is drawn within to use as an
	// anchor point.
	Dst geo.Vec
	// Offset positions the source's Anchor point at this offset relative to the destination
	// anchor point.
	Offset geo.Vec
}

// TopLeft returns the top left position that the source element should be to align
// the anchors with the offset.
func (a Anchor) TopLeft(srcBounds, dstBounds geo.Rect) geo.Vec {
	srcAnchor := geo.VecXY(srcBounds.W*a.Src.X, srcBounds.H*a.Src.Y)

	dstAnchorOffset := geo.VecXY(dstBounds.W*a.Dst.X, dstBounds.H*a.Dst.Y)
	dstAnchor := geo.VecXY(dstBounds.TopLeft()).Plus(dstAnchorOffset)

	return dstAnchor.Minus(srcAnchor).Plus(a.Offset)
}
