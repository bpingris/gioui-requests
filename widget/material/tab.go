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
	label                  material.LabelStyle
	key                    string
	tabcolor, hoveredcolor color.NRGBA
	group                  *widget.Enum
	clickable              *widget.Clickable
}

func TabButton(th *material.Theme, group *widget.Enum, key, label string) TabButtonStyle {
	c := th.ContrastBg
	c.A = 50
	return TabButtonStyle{group: group, label: material.Body1(th, label), tabcolor: th.ContrastBg, key: key, clickable: &widget.Clickable{}, hoveredcolor: c}
}

func (t TabButtonStyle) layout(gtx layout.Context, checked, hovered bool) layout.Dimensions {
	var tabWidth int
	return layout.Stack{Alignment: layout.SW}.Layout(gtx,
		layout.Expanded(func(gtx layout.Context) layout.Dimensions {
			if !hovered {
				return layout.Dimensions{}
			}
			paint.FillShape(gtx.Ops, t.hoveredcolor, clip.Rect{Max: gtx.Constraints.Min}.Op())
			return layout.Dimensions{Size: gtx.Constraints.Min}
		}),
		layout.Stacked(func(gtx layout.Context) layout.Dimensions {
			dims := layout.UniformInset(unit.Sp(10)).Layout(gtx, t.label.Layout)
			tabWidth = dims.Size.X
			return dims
		}),
		layout.Stacked(func(gtx layout.Context) layout.Dimensions {
			if !checked {
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

func (t TabButtonStyle) Layout(gtx layout.Context) layout.Dimensions {
	dims := t.layout(gtx, t.group.Value == t.key, t.group.Hovered == t.key)
	gtx.Constraints.Min = dims.Size
	return t.group.Layout(gtx, t.key)
	// var tabWidth int
	// return layout.Stack{Alignment: layout.SW}.Layout(gtx,
	// 	layout.Stacked(func(gtx layout.Context) layout.Dimensions {
	// 		dims := material.Clickable(gtx, t.clickable, func(gtx layout.Context) layout.Dimensions {
	// 			return layout.UniformInset(unit.Sp(10)).Layout(gtx, t.label.Layout)
	// 		})
	// 		tabWidth = dims.Size.X
	// 		return dims
	// 	}),
	// 	layout.Stacked(func(gtx layout.Context) layout.Dimensions {
	// 		if !t.Active {
	// 			return layout.Dimensions{}
	// 		}
	// 		tabHeight := gtx.Px(unit.Dp(4))
	// 		tabRect := image.Rect(0, 0, tabWidth, tabHeight)
	// 		paint.FillShape(gtx.Ops, t.tabcolor, clip.Rect(tabRect).Op())
	// 		return layout.Dimensions{
	// 			Size: image.Point{X: tabWidth, Y: tabHeight},
	// 		}
	// 	}),
	// )
}
