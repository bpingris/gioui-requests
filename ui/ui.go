package ui

import (
	"image/color"

	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

type InputStyle struct {
	Editor material.EditorStyle
}

func Input(th *material.Theme, ed *widget.Editor, hint string) InputStyle {
	return InputStyle{Editor: material.Editor(th, ed, hint)}
}

func (inp InputStyle) Layout(gtx layout.Context) layout.Dimensions {
	// TODO: These values should be moved to InputStyle.
	border := widget.Border{Color: color.NRGBA{A: 200}, CornerRadius: unit.Dp(3), Width: unit.Px(1)}
	return border.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.UniformInset(unit.Dp(8)).Layout(gtx, inp.Editor.Layout)
	})
}
