package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime/trace"
	"sync"
	"time"
)

type TraceDumper struct {
	fr   *trace.FlightRecorder
	once sync.Once // we'll only do this once (per TraceDumper instance)
}

func (td *TraceDumper) Dump() {
	td.once.Do(func() {
		go func() {
			// spawning this goroutine is on the hot path of the request handler, but writing the trace is not
			// even more efficient, but more verbose: we could set up a buffered channel and have a near permanently parked goroutine for trace dumping, but it's probably not worth the extra complexity
			if err := writeTrace(td.fr); err != nil {
				log.Printf("Failed to write trace: %v", err)
			} else {
				log.Println("Trace written successfully")
			}
		}()
	})
}

func main() {
	cfg := trace.FlightRecorderConfig{
		MinAge:   2 * time.Second,
		MaxBytes: 2 * 1024 * 1024, // 2 MB
	}

	fr := trace.NewFlightRecorder(cfg)
	if err := fr.Start(); err != nil {
		log.Fatalf("Failed to start flight recorder: %v", err)
	}
	defer fr.Stop()

	td := &TraceDumper{fr: fr}

	http.HandleFunc("/work", workHandler(td))
	http.ListenAndServe(":8080", nil)
}

func workHandler(td *TraceDumper) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// handle and time the request
		start := time.Now()
		handleRequest(r.URL.Query().Get("type"))
		duration := time.Since(start)

		// if slow request, write trace (but only once)
		if duration > 200*time.Millisecond {
			log.Printf("Slow request detected! Request took %v, writing trace...", duration)
			td.Dump()
		}
	}
}

func handleRequest(requestType string) {
	if requestType == "heavy" {
		doHeavyWork()
	} else {
		doLightWork()
	}
}

func doLightWork() {
	doWork(10_000)
}

func doHeavyWork() {
	doWork(10_000_000)
}

func doWork(iterations int) {
	for i := 0; i < iterations; i++ {
		_ = fmt.Sprintf("Working... %d\n", i)
	}
}

func writeTrace(fr *trace.FlightRecorder) error {
	if !fr.Enabled() {
		return fmt.Errorf("flight recorder is not enabled")
	}

	file, err := os.Create("trace.out")
	if err != nil {
		return fmt.Errorf("failed to create trace file: %v", err)
	}
	defer file.Close()

	_, err = fr.WriteTo(file)
	return err
}
