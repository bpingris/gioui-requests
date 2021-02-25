package material

import (
	"image/color"

	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

type InputStyle struct{ material.EditorStyle }

func Input(th *material.Theme, ed *widget.Editor, hint string) InputStyle {
	return InputStyle{EditorStyle: material.Editor(th, ed, hint)}
}

func (inp InputStyle) Layout(gtx layout.Context) layout.Dimensions {
	// TODO: These values should be moved to InputStyle.
	border := widget.Border{Color: color.NRGBA{A: 200}, CornerRadius: unit.Dp(3), Width: unit.Px(1)}
	return border.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.UniformInset(unit.Dp(8)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			defer op.Save(gtx.Ops).Load()
			clip.Rect{Max: gtx.Constraints.Max}.Add(gtx.Ops)
			return inp.EditorStyle.Layout(gtx)
		})
	})
}
