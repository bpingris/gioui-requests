package state

import "gioman/service"

type Request struct {
	Method service.Method
	URL    string
	Name   string
}

type Requests []Request
