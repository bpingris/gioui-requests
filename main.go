package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
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
	var (
		fetcher       fetcher
		fetchResponse chan string
	)
	fetch := func(url string) {
		fetchResponse = make(chan string, 1)
		go func() {
			fetchResponse <- fetcher.fetch(state.Request{URL: url})
		}()
	}

	th := material.NewTheme(gofont.Collection())

	// Home state. We haven't wrapped it with anything since it's the only value so far.
	requests := state.Requests{
		state.Request{Method: state.GET, URL: "https://typicode.jsonplaceholder.com/todos/1", Name: "jsonplaceholder"},
		state.Request{Method: state.POST, URL: "https://typicode.jsonplaceholder.com/comments/1", Name: "/comments"},
	}

	response := "Last response N/A"

	home := views.HomeScreen(th, requests, fetch)

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

type fetcher struct {
	cnt uint64
}

func (f *fetcher) fetch(r state.Request) string {
	log.Printf("Fetching %v", r)
	f.cnt++
	// Emulate fetching: 500-1500ms delay.
	time.Sleep(time.Millisecond * time.Duration(500+rand.Intn(1000)))
	resp := fmt.Sprintf("Response #%d", f.cnt)
	log.Printf("Fetched %d bytes", len([]byte(resp)))
	return resp
}
