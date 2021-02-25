package view

import (
	"fmt"
	"gioman/service"
	"gioman/state"
	mat "gioman/widget/material"
	"strings"

	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

// requestStorage is HomeScreen -expected interface. It exists for as a
// demonstration and will require an adaptor.
type requestStorage interface {
	All() state.Requests
	Save(state.Request)
	At(index int) state.Request
}

type homeStyleState struct {
	URL, Name   widget.Editor
	Fetch, Save widget.Clickable
	Items       []*widget.Clickable
	ItemsLayout []material.ButtonStyle // Cached Items buttons.
	btnStyle    material.ButtonStyle   // To create new buttons only.
}

func (w *homeStyleState) saveRequest(rs requestStorage) {
	r := state.Request{
		Method: service.GET, // TODO: Change this.
		URL:    w.URL.Text(),
		Name:   w.Name.Text(),
	}
	rs.Save(r)
	w.addSavedRequestButton(r)
}

func (w *homeStyleState) addSavedRequestButton(r state.Request) {
	clk := new(widget.Clickable)
	w.Items = append(w.Items, clk)
	txt := fmt.Sprintf("%s %s", r.Method, r.Name)
	btn := w.btnStyle
	btn.Button = clk
	btn.Text = txt
	w.ItemsLayout = append(w.ItemsLayout, btn)
}

func (w *homeStyleState) setRequest(r state.Request) {
	w.URL.SetText(r.URL)
	w.Name.SetText(r.Name)
}

type HomeStyle struct {
	widgets *homeStyleState
	home    homeLayoutStyle
	fetch   func(m service.Method, url string)
	reqStor requestStorage
}

func Home(th *material.Theme, fetch func(m service.Method, url string), rs requestStorage) HomeStyle {
	widgets := &homeStyleState{
		URL: widget.Editor{
			SingleLine: true,
			Submit:     true,
		},
		Name: widget.Editor{
			SingleLine: true,
			Submit:     true,
		},
		btnStyle: material.Button(th, nil, ""), // Store as a style only.
	}
	if all := rs.All(); len(all) > 0 {
		widgets.setRequest(all[0])
	}
	for _, r := range rs.All() {
		widgets.addSavedRequestButton(r)
	}
	return HomeStyle{
		widgets: widgets,
		home:    homeLayout(th, widgets),
		fetch:   fetch,
		reqStor: rs,
	}
}

func (h HomeStyle) Layout(gtx layout.Context, fetching bool, response string) layout.Dimensions {
	if hasSubmitEvent(h.widgets.URL.Events()) || h.widgets.Fetch.Clicked() {
		h.fetch(service.GET, h.widgets.URL.Text())
	}
	if hasSubmitEvent(h.widgets.Name.Events()) || h.widgets.Save.Clicked() {
		h.widgets.saveRequest(h.reqStor)
	}
	for i, c := range h.widgets.Items {
		if c.Clicked() {
			h.widgets.setRequest(h.reqStor.At(i))
		}
	}
	return h.home.Layout(gtx, homeLayoutStyleContext{
		fetching: fetching,
		response: response,
		saved:    h.widgets.ItemsLayout,
	})
}

type homeLayoutStyle struct {
	loader material.LoaderStyle
	resp   mat.InputStyle

	url, name mat.InputStyle

	fetchStyle, saveStyle material.ButtonStyle
	list                  *layout.List
}

func homeLayout(th *material.Theme, state *homeStyleState) homeLayoutStyle {
	return homeLayoutStyle{
		loader: material.Loader(th),
		resp:   mat.Input(th, new(widget.Editor), "Response N/A"),

		url:  mat.Input(th, &state.URL, "URL"),
		name: mat.Input(th, &state.Name, "Name"),

		fetchStyle: material.Button(th, &state.Fetch, "Fetch"),
		saveStyle:  material.Button(th, &state.Save, "Save"),
		list:       &layout.List{Axis: layout.Vertical},
	}
}

type homeLayoutStyleContext struct {
	fetching bool
	response string
	saved    []material.ButtonStyle
}

func (h homeLayoutStyle) Layout(gtx layout.Context, ctx homeLayoutStyleContext) layout.Dimensions {
	homeLayout := func(gtx layout.Context) layout.Dimensions {
		if ctx.fetching {
			gtx = gtx.Disabled()
		}
		return h.layout(gtx, ctx)
	}
	if !ctx.fetching {
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

func (h homeLayoutStyle) layout(gtx layout.Context, ctx homeLayoutStyleContext) layout.Dimensions {
	methods := func(gtx layout.Context) layout.Dimensions {
		return h.list.Layout(gtx, len(ctx.saved), func(gtx layout.Context, index int) layout.Dimensions {
			return layout.UniformInset(unit.Dp(4)).Layout(gtx, ctx.saved[index].Layout)
		})
	}
	inset := func(w layout.Widget) layout.Widget {
		return func(gtx layout.Context) layout.Dimensions {
			return layout.UniformInset(unit.Dp(4)).Layout(gtx, w)
		}
	}
	inputs := func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(inset(h.url.Layout)),
			layout.Rigid(inset(h.name.Layout)),
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
		hasURL := len(strings.TrimSpace(h.url.Editor.Text())) > 0
		hasName := len(strings.TrimSpace(h.name.Editor.Text())) > 0
		return layout.Flex{Axis: layout.Horizontal, Spacing: layout.SpaceStart}.Layout(gtx,
			layout.Rigid(enableIf(inset(h.fetchStyle.Layout), hasURL)),
			layout.Rigid(enableIf(inset(h.saveStyle.Layout), hasName)),
		)
	}
	if h.resp.Editor.Text() != ctx.response {
		h.resp.Editor.SetText(ctx.response)
	}

	controls := func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(inputs),
			layout.Rigid(buttons),
		)
	}

	response := func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{}.Layout(gtx,
			layout.Flexed(1, controls),
			layout.Flexed(1, inset(h.resp.Layout)),
		)
	}

	return layout.Flex{}.Layout(gtx,
		layout.Flexed(1, methods),
		layout.Flexed(3, response),
	)
}

func hasSubmitEvent(evts []widget.EditorEvent) bool {
	for _, e := range evts {
		if _, ok := e.(widget.SubmitEvent); ok {
			return true
		}
	}
	return false
}
