package main

import (
	"strings"
	"testing"
)

func TestLineCount(t *testing.T) {
	in := strings.NewReader(
		`Line 1
	This is line 2

	Oh, line 4`)
	result, err := countLines(in)
	if err != nil {
		t.Errorf("Unexpected error while counting lines in test: %v", err)
	}
	if result != 4 {
		t.Errorf("Expected to read 4 lines, but reported %d", result)
	}
}
