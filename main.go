package main

import (
	"gioman/routes"
	"os"

	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/widget/material"
	"github.com/BenoitPingris/giorouter"
)

func main() {

	th := material.NewTheme(gofont.Collection())
	router := giorouter.NewRouter(th)

	home := routes.NewHome(&router)
	config := routes.NewConfig(&router)

	router.SetRoutes(giorouter.Routes{
		"home":   home,
		"config": config,
	}, "home")

	w := app.NewWindow()
	appl := newApp(&router)
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
			case <-router.C:
				w.Invalidate()
			}
		}
	}()
	app.Main()
}
