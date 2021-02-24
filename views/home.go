package views

import (
	"fmt"
	"sandbox/state"
	sm "sandbox/widget/material"
	"strings"

	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget/material"
)

type HomeStyle struct {
	widgets *HomeScreenWidgets

	th     *material.Theme
	loader material.LoaderStyle
	lbl    material.LabelStyle

	urlInp, nameInp sm.InputStyle

	fetchStyle, saveStyle material.ButtonStyle
}

// TODO change the awesome names
type Requests struct {
	ReqList state.Requests
	Current state.Request
}

func Home(th *material.Theme, widgets *HomeScreenWidgets) HomeStyle {
	return HomeStyle{
		th:     th,
		loader: material.Loader(th),
		lbl:    material.Body1(th, ""),

		widgets: widgets,

		urlInp:  sm.Input(th, &widgets.url, "URL"),
		nameInp: sm.Input(th, &widgets.name, "Name"),

		fetchStyle: material.Button(th, &widgets.fetchBtn, "Fetch"),
		saveStyle:  material.Button(th, &widgets.saveBtn, "Save"),
	}
}

func (h HomeStyle) Layout(gtx layout.Context, r Requests, fetching bool, response string) layout.Dimensions {
	homeLayout := func(gtx layout.Context) layout.Dimensions {
		if fetching {
			gtx = gtx.Disabled()
		}
		return h.layout(gtx, r.ReqList, r.Current, response)
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

func (h HomeStyle) layout(gtx layout.Context, r state.Requests, current state.Request, response string) layout.Dimensions {

	methods := func(gtx layout.Context) layout.Dimensions {
		list := layout.List{Axis: layout.Vertical}
		return list.Layout(gtx, len(r), func(gtx layout.Context, index int) layout.Dimensions {
			label := fmt.Sprintf("%s %s", r[index].Method, r[index].Name)
			return layout.UniformInset(unit.Dp(4)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return material.Button(h.th, &h.widgets.itemsBtn[index], label).Layout(gtx)
			})
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
	enableIf := func(w layout.Widget, enable bool) layout.Widget {
		return func(gtx layout.Context) layout.Dimensions {
			if !enable {
				gtx = gtx.Disabled()
			}
			return w(gtx)
		}
	}
	buttons := func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{Axis: layout.Horizontal, Spacing: layout.SpaceStart}.Layout(gtx,
			layout.Rigid(enableIf(inset(h.fetchStyle.Layout), len(strings.TrimSpace(h.widgets.url.Text())) > 0)),
			layout.Rigid(enableIf(inset(h.saveStyle.Layout), len(strings.TrimSpace(h.widgets.name.Text())) > 0)),
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
