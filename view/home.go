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

type headerInput struct {
	key, value widget.Editor
}

type homeStyleState struct {
	URL, Name    widget.Editor
	Fetch, Save  widget.Clickable
	TabsGroup    widget.Enum
	Items        []*widget.Clickable
	ItemsLayout  []material.ButtonStyle // Cached Items buttons.
	btnStyle     material.ButtonStyle   // To create new buttons only.
	addHeader    widget.Clickable
	headerInputs []headerInput
}

func (w *homeStyleState) saveRequest(rs requestStorage) {
	var headers []state.Header
	for _, h := range w.headerInputs {
		headers = append(headers, state.Header{Value: h.value.Text(), Key: h.key.Text()})
	}
	r := state.Request{
		Method:  service.GET, // TODO: Change this.
		URL:     w.URL.Text(),
		Name:    w.Name.Text(),
		Headers: headers,
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
	for i, r := range r.Headers {
		w.headerInputs[i].key.SetText(r.Key)
		w.headerInputs[i].value.SetText(r.Value)
	}
}

type HomeStyle struct {
	widgets *homeStyleState
	home    homeLayoutStyle
	fetch   func(service.FetchPayload)
	reqStor requestStorage
}

func Home(th *material.Theme, fetch func(service.FetchPayload), rs requestStorage) HomeStyle {
	widgets := &homeStyleState{
		URL: widget.Editor{
			SingleLine: true,
			Submit:     true,
		},
		Name: widget.Editor{
			SingleLine: true,
			Submit:     true,
		},
		btnStyle:     material.Button(th, nil, ""), // Store as a style only.
		headerInputs: make([]headerInput, 1),
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
		var headers service.Headers
		for _, hh := range h.widgets.headerInputs {
			headers = append(headers, service.Header{Key: hh.key.Text(), Value: hh.value.Text()})
		}
		h.fetch(service.FetchPayload{URL: h.widgets.URL.Text(), Method: service.GET, Headers: headers})
	}
	if hasSubmitEvent(h.widgets.Name.Events()) || h.widgets.Save.Clicked() {
		h.widgets.saveRequest(h.reqStor)
	}
	for i, c := range h.widgets.Items {
		if c.Clicked() {
			h.widgets.setRequest(h.reqStor.At(i))
		}
	}

	if h.widgets.addHeader.Clicked() {
	}

	return h.home.Layout(gtx, homeLayoutStyleContext{
		fetching: fetching,
		response: response,
		saved:    h.widgets.ItemsLayout,
	})
}

type homeLayoutStyle struct {
	state  *homeStyleState
	loader material.LoaderStyle
	resp   mat.InputStyle

	url, name mat.InputStyle

	fetchStyle, saveStyle  material.ButtonStyle
	headerStyle, bodyStyle mat.TabButtonStyle
	list                   *layout.List

	minSZ *image.Point

	headerTab headerTabView
	// headers []
}

type inputKeyValue struct {
	header mat.InputStyle
	value  mat.InputStyle
}

type headerTabView struct {
	fields []inputKeyValue
	add    material.ButtonStyle
}

func homeLayout(th *material.Theme, state *homeStyleState) homeLayoutStyle {
	state.TabsGroup.Value = "body"
	bodyFields := make([]inputKeyValue, 1)
	for i := range bodyFields {
		bodyFields[i] = inputKeyValue{
			header: mat.Input(th, &state.headerInputs[i].key, "header"),
			value:  mat.Input(th, &state.headerInputs[i].value, "value"),
		}
	}

	return homeLayoutStyle{
		state:  state,
		loader: material.Loader(th),
		resp:   mat.Input(th, new(widget.Editor), "Response N/A"),

		url:  mat.Input(th, &state.URL, "URL"),
		name: mat.Input(th, &state.Name, "Name"),

		fetchStyle:  material.Button(th, &state.Fetch, "Fetch"),
		saveStyle:   material.Button(th, &state.Save, "Save"),
		headerStyle: mat.TabButton(th, &state.TabsGroup, "header", "Header"),
		bodyStyle:   mat.TabButton(th, &state.TabsGroup, "body", "Body"),
		list:        &layout.List{Axis: layout.Vertical},

		minSZ: new(image.Point),
		headerTab: headerTabView{
			fields: bodyFields,
			add:    material.Button(th, &state.addHeader, "Add"),
		}}
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
func (h homeLayoutStyle) bodyTabLayout(gtx layout.Context) layout.Dimensions {
	return layout.Dimensions{}
}
func (h homeLayoutStyle) headerTabLayout(gtx layout.Context) layout.Dimensions {
	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			l := &layout.List{Axis: layout.Vertical}
			return l.Layout(gtx, len(h.headerTab.fields), func(gtx layout.Context, index int) layout.Dimensions {
				return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
					layout.Flexed(1, inset(h.headerTab.fields[index].header.Layout)),
					layout.Flexed(1, inset(h.headerTab.fields[index].value.Layout)),
				)
			})
		}),
		layout.Rigid(inset(h.headerTab.add.Layout)))
}

func (h homeLayoutStyle) tabContentLayout(gtx layout.Context) layout.Dimensions {
	switch h.state.TabsGroup.Value {
	case "body":
		return h.bodyTabLayout(gtx)
	case "header":
		return h.headerTabLayout(gtx)
	}
	return layout.Dimensions{}
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
		layout.Rigid(h.tabContentLayout),
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
