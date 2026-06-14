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

	http.HandleFunc("/work", workHandler(fr))
	http.ListenAndServe(":8080", nil)
}

func workHandler(fr *trace.FlightRecorder) http.HandlerFunc {
	var traceWritten sync.Once // notice the Once, to ensure we only write the trace once - and notice the use of closure to capture the flight recorder instance for use in each request handler
	return func(w http.ResponseWriter, r *http.Request) {
		// handle and time the request
		start := time.Now()
		handleRequest(r.URL.Query().Get("type"))
		duration := time.Since(start)

		// if slow request, write trace (but only once)
		if duration > 200*time.Millisecond {
			traceWritten.Do(func() {
				log.Printf("Slow request detected! Request took %v, writing trace...", duration)
				if err := writeTrace(fr); err != nil {
					log.Printf("Failed to write trace: %v", err)
				} else {
					log.Println("Trace written successfully")
				}
			})
		}
	}
}

func handleRequest(requestType string) {
	if requestType == "heavy" {
		heavyWork()
	} else {
		lightWork()
	}
}

func lightWork() {
	doWork(10_000)
}

func heavyWork() {
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
