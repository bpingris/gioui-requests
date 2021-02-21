package main

import (
	"log"
	"os"
	"sandbox/state"
	"sandbox/views"

	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
	"gioui.org/widget/material"
)

func loop(w *app.Window) error {
	th := material.NewTheme(gofont.Collection())

	s := state.NewState()
	s.Th = th
	s.Set("requests", state.Requests{
		state.Request{Method: state.GET, URL: "https://typicode.jsonplaceholder.com/todos/1", Name: "jsonplaceholder"},
		state.Request{Method: state.POST, URL: "https://typicode.jsonplaceholder.com/comments/1", Name: "/comments"}})

	home := views.NewHome(s)

	var ops op.Ops
	for {
		select {
		case e := <-w.Events():
			switch e := e.(type) {
			case system.DestroyEvent:
				return e.Err
			case system.FrameEvent:
				s.Gtx = layout.NewContext(&ops, e)
				home.Layout()
				e.Frame(s.Gtx.Ops)
			}
		}
	}
}

func main() {
	go func() {
		w := app.NewWindow(app.Size(unit.Dp(600), unit.Dp(400)))
		if err := loop(w); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()
	app.Main()
}
