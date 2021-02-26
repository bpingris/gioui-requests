package main

import (
	"flag"
	"log"
	"math/rand"
	"os"
	"time"

	"gioman/service"
	"gioman/state"
	"gioman/view"
	mat "gioman/widget/material"

	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
	"gioui.org/widget/material"
)

var configPath = flag.String("config", "config.json", "a configuration file")

func main() {
	flag.Parse()
	rand.Seed(time.Now().UnixNano())

	if *configPath == "" {
		log.Fatalf("-config must be specified")
	}

	cfg, err := configFromFilepath(*configPath)
	if err != nil && !os.IsNotExist(err) {
		log.Fatalf("read config: %v", err)
	}

	requests := make(chan state.Requests)
	go func() {
		save := func(cfg *config) {
			f, err := os.Create(*configPath)
			if err != nil {
				// TODO: Show error notification.
				log.Printf("create config failed: %v", err)
				return
			}
			defer f.Close()
			cfg.save(f)
			log.Printf("saved config to %q", *configPath)
		}
		for r := range requests {
			cfg.setRequests(r)
			save(&cfg)
		}
	}()

	storage := requestStorage{
		requests: cfg.requests(),
		save:     func(r state.Requests) { requests <- r },
	}

	go func() {
		w := app.NewWindow(app.Size(unit.Dp(1000), unit.Dp(600)))
		if err := loop(w, &storage); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()
	app.Main()
}

func loop(w *app.Window, requests *requestStorage) error {
	var (
		fetcher       service.Fetcher
		fetchResponse chan string
	)

	fetch := func(m service.Method, url string) {
		fetchResponse = make(chan string, 1)
		// Ensure closure has its own reference. We need this to guarantee
		// the buffer of size 1 will be used once and only once.
		fetchResponse := fetchResponse
		go func() {
			fetchResponse <- fetcher.Fetch(m, url)
		}()
	}

	th := material.NewTheme(gofont.Collection())

	response := "Last response N/A"

	appbar := mat.Appbar(th)
	home := view.Home(th, fetch, (*homeScreenRequestStorageAdaptor)(requests))

	var ops op.Ops
	for {
		select {
		case e := <-w.Events():
			switch e := e.(type) {
			case system.DestroyEvent:
				return e.Err
			case system.FrameEvent:
				gtx := layout.NewContext(&ops, e)
				fetching := fetchResponse != nil
				appbar.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return home.Layout(gtx, fetching, response)
				})
				e.Frame(gtx.Ops)
			}
		case response = <-fetchResponse:
			fetchResponse = nil
		}
	}
}

// requestStorage and requestProviderAdaptor exist for demonstration purpose.
type requestStorage struct {
	requests state.Requests
	save     func(state.Requests)
}

func (rs *requestStorage) add(m service.Method, url, name string) {
	rs.addRequest(state.Request{
		Method: m,
		URL:    url,
		Name:   name,
	})
}

func (rs *requestStorage) addRequest(r state.Request) {
	rs.requests = append(rs.requests, r)
	rs.save(rs.requests)
}

// homeScreenRequestStorageAdaptor and requestStorage exist for demonstration purpose.
type homeScreenRequestStorageAdaptor requestStorage

func (rp *homeScreenRequestStorageAdaptor) All() state.Requests {
	return (*requestStorage)(rp).requests
}

func (rp *homeScreenRequestStorageAdaptor) Save(r state.Request) {
	(*requestStorage)(rp).addRequest(r)
}

func (rp *homeScreenRequestStorageAdaptor) At(index int) state.Request {
	return (*requestStorage)(rp).requests[index]
}
