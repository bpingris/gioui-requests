package views

import (
	"fmt"
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
	Current() state.Request
}

type HomeScreenStyle struct {
	url, name         *widget.Editor
	fetchBtn, saveBtn *widget.Clickable
	itemsBtn          []*widget.Clickable
	home              HomeStyle
	fetch             func(url string)
	reqStor           requestStorage
}

func HomeScreen(th *material.Theme, fetch func(url string), rs requestStorage) HomeScreenStyle {
	url := new(widget.Editor)
	name := new(widget.Editor)
	fetchBtn := new(widget.Clickable)
	saveBtn := new(widget.Clickable)
	var itemsBtn []*widget.Clickable
	for range rs.All() {
		itemsBtn = append(itemsBtn, new(widget.Clickable))
	}
	return HomeScreenStyle{
		url:      url,
		name:     name,
		fetchBtn: fetchBtn,
		saveBtn:  saveBtn,
		itemsBtn: itemsBtn,
		home:     Home(th, url, name, fetchBtn, saveBtn),
		fetch:    fetch,
		reqStor:  rs,
	}
}

func (h HomeScreenStyle) Layout(gtx layout.Context, fetching bool, response string) layout.Dimensions {
	if h.fetchBtn.Clicked() {
		h.fetch(h.url.Text())
	}
	if h.saveBtn.Clicked() {
		h.reqStor.Save(state.Request{
			Method: state.GET, // TODO: Change this.
			URL:    h.url.Text(),
			Name:   h.name.Text(),
		})
	}
	for _, i := range h.itemsBtn {
		if i.Clicked() {
			fmt.Println("clicked")
		}
	}
	return h.home.Layout(gtx, Requests{ReqList: h.reqStor.All(), Current: h.reqStor.Current()}, fetching, response)
}
