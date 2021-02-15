package ui

import (
	"gioui.org/layout"
	"gioui.org/unit"
)

// Padding wrap a Widget with an uniform padding around it
func Padding(gtx layout.Context, padding float32, w layout.Widget) layout.Dimensions {
	return layout.Inset{
		Top:    unit.Dp(padding),
		Right:  unit.Dp(padding),
		Bottom: unit.Dp(padding),
		Left:   unit.Dp(padding),
	}.Layout(gtx, w)
}
