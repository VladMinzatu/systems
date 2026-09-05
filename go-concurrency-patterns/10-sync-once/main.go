package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	loader := newConfigLoader()
	wg := sync.WaitGroup{}

	for i := 1; i <= 5; i++ {
		id := i
		wg.Go(func() {
			cfg := loader.get()
			fmt.Printf("worker %d using config: %+v\n", id, cfg)
		})
	}

	wg.Wait()
}

// configLoader fetches its config from a slow source lazily, on first use, and only
// once, no matter how many goroutines call get() concurrently: sync.Once.Do blocks
// every caller until the first one finishes running load, then lets them all through,
// so there's no need for a separate mutex to guard cfg once it's been set.
type configLoader struct {
	once sync.Once
	cfg  config
}

type config struct {
	Endpoint string
}

func newConfigLoader() *configLoader {
	return &configLoader{}
}

func (l *configLoader) get() config {
	l.once.Do(l.load) // if load panics, Once still counts as "done", so later callers would silently get the zero value instead of retrying
	return l.cfg
}

func (l *configLoader) load() {
	fmt.Println("loading config... (this should print exactly once)")
	time.Sleep(200 * time.Millisecond) // simulate a slow network/disk read
	l.cfg = config{Endpoint: "https://example.com/api"}
}
