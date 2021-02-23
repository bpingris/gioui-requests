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
	Current() state.Request
	SetCurrent(i int)
}

type HomeScreenWidgets struct {
	url, name         widget.Editor
	fetchBtn, saveBtn widget.Clickable
	itemsBtn          []widget.Clickable
}

type HomeScreenStyle struct {
	widgets *HomeScreenWidgets
	home    HomeStyle
	fetch   func(url string)
	reqStor requestStorage
}

func HomeScreen(th *material.Theme, fetch func(url string), rs requestStorage) HomeScreenStyle {
	var url widget.Editor
	url.SetText(rs.Current().URL)
	var name widget.Editor
	name.SetText(rs.Current().Name)
	var fetchBtn widget.Clickable
	var saveBtn widget.Clickable
	var itemsBtn []widget.Clickable
	for range rs.All() {
		itemsBtn = append(itemsBtn, widget.Clickable{})
	}
	widgets := &HomeScreenWidgets{
		url:      url,
		name:     name,
		fetchBtn: fetchBtn,
		saveBtn:  saveBtn,
		itemsBtn: itemsBtn,
	}
	return HomeScreenStyle{
		widgets: widgets,
		home:    Home(th, widgets),
		fetch:   fetch,
		reqStor: rs,
	}
}

func (h HomeScreenStyle) Layout(gtx layout.Context, fetching bool, response string) layout.Dimensions {
	if h.widgets.fetchBtn.Clicked() {
		h.fetch(h.widgets.url.Text())
	}
	if h.widgets.saveBtn.Clicked() {
		h.widgets.itemsBtn = append(h.widgets.itemsBtn, widget.Clickable{})
		h.reqStor.Save(state.Request{
			Method: state.GET, // TODO: Change this.
			URL:    h.widgets.url.Text(),
			Name:   h.widgets.name.Text(),
		})
	}
	for i, c := range h.widgets.itemsBtn {
		if c.Clicked() {
			h.reqStor.SetCurrent(i)
			h.widgets.url.SetText(h.reqStor.Current().URL)
			h.widgets.name.SetText(h.reqStor.Current().Name)
		}
	}
	return h.home.Layout(gtx, Requests{ReqList: h.reqStor.All(), Current: h.reqStor.Current()}, fetching, response)
}
