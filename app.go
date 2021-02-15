package main

import (
	"gioman/ui"

	"gioui.org/layout"
	"github.com/BenoitPingris/giorouter"
)

type App struct {
	appbar *ui.Appbar
	router *giorouter.Router
}

func newApp(router *giorouter.Router) *App {
	return &App{
		ui.NewAppbar(router, "Requestsss"),
		router,
	}
}

func (a App) Layout(gtx layout.Context) layout.Dimensions {
	return layout.Flex{Axis: layout.Vertical}.Layout(
		gtx,
		layout.Rigid(a.appbar.Layout),
		layout.Rigid(a.router.Layout),
	)
}
