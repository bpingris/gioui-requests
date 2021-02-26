package view

import (
	"fmt"
	"gioman/service"
	"gioman/state"
	mat "gioman/widget/material"
	"image"
	"strings"

	"gioui.org/layout"
	"gioui.org/op"
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
	URL, Name    widget.Editor
	Fetch, Save  widget.Clickable
	Header, Body widget.Clickable
	Items        []*widget.Clickable
	ItemsLayout  []material.ButtonStyle // Cached Items buttons.
	btnStyle     material.ButtonStyle   // To create new buttons only.
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
	home    *homeLayoutStyle
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

	if h.widgets.Body.Clicked() {
		h.home.bodyStyle.Active = true
		h.home.headerStyle.Active = false
	}
	if h.widgets.Header.Clicked() {
		h.home.headerStyle.Active = true
		h.home.bodyStyle.Active = false
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

	fetchStyle, saveStyle  material.ButtonStyle
	headerStyle, bodyStyle mat.TabButtonStyle
	list                   *layout.List

	minSZ *image.Point
}

func homeLayout(th *material.Theme, state *homeStyleState) *homeLayoutStyle {
	bs := mat.TabButton(th, &state.Body, "Body")
	hs := mat.TabButton(th, &state.Header, "Header")
	bs.Active = true

	return &homeLayoutStyle{
		loader: material.Loader(th),
		resp:   mat.Input(th, new(widget.Editor), "Response N/A"),

		url:  mat.Input(th, &state.URL, "URL"),
		name: mat.Input(th, &state.Name, "Name"),

		fetchStyle:  material.Button(th, &state.Fetch, "Fetch"),
		saveStyle:   material.Button(th, &state.Save, "Save"),
		headerStyle: hs,
		bodyStyle:   bs,
		list:        &layout.List{Axis: layout.Vertical},

		minSZ: new(image.Point),
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
	if h.resp.Editor.Text() != ctx.response {
		h.resp.Editor.SetText(ctx.response)
	}
	return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
		layout.Flexed(2, methods),
		layout.Flexed(3, h.controlsLayout),
		layout.Flexed(3, inset(h.resp.Layout)),
	)
}

func (h homeLayoutStyle) controlsLayout(gtx layout.Context) layout.Dimensions {
	if *h.minSZ == image.ZP {
		*h.minSZ = maxSize(gtx, h.saveStyle.Layout, h.fetchStyle.Layout)
	}
	fetchButton := func(gtx layout.Context) layout.Dimensions {
		hasURL := len(strings.TrimSpace(h.url.Editor.Text())) > 0
		return disableIf(inset(h.fetchStyle.Layout), !hasURL)(gtx)
	}
	saveButton := func(gtx layout.Context) layout.Dimensions {
		hasName := len(strings.TrimSpace(h.name.Editor.Text())) > 0
		return disableIf(inset(h.saveStyle.Layout), !hasName)(gtx)
	}
	editAndButton := func(e mat.InputStyle, b layout.Widget) layout.Widget {
		return func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
				layout.Flexed(1, e.Layout),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					gtx.Constraints.Min = *h.minSZ
					return b(gtx)
				}),
			)
		}
	}
	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(editAndButton(h.url, fetchButton)),
				layout.Rigid(editAndButton(h.name, saveButton)),
			)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
				layout.Rigid(inset(h.bodyStyle.Layout)),
				layout.Rigid(inset(h.headerStyle.Layout)),
			)
		}),
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

func inset(w layout.Widget) layout.Widget {
	return func(gtx layout.Context) layout.Dimensions {
		return layout.UniformInset(unit.Dp(4)).Layout(gtx, w)
	}
}

func disableIf(w layout.Widget, disable bool) layout.Widget {
	return func(gtx layout.Context) layout.Dimensions {
		if disable {
			gtx = gtx.Disabled()
		}
		return w(gtx)
	}
}

func maxSize(gtx layout.Context, widgets ...layout.Widget) (max image.Point) {
	gtx.Constraints.Min = image.ZP
	defer op.Record(gtx.Ops).Stop()
	for _, w := range widgets {
		sz := w(gtx).Size
		if sz.X > max.X {
			max.X = sz.X
		}
		if sz.Y > max.Y {
			max.Y = sz.Y
		}
	}
	return
}
