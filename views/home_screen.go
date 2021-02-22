package views

import (
	"sandbox/state"

	"gioui.org/layout"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

type HomeScreenStyle struct {
	url, name         *widget.Editor
	fetchBtn, saveBtn *widget.Clickable
	home              HomeStyle
	fetch             func(url string)
	save              func(state.Request)
}

func HomeScreen(th *material.Theme, fetch func(url string), save func(state.Request)) HomeScreenStyle {
	url := new(widget.Editor)
	name := new(widget.Editor)
	fetchBtn := new(widget.Clickable)
	saveBtn := new(widget.Clickable)
	return HomeScreenStyle{
		url:      url,
		name:     name,
		fetchBtn: fetchBtn,
		saveBtn:  saveBtn,
		home:     Home(th, url, name, fetchBtn, saveBtn),
		fetch:    fetch,
		save:     save,
	}
}

func (h HomeScreenStyle) Layout(gtx layout.Context, r state.Requests, fetching bool, response string) layout.Dimensions {
	if h.fetchBtn.Clicked() {
		h.fetch(h.url.Text())
	}
	if h.saveBtn.Clicked() {
		h.save(state.Request{
			Method: state.GET, // TODO: Change this.
			URL:    h.url.Text(),
			Name:   h.name.Text(),
		})
	}
	return h.home.Layout(gtx, r, fetching, response)
}
