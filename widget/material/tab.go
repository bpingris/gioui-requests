package material

import (
	"image"
	"image/color"

	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

type TabButtonStyle struct {
	clickable *widget.Clickable
	label     material.LabelStyle
	tabcolor  color.NRGBA

	Active bool
}

func TabButton(th *material.Theme, button *widget.Clickable, txt string) TabButtonStyle {
	return TabButtonStyle{clickable: button, label: material.Body1(th, txt), tabcolor: th.ContrastBg}
}

func (t TabButtonStyle) Layout(gtx layout.Context) layout.Dimensions {
	var tabWidth int
	return layout.Stack{Alignment: layout.SW}.Layout(gtx,
		layout.Stacked(func(gtx layout.Context) layout.Dimensions {
			dims := material.Clickable(gtx, t.clickable, func(gtx layout.Context) layout.Dimensions {
				return layout.UniformInset(unit.Sp(10)).Layout(gtx, t.label.Layout)
			})
			tabWidth = dims.Size.X
			return dims
		}),
		layout.Stacked(func(gtx layout.Context) layout.Dimensions {
			if !t.Active {
				return layout.Dimensions{}
			}
			tabHeight := gtx.Px(unit.Dp(4))
			tabRect := image.Rect(0, 0, tabWidth, tabHeight)
			paint.FillShape(gtx.Ops, t.tabcolor, clip.Rect(tabRect).Op())
			return layout.Dimensions{
				Size: image.Point{X: tabWidth, Y: tabHeight},
			}
		}),
	)
}
