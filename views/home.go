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
	loader material.LoaderStyle
	lbl    material.LabelStyle

	name            *widget.Editor
	urlInp, nameInp ui.InputStyle

	save                  *widget.Clickable
	fetchStyle, saveStyle material.ButtonStyle
}

func Home(th *material.Theme, url *widget.Editor, fetch *widget.Clickable) HomeStyle {
	name := new(widget.Editor)
	save := new(widget.Clickable)
	return HomeStyle{
		loader: material.Loader(th),
		lbl:    material.Body1(th, ""),

		name:    name,
		urlInp:  ui.Input(th, url, "URL"),
		nameInp: ui.Input(th, name, "Name"),

		save:       save,
		fetchStyle: material.Button(th, fetch, "Fetch"),
		saveStyle:  material.Button(th, save, "Save"),
	}
}

func (h HomeStyle) Layout(gtx layout.Context, r state.Requests, fetching bool, response string) layout.Dimensions {
	homeLayout := func(gtx layout.Context) layout.Dimensions {
		if fetching {
			gtx = gtx.Disabled()
		}
		return h.layout(gtx, r, response)
	}
	if !fetching {
		return homeLayout(gtx)
	}
	min := gtx.Constraints.Min
	return layout.Stack{}.Layout(gtx,
		layout.Stacked(homeLayout),
		layout.Expanded(func(gtx layout.Context) layout.Dimensions {
			gtx.Constraints.Min = min
			return layout.Center.Layout(gtx, h.loader.Layout)
		}),
	)
}

func (h HomeStyle) layout(gtx layout.Context, r state.Requests, response string) layout.Dimensions {
	methods := func(gtx layout.Context) layout.Dimensions {
		list := layout.List{Axis: layout.Vertical}
		return list.Layout(gtx, len(r), func(gtx layout.Context, index int) layout.Dimensions {
			h.lbl.Text = fmt.Sprintf("%s %s", r[index].Method, r[index].Name)
			return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
				layout.Rigid(h.lbl.Layout),
			)
		})
	}
	inset := func(w layout.Widget) layout.Widget {
		return func(gtx layout.Context) layout.Dimensions {
			return layout.UniformInset(unit.Dp(4)).Layout(gtx, w)
		}
	}
	inputs := func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(inset(h.urlInp.Layout)),
			layout.Rigid(inset(h.nameInp.Layout)),
		)
	}
	buttons := func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
			layout.Rigid(inset(h.fetchStyle.Layout)),
			layout.Rigid(inset(h.saveStyle.Layout)),
		)
	}
	// TODO: May require different style, in which case move it to the constructor.
	resp := h.lbl
	resp.Text = response
	controls := func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(inputs),
			layout.Rigid(buttons),
			layout.Rigid(inset(resp.Layout)),
		)
	}
	return layout.Flex{}.Layout(gtx,
		layout.Flexed(1, methods),
		layout.Flexed(2, controls),
	)
}
