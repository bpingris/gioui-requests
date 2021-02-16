package routes

import (
	"gioman/ui"
	"gioman/utils"
	"image/color"
	"io/ioutil"
	"log"
	"net/http"

	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	gomaterial "github.com/BenoitPingris/go-material-colors"
)

type Request struct {
	id        string
	name      string
	url       string
	clickable *widget.Clickable
}

type Home struct {
	Th             *material.Theme
	send           *widget.Clickable
	save           *widget.Clickable
	editor         *widget.Editor
	name           *widget.Editor
	response       string
	requests       []Request
	currentRequest Request
}

func NewHome(th *material.Theme) *Home {
	reqs := []Request{{utils.RandString(16), "JSON Placeholder", "https://jsonplaceholder.typicode.com/todos/1", &widget.Clickable{}},
		{utils.RandString(16), "Another API", "https://jsonplaceholder.typicode.com/comments/1", &widget.Clickable{}}}

	h := &Home{
		Th:             th,
		send:           &widget.Clickable{},
		save:           &widget.Clickable{},
		editor:         &widget.Editor{},
		name:           &widget.Editor{},
		requests:       reqs,
		currentRequest: reqs[0],
	}
	h.editor.SetText(reqs[0].url)
	h.name.SetText(reqs[0].name)
	return h
}

func (h *Home) makeRequest() {
	url := h.editor.Text()
	go func() {
		res, err := http.Get(url)
		if err != nil {
			h.response = "An error occured: \n" + err.Error()
		}
		responseData, err := ioutil.ReadAll(res.Body)
		if err != nil {
			log.Fatal(err)
		}
		h.response = string(responseData)
	}()
}

func (h *Home) DrawRequests() layout.FlexChild {
	return layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {

		return layout.Stack{}.Layout(gtx, layout.Expanded(ui.Fill{Color: gomaterial.Grey800}.Layout))
		return layout.UniformInset(unit.Dp(8)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			l := layout.List{Axis: layout.Vertical}
			return l.Layout(gtx, len(h.requests), func(gtx layout.Context, index int) layout.Dimensions {
				request := h.requests[index]
				if request.clickable.Clicked() {
					h.currentRequest = request
					h.editor.SetText(request.url)
				}
				btn := material.Button(h.Th, request.clickable, request.name)
				if h.currentRequest.id != request.id {
					btn.Background.A = 0
					btn.Color = gomaterial.Black
				}
				return layout.Inset{Bottom: unit.Dp(5)}.Layout(gtx, btn.Layout)
			})
		})
	})
}

func (h *Home) Layout(gtx layout.Context) layout.Dimensions {
	if h.save.Clicked() {
		h.requests = append(h.requests, Request{utils.RandString(16), h.name.Text(), h.currentRequest.url, &widget.Clickable{}})
	}
	if h.send.Clicked() {
		h.makeRequest()
	}
	widgets := []layout.Widget{
		func(gtx layout.Context) layout.Dimensions {
			ed := material.Editor(h.Th, h.editor, "URL")
			border := widget.Border{Color: color.NRGBA{A: 200}, CornerRadius: unit.Dp(3), Width: unit.Px(1)}
			return border.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.UniformInset(unit.Dp(8)).Layout(gtx, ed.Layout)
			})
		},
		func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Spacing: layout.SpaceBetween}.Layout(
				gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					border := widget.Border{Color: color.NRGBA{A: 200}, CornerRadius: unit.Dp(3), Width: unit.Px(1)}
					return border.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						return layout.UniformInset(unit.Dp(8)).Layout(gtx, material.Editor(h.Th, h.name, "Name of the request").Layout)
					})
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return layout.Flex{}.Layout(
						gtx,
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							return material.Button(h.Th, h.send, "Send").Layout(gtx)
						}),
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							return layout.Inset{Left: unit.Dp(16)}.Layout(gtx, material.Button(h.Th, h.save, "Save").Layout)
						}),
					)
				}),
			)
		},
		func(gtx layout.Context) layout.Dimensions {
			return material.Body1(h.Th, h.response).Layout(gtx)
		},
	}
	row := layout.Flex{}
	l := layout.List{Axis: layout.Vertical}
	return row.Layout(
		gtx,
		h.DrawRequests(),
		layout.Flexed(3, func(gtx layout.Context) layout.Dimensions {
			return l.Layout(gtx, len(widgets), func(gtx layout.Context, index int) layout.Dimensions {
				return layout.UniformInset(unit.Dp(8)).Layout(gtx, widgets[index])
			})
		}),
	)
}
