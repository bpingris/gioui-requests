package service

import (
	"io/ioutil"
	"net/http"
)

type Method int

const (
	GET Method = iota
	POST
	DELETE
	PUT
	PATCH
	OPTIONS
	HEAD
	TRACE
	CONNECT
)

var Methods = map[string]Method{
	"GET":     GET,
	"POST":    POST,
	"DELETE":  DELETE,
	"PUT":     PUT,
	"PATCH":   PATCH,
	"OPTIONS": OPTIONS,
	"HEAD":    HEAD,
	"TRACE":   TRACE,
	"CONNECT": CONNECT,
}

func (m Method) String() string {
	return map[Method]string{
		GET:     "GET",
		POST:    "POST",
		DELETE:  "DELETE",
		PUT:     "PUT",
		PATCH:   "PATCH",
		OPTIONS: "OPTIONS",
		HEAD:    "HEAD",
		TRACE:   "TRACE",
		CONNECT: "CONNECT",
	}[m]
}

func (m Method) Request(url string, headers Headers) (string, error) {
	c := &http.Client{}
	req, err := http.NewRequest(m.String(), url, nil)
	if err != nil {
		return "", err
	}
	for _, h := range headers {
		req.Header.Add(h.Key, h.Value)
	}
	res, err := c.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	bodyBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	return string(bodyBytes), nil
}
