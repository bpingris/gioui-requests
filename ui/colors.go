package ui

import (
	"image"
	"image/color"

	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
)

func Rgb(c uint32) color.NRGBA {
	return argb((0xff << 24) | c)
}

func argb(c uint32) color.NRGBA {
	return color.NRGBA{A: uint8(c >> 24), R: uint8(c >> 16), G: uint8(c >> 8), B: uint8(c)}
}

type Fill struct {
	Color color.NRGBA
}

func (f Fill) Layout(gtx layout.Context) layout.Dimensions {
	cs := gtx.Constraints
	d := cs.Min
	paint.FillShape(gtx.Ops, f.Color, clip.Rect(image.Rectangle{Max: d}).Op())
	return layout.Dimensions{Size: d, Baseline: d.Y}
}
