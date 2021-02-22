package views

import (
	"fmt"
	"sandbox/state"
	"sandbox/ui"

	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

type HomeStyle struct {
	url  *widget.Editor
	name *widget.Editor

	lbl             material.LabelStyle
	urlInp, nameInp ui.InputStyle
}

func Home(th *material.Theme) HomeStyle {
	url := new(widget.Editor)
	name := new(widget.Editor)
	return HomeStyle{
		url:     url,
		name:    name,
		lbl:     material.Body1(th, ""),
		urlInp:  ui.Input(th, url, "URL"),
		nameInp: ui.Input(th, name, "Name"),
	}
}

func (h HomeStyle) Layout(gtx layout.Context, r state.Requests) layout.Dimensions {
	methods := func(gtx layout.Context) layout.Dimensions {
		list := layout.List{Axis: layout.Vertical}
		return list.Layout(gtx, len(r), func(gtx layout.Context, index int) layout.Dimensions {
			h.lbl.Text = fmt.Sprintf("%s %s", r[index].Method, r[index].Name)
			return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
				layout.Rigid(h.lbl.Layout),
			)
		})
	}
	inputs := func(gtx layout.Context) layout.Dimensions {
		inset := func(w layout.Widget) layout.Widget {
			return func(gtx layout.Context) layout.Dimensions {
				return layout.UniformInset(unit.Dp(4)).Layout(gtx, w)
			}
		}
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(inset(h.urlInp.Layout)),
			layout.Rigid(inset(h.nameInp.Layout)),
		)
	}
	return layout.Flex{}.Layout(gtx,
		layout.Flexed(1, methods),
		layout.Flexed(2, inputs),
	)
}
