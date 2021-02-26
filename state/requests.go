package state

import "gioman/service"

type Request struct {
	Method service.Method `json:"method"`
	URL    string         `json:"url"`
	Name   string         `json:"name"`
}

type Requests []Request
