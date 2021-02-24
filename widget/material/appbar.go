package material

import (
	"image/color"

	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget/material"
)

type AppbarStyle struct {
	Background color.NRGBA
	Label      material.LabelStyle
	Inset      layout.Inset
}

func Appbar(th *material.Theme) AppbarStyle {
	lbl := material.Body1(th, "Gioman")
	lbl.Color = th.ContrastFg

	return AppbarStyle{
		Background: th.ContrastBg,
		Label:      lbl,
		Inset:      layout.UniformInset(unit.Dp(15)),
	}
}

func (a AppbarStyle) Layout(gtx layout.Context, w layout.Widget) layout.Dimensions {
	appbar := func(gtx layout.Context) layout.Dimensions {
		min := gtx.Constraints.Min
		return layout.Stack{Alignment: layout.NW}.Layout(gtx,
			layout.Expanded(func(gtx layout.Context) layout.Dimensions {
				paint.FillShape(gtx.Ops, a.Background, clip.Rect{Max: gtx.Constraints.Min}.Op())
				return layout.Dimensions{Size: gtx.Constraints.Min}
			}),
			layout.Stacked(func(gtx layout.Context) layout.Dimensions {
				gtx.Constraints.Min = min
				return a.Inset.Layout(gtx, a.Label.Layout)
			}),
		)
	}
	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(appbar),
		layout.Flexed(1, w),
	)
}
