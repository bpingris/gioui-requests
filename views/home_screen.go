package views

import (
	"sandbox/state"

	"gioui.org/layout"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

// requestStorage is HomeScreen -expected interface. It exists for as a
// demonstration and will require an adaptor.
type requestStorage interface {
	All() state.Requests
	Save(state.Request)
}

type HomeScreenStyle struct {
	url, name         *widget.Editor
	fetchBtn, saveBtn *widget.Clickable
	home              HomeStyle
	fetch             func(url string)
	reqProv           requestStorage
}

func HomeScreen(th *material.Theme, fetch func(url string), rp requestStorage) HomeScreenStyle {
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
		reqProv:  rp,
	}
}

func (h HomeScreenStyle) Layout(gtx layout.Context, fetching bool, response string) layout.Dimensions {
	if h.fetchBtn.Clicked() {
		h.fetch(h.url.Text())
	}
	if h.saveBtn.Clicked() {
		h.reqProv.Save(state.Request{
			Method: state.GET, // TODO: Change this.
			URL:    h.url.Text(),
			Name:   h.name.Text(),
		})
	}
	return h.home.Layout(gtx, h.reqProv.All(), fetching, response)
}
