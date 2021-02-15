package ui

import (
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/BenoitPingris/giorouter"
	gomaterial "github.com/BenoitPingris/go-material-colors"
	"golang.org/x/exp/shiny/materialdesign/icons"
)

type Appbar struct {
	Router   *giorouter.Router
	Title    string
	BackIcon *widget.Icon
	Back     *widget.Clickable
}

func mustIcon(ic *widget.Icon, err error) *widget.Icon {
	if err != nil {
		panic(err)
	}
	return ic
}

func NewAppbar(router *giorouter.Router, title string) *Appbar {
	backIcon := mustIcon(widget.NewIcon(icons.NavigationChevronLeft))
	backIcon.Color = gomaterial.White
	return &Appbar{
		Router:   router,
		Title:    title,
		BackIcon: backIcon,
		Back:     new(widget.Clickable),
	}
}

func (a *Appbar) Layout(gtx layout.Context) layout.Dimensions {
	b := material.IconButton(a.Router.Th, a.Back, a.BackIcon)
	if a.Back.Clicked() {
		a.Router.Pop()
	}
	return layout.Stack{Alignment: layout.NW}.Layout(gtx,
		layout.Expanded(fill{gomaterial.Grey800}.Layout),
		layout.Stacked(func(gtx layout.Context) layout.Dimensions {
			return Padding(gtx, 16, func(gtx layout.Context) layout.Dimensions {
				insets := layout.Inset{
					Right: unit.Dp(16),
				}
				return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return insets.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							disabled := gtx
							b.Background = gomaterial.Blue600
							if !a.Router.CanPop() {
								disabled = gtx.Disabled()
								b.Background.A = 0
								b.Color.A = 0
							}
							b.Inset = layout.UniformInset(unit.Dp(3))
							return b.Layout(disabled)
						})
					}),
					layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
						text := material.Body1(a.Router.Th, a.Title)
						text.Color = gomaterial.White
						return text.Layout(gtx)
					}),
				)
			})
		}),
	)
}
