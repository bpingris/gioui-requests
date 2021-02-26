package state

import "gioman/service"

type Header struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type Request struct {
	Method  service.Method `json:"method"`
	URL     string         `json:"url"`
	Name    string         `json:"name"`
	Headers []Header       `json:"headers"`
}

type Requests []Request
