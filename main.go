package main

import (
	"log"
	"math/rand"
	"os"
	"sandbox/service"
	"sandbox/state"
	"sandbox/view"
	mat "sandbox/widget/material"
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
	var (
		fetcher       service.Fetcher
		fetchResponse chan string
	)

	fetch := func(m service.Method, url string) {
		fetchResponse = make(chan string, 1)
		// Ensure closure has its own reference. We need this to guarantee
		// the buffer of size 1 will be used once and only once.
		fetchResponse := fetchResponse
		go func() {
			fetchResponse <- fetcher.Fetch(m, url)
		}()
	}

	th := material.NewTheme(gofont.Collection())

	var requests requestStorage
	requests.add(service.GET, "https://jsonplaceholder.typicode.com/todos/1", "jsonplaceholder")
	requests.add(service.POST, "https://jsonplaceholder.typicode.com/comments/1", "/comments")

	response := "Last response N/A"

	appbar := mat.Appbar{Th: th}
	home := view.Home(th, fetch, (*homeScreenRequestStorageAdaptor)(&requests))

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
				appbar.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return home.Layout(gtx, fetching, response)
				})
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
}

func (rs *requestStorage) add(m service.Method, url, name string) {
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

func (rp *homeScreenRequestStorageAdaptor) At(index int) state.Request {
	return (*requestStorage)(rp).requests[index]
}
