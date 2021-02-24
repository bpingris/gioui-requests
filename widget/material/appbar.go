package material

import (
	"image"
	"image/color"

	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget/material"
)

type Appbar struct {
	Th *material.Theme
}

type fill struct {
	color color.NRGBA
}

func (f fill) Layout(gtx layout.Context) layout.Dimensions {
	cs := gtx.Constraints
	d := cs.Min
	paint.FillShape(gtx.Ops, f.color, clip.Rect(image.Rectangle{Max: d}).Op())
	return layout.Dimensions{Size: d, Baseline: d.Y}
}

func (a *Appbar) Layout(gtx layout.Context, w layout.Widget) layout.Dimensions {
	bg := color.NRGBA{R: 97, G: 97, B: 97, A: 255}
	fg := color.NRGBA{R: 255, G: 255, B: 255, A: 255}

	appbar := layout.Rigid(func(gtx layout.Context) layout.Dimensions {
		return layout.Stack{Alignment: layout.NW}.Layout(gtx,
			layout.Expanded(fill{bg}.Layout),
			layout.Stacked(func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{}.Layout(gtx, layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
					lbl := material.Body1(a.Th, "Gioman")
					lbl.Color = fg
					return layout.UniformInset(unit.Dp(15)).Layout(gtx, lbl.Layout)
				}))
			}),
		)
	})

	return layout.Flex{Axis: layout.Vertical}.Layout(gtx, appbar, layout.Flexed(1, w))
}
