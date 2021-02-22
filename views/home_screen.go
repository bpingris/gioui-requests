package views

import (
	"sandbox/state"

	"gioui.org/layout"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

type HomeScreenStyle struct {
	url      *widget.Editor
	fetchBtn *widget.Clickable
	home     HomeStyle
	requests state.Requests
	fetch    func(url string)
}

func HomeScreen(th *material.Theme, r state.Requests, fetch func(url string)) HomeScreenStyle {
	url := new(widget.Editor)
	fetchBtn := new(widget.Clickable)
	return HomeScreenStyle{
		url:      url,
		fetchBtn: fetchBtn,
		home:     Home(th, url, fetchBtn),
		requests: r,
		fetch:    fetch,
	}
}

func (h HomeScreenStyle) Layout(gtx layout.Context, fetching bool, response string) layout.Dimensions {
	if h.fetchBtn.Clicked() {
		h.fetch(h.url.Text())
	}
	return h.home.Layout(gtx, h.requests, fetching, response)
}
