package main

import (
	"encoding/json"
	"io"
	"log"
	"os"
	"strings"

	"gioman/service"
	"gioman/state"
)

// [
//     {
//         name: 'qeweqqw',
//         url: 'weqweqweqwe',
//         method: 0
//     }
// ]
// TODO: Need to validate the config, if `method: 921`, it's not valid

type (
	config struct {
		Requests []requestConfig `json: 'requests'`
	}
	requestConfig struct {
		Name   string `json: "name"`
		URL    string `json: "url"`
		Method string `json: "method"`
	}
)

func configFromFilepath(path string) (cfg config, err error) {
	// TODO change the filename path, look for the config folder of the system?
	file, err := os.Open(path)
	if err != nil {
		return
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&cfg)
	return
}

func (cfg *config) save(w io.Writer) error {
	encoder := json.NewEncoder(w)
	encoder.SetIndent("" /*prefix*/, "  " /*indent*/)
	return encoder.Encode(cfg)
}

func (cfg *config) setRequests(requests []state.Request) {
	cfg.Requests = cfg.Requests[:0]
	for _, r := range requests {
		cfg.Requests = append(cfg.Requests, requestConfig{
			Name:   r.Name,
			URL:    r.URL,
			Method: r.Method.String(),
		})
	}
}

func (cfg *config) requests() (requests []state.Request) {
	for _, r := range cfg.Requests {
		method, ok := service.Methods[strings.ToUpper(r.Method)]
		if !ok {
			log.Printf("requests: unknown method %q", method)
		}
		requests = append(requests, state.Request{
			Name:   r.Name,
			URL:    r.URL,
			Method: method,
		})
	}
	return
}
