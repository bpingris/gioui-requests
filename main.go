package main

import (
	"log"
	"math/rand"
	"os"
	"sandbox/services"
	"sandbox/state"
	"sandbox/views"
	"time"

	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
	"gioui.org/widget/material"
)

func loop(w *app.Window) error {
	var fetchResponse chan string

	fetch := func(url string) {
		fetchResponse = make(chan string, 1)
		services.Fetch(fetchResponse, url)
	}

	th := material.NewTheme(gofont.Collection())

	var requests requestStorage
	requests.add(state.GET, "https://typicode.jsonplaceholder.com/todos/1", "jsonplaceholder")
	requests.add(state.POST, "https://typicode.jsonplaceholder.com/comments/1", "/comments")

	response := "Last response N/A"

	home := views.HomeScreen(th, fetch, (*homeScreenRequestStorageAdaptor)(&requests))

	var ops op.Ops
	for {
		select {
		case e := <-w.Events():
			switch e := e.(type) {
			case system.DestroyEvent:
				return e.Err
			case system.FrameEvent:
				gtx := layout.NewContext(&ops, e)
				fetching := fetchResponse != nil
				home.Layout(gtx, fetching, response)
				e.Frame(gtx.Ops)
			}
		case response = <-fetchResponse:
			fetchResponse = nil
		}
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())

	go func() {
		w := app.NewWindow(app.Size(unit.Dp(600), unit.Dp(400)))
		if err := loop(w); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()
	app.Main()
}

// requestStorage and requestProviderAdaptor exist for demonstration purpose.
type requestStorage struct {
	requests state.Requests
	current  int
}

func (rs *requestStorage) add(m state.Method, url, name string) {
	rs.requests = append(rs.requests, state.Request{
		Method: m,
		URL:    url,
		Name:   name,
	})
}

func (rs *requestStorage) addRequest(r state.Request) {
	rs.requests = append(rs.requests, r)
}

// homeScreenRequestStorageAdaptor and requestStorage exist for demonstration purpose.
type homeScreenRequestStorageAdaptor requestStorage

func (rp *homeScreenRequestStorageAdaptor) All() state.Requests {
	return (*requestStorage)(rp).requests
}

func (rp *homeScreenRequestStorageAdaptor) Save(r state.Request) {
	(*requestStorage)(rp).addRequest(r)
}

func (rp *homeScreenRequestStorageAdaptor) Current() state.Request {
	return rp.requests[rp.current]
}

func (rp *homeScreenRequestStorageAdaptor) SetCurrent(i int) {
	rp.current = i
}
