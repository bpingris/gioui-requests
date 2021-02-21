package views

import (
	"sandbox/state"
	"sandbox/ui"

	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

type Home struct {
	url  widget.Editor
	name widget.Editor
	s    *state.State
}

func NewHome(s *state.State) *Home {
	return &Home{
		s: s,
	}
}

func (h *Home) Layout() layout.Dimensions {
	list := layout.List{Axis: layout.Vertical}
	s := h.s
	r := s.MustGet("requests").(state.Requests)

	return layout.Flex{}.Layout(
		s.Gtx,
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			return list.Layout(
				gtx,
				len(r),
				func(gtx layout.Context, index int) layout.Dimensions {
					return layout.Flex{Alignment: layout.Middle}.Layout(
						gtx,
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							return material.Body1(s.Th, r[index].Method.String()+" "+r[index].Name).Layout(gtx)
						}),
					)
				},
			)
		}),
		layout.Flexed(2, func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return layout.UniformInset(unit.Dp(4)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						return ui.Input(gtx, s.Th, &h.url, "URL")
					})
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return layout.UniformInset(unit.Dp(4)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						return ui.Input(gtx, s.Th, &h.name, "Name")
					})
				}),
			)
		}),
	)
}
