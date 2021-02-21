package main

import (
	"giorequests/routes"
	"giorequests/ui"

	"gioui.org/layout"
)

type App struct {
	appbar *ui.Appbar
	home   *routes.Home
}

func newApp(home *routes.Home) *App {
	return &App{
		ui.NewAppbar(home.Th, "Requestsss"),
		home,
	}
}

func (a App) Layout(gtx layout.Context) layout.Dimensions {
	return layout.Flex{Axis: layout.Vertical}.Layout(
		gtx,
		layout.Rigid(a.appbar.Layout),
		layout.Rigid(a.home.Layout),
	)
}
