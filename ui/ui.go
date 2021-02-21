package ui

import (
	"image/color"

	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

func Input(gtx layout.Context, th *material.Theme, ed *widget.Editor, hint string) layout.Dimensions {
	e := material.Editor(th, ed, hint)
	border := widget.Border{Color: color.NRGBA{A: 200}, CornerRadius: unit.Dp(3), Width: unit.Px(1)}
	return border.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.UniformInset(unit.Dp(8)).Layout(gtx, e.Layout)
	})
}
