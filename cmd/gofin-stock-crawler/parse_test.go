package main

import (
	"log"
	"testing"
)

func TestParse(t *testing.T) {
	symbols, err := parseFile("test.csv")
	if err != nil {
		log.Fatal(err)
	}

	if len(symbols) != 4 {
		t.Errorf("expected len(symbols) to be 4, go %d", len(symbols))
	}

	if symbols[0] != "TSLA" {
		t.Errorf("expected TSLA, got %s", symbols[0])
	}
	if symbols[1] != "ISNS" {
		t.Errorf("expected ISNS, got %s", symbols[1])
	}
	if symbols[2] != "APPL" {
		t.Errorf("expected APPL, got %s", symbols[2])
	}
	if symbols[3] != "MSFT" {
		t.Errorf("expected MSFT, got %s", symbols[3])
	}
}
