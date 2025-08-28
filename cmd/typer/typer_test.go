package main

import (
	"testing"
)

func TestInitQuotes(t *testing.T) {
	err := initQuotes("quotes.json")

	if err != nil {
		t.Fatalf("Failed to initialize quotes: %v", err)
	}

	if quotes.Count() == 0 {
		t.Fatal("No quotes loaded")
	}
}
