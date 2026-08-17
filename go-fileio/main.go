package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	tmpDir, err := os.MkdirTemp("", "go-fileio-buffered-*")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	path := filepath.Join(tmpDir, "lines.txt")
	if err := writeBuffered(path, 10_000); err != nil {
		log.Fatal(err)
	}

	n, err := countLines(path)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("wrote, then read and processed %d lines via buffered io\n", n)
}

func writeBuffered(path string, lines int) (err error) {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer func() {
		if cerr := f.Close(); err == nil {
			err = cerr
		}
	}()

	w := bufio.NewWriter(f)
	for i := range lines {
		if _, err := fmt.Fprintf(w, "line %d\n", i); err != nil {
			return err
		}
	}

	return w.Flush()
}

func countLines(path string) (int, error) {
	f, err := os.Open(path)
	if err != nil {
		return 0, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	count := 0
	for scanner.Scan() {
		line := scanner.Text()
		strings.Split(line, " ")
		count++
	}
	if err := scanner.Err(); err != nil {
		return 0, err
	}
	return count, nil
}
