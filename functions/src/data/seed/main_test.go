package main

import (
	"strings"
	"testing"
)

func TestSeedEscapeSQL(t *testing.T) {
	input := "Chicago O'Hare"
	got := escapeSQL(input)

	if got != "Chicago O''Hare" {
		t.Fatalf("expected doubled single quote escaping, got %q", got)
	}
}

func TestBuildSeedUpsertQuery(t *testing.T) {
	query := buildSeedUpsertQuery(airportSeed{
		LocalCode:    "ORD",
		Name:         "Chicago O'Hare International Airport",
		Type:         "large_airport",
		TypeEmoji:    "🛫",
		Continent:    "NA",
		ISOCountry:   "US",
		ISORegion:    "US-IL",
		Municipality: "Chicago",
		Emoji:        "🇺🇸",
		Status:       "open",
	})

	if !strings.Contains(query, "INSERT INTO airports") {
		t.Fatalf("expected insert statement, got %q", query)
	}
	if !strings.Contains(query, "ON DUPLICATE KEY UPDATE") {
		t.Fatalf("expected upsert statement, got %q", query)
	}
	if !strings.Contains(query, "Chicago O''Hare International Airport") {
		t.Fatalf("expected escaped airport name, got %q", query)
	}
}

func TestMarshalSeedSummary(t *testing.T) {
	payload := marshalSeedSummary(seedSummary{
		SeededAirports: 20,
		SuccessUpsert:  20,
		FailedUpsert:   0,
	})

	got := string(payload)
	for _, expected := range []string{
		`"seeded_airports":20`,
		`"successful_upsert":20`,
		`"failed_upsert":0`,
	} {
		if !strings.Contains(got, expected) {
			t.Fatalf("expected summary to contain %q, got %q", expected, got)
		}
	}
}
