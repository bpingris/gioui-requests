package service

import (
	"fmt"
	"log"
	"math/rand"
	"sandbox/state"
	"time"
)

type Fetcher struct {
	cnt uint64
}

func (f *Fetcher) Fetch(m state.Method, url string) string {
	log.Printf("Fetching %v %v", m, url)
	f.cnt++
	// Emulate fetching: 500-1500ms delay.
	time.Sleep(time.Millisecond * time.Duration(500+rand.Intn(1000)))
	resp := fmt.Sprintf("Response #%d", f.cnt)
	log.Printf("Fetched %d bytes", len([]byte(resp)))
	return resp
}
