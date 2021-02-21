package main

import (
	"giorequests/routes"
	"os"

	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/widget/material"
)

func main() {

	th := material.NewTheme(gofont.Collection())

	home := routes.NewHome(th)

	w := app.NewWindow()
	appl := newApp(home)
	go func() {
		var ops op.Ops

		for {
			select {
			case e := <-w.Events():
				switch e := e.(type) {
				case system.DestroyEvent:
					os.Exit(0)
				case system.FrameEvent:
					gtx := layout.NewContext(&ops, e)
					appl.Layout(gtx)
					e.Frame(gtx.Ops)
				}
			}
		}
	}()
	app.Main()
}
