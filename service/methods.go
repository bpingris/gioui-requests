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

func (m Method) Request(url string) (string, error) {
	c := &http.Client{}
	req, err := http.NewRequest(m.String(), url, nil)
	if err != nil {
		return "", err
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
