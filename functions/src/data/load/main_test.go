package main

import (
	"strings"
	"testing"

	airportpkg "github.com/tarmac-project/example-airport-lookup-go/pkg/airport"
)

func TestEscapeSQL(t *testing.T) {
	input := "O'Hare \"line\"\nnext"
	got := escapeSQL(input)

	if !strings.Contains(got, "\\'") {
		t.Fatalf("expected escaped single quote, got %q", got)
	}
	if !strings.Contains(got, "\\\"") {
		t.Fatalf("expected escaped double quote, got %q", got)
	}
	if !strings.Contains(got, "\\n") {
		t.Fatalf("expected escaped newline, got %q", got)
	}
}

func TestBuildAirportUpsertQuery(t *testing.T) {
	ap := airportpkg.Airport{
		LocalCode:    "ORD",
		Name:         "Chicago O'Hare International Airport",
		Type:         "large_airport",
		TypeEmoji:    "✈️",
		Continent:    "NA",
		ISOCountry:   "US",
		ISORegion:    "US-IL",
		Municipality: "Chicago",
		Emoji:        "🇺🇸",
		Status:       "open",
	}

	query := buildAirportUpsertQuery(ap)

	if !strings.Contains(query, "INSERT INTO airports") {
		t.Fatalf("expected insert statement, got %q", query)
	}
	if !strings.Contains(query, "ON DUPLICATE KEY UPDATE") {
		t.Fatalf("expected upsert statement, got %q", query)
	}
	if !strings.Contains(query, "Chicago O\\'Hare International Airport") {
		t.Fatalf("expected escaped airport name, got %q", query)
	}
}

func TestMarshalLoadSummary(t *testing.T) {
	payload := marshalLoadSummary(loadSummary{
		FetchedBytes:   100,
		ParsedAirports: 20,
		SuccessUpsert:  19,
		FailedUpsert:   1,
	})

	got := string(payload)
	for _, expected := range []string{
		`"fetched_bytes":100`,
		`"parsed_airports":20`,
		`"successful_upsert":19`,
		`"failed_upsert":1`,
	} {
		if !strings.Contains(got, expected) {
			t.Fatalf("expected summary to contain %q, got %q", expected, got)
		}
	}
}
