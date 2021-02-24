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
	Foreground, Background color.NRGBA
	Label                  material.LabelStyle
	Inset                  layout.Inset
}

func Appbar(th *material.Theme) AppbarStyle {
	return AppbarStyle{
		// TODO: Take these from the theme.
		Foreground: color.NRGBA{R: 255, G: 255, B: 255, A: 255},
		Background: color.NRGBA{R: 97, G: 97, B: 97, A: 255},
		Label:      material.Body1(th, "Gioman"),
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
