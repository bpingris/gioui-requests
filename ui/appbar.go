package ui

import (
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	gomaterial "github.com/BenoitPingris/go-material-colors"
	"golang.org/x/exp/shiny/materialdesign/icons"
)

type Appbar struct {
	Title    string
	BackIcon *widget.Icon
	Back     *widget.Clickable
	Th       *material.Theme
}

func mustIcon(ic *widget.Icon, err error) *widget.Icon {
	if err != nil {
		panic(err)
	}
	return ic
}

func NewAppbar(th *material.Theme, title string) *Appbar {
	backIcon := mustIcon(widget.NewIcon(icons.NavigationChevronLeft))
	backIcon.Color = gomaterial.White
	return &Appbar{
		Title:    title,
		BackIcon: backIcon,
		Back:     new(widget.Clickable),
		Th:       th,
	}
}

func (a *Appbar) Layout(gtx layout.Context) layout.Dimensions {
	min := gtx.Constraints.Min
	return layout.Stack{Alignment: layout.NW}.Layout(gtx,
		layout.Expanded(Fill{gomaterial.Grey800}.Layout),
		layout.Stacked(func(gtx layout.Context) layout.Dimensions {
			gtx.Constraints.Min = min
			text := material.Body1(a.Th, a.Title)
			text.Color = gomaterial.White
			return layout.UniformInset(unit.Dp(16)).Layout(gtx, text.Layout)
		}),
	)
}
