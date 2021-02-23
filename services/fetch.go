package services

import (
	"fmt"
	"log"
	"math/rand"
	"sandbox/state"
	"time"
)

type fetcher struct {
	cnt uint64
}

func (f *fetcher) fetch(r state.Request) string {
	log.Printf("Fetching %v", r)
	f.cnt++
	// Emulate fetching: 500-1500ms delay.
	time.Sleep(time.Millisecond * time.Duration(500+rand.Intn(1000)))
	resp := fmt.Sprintf("Response #%d", f.cnt)
	log.Printf("Fetched %d bytes", len([]byte(resp)))
	return resp
}

func Fetch(response chan string, url string) {
	var fetcher fetcher
	go func() {
		response <- fetcher.fetch(state.Request{URL: url})
	}()
}
