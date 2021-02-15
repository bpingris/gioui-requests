package routes

import (
	"gioman/ui"

	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/BenoitPingris/giorouter"
)

type Config struct {
	Router *giorouter.Router
	Appbar *ui.Appbar
	click  *widget.Clickable
}

func NewConfig(router *giorouter.Router) *Config {
	return &Config{
		Router: router,
		Appbar: ui.NewAppbar(router, "Gioman"),
		click:  &widget.Clickable{},
	}
}

func (c *Config) Layout(gtx layout.Context) layout.Dimensions {
	if c.click.Clicked() {
		c.Router.Pop()
	}
	widgets := []layout.Widget{
		func(gtx layout.Context) layout.Dimensions {
			return material.H1(c.Router.Th, "Title").Layout(gtx)
		},
		func(gtx layout.Context) layout.Dimensions {
			return material.H2(c.Router.Th, "Another title").Layout(gtx)
		},
		func(gtx layout.Context) layout.Dimensions {
			return material.Button(c.Router.Th, c.click, "Go back").Layout(gtx)
		},
	}
	l := layout.List{Axis: layout.Vertical}
	return l.Layout(gtx, len(widgets), func(gtx layout.Context, index int) layout.Dimensions {
		return layout.UniformInset(unit.Dp(8)).Layout(gtx, widgets[index])
	})
}
