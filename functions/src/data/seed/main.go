/*
The purpose of this function is to load a fixed set of major US airports into
the SQL database. It avoids large CSV parsing and gives deterministic startup
data for local/dev runs.
*/
package main

import (
	"fmt"
	"strings"

	sdk "github.com/tarmac-project/sdk"
	"github.com/tarmac-project/sdk/logging"
	sdksql "github.com/tarmac-project/sdk/sql"
)

type airportSeed struct {
	LocalCode    string
	Name         string
	Type         string
	TypeEmoji    string
	Continent    string
	ISOCountry   string
	ISORegion    string
	Municipality string
	Emoji        string
	Status       string
}

var primaryUSAirports = []airportSeed{
	{
		LocalCode:    "ATL",
		Name:         "Hartsfield Jackson Atlanta International Airport",
		Type:         "large_airport",
		TypeEmoji:    "🛫",
		Continent:    "NA",
		ISOCountry:   "US",
		ISORegion:    "US-GA",
		Municipality: "Atlanta",
		Emoji:        "🇺🇸",
		Status:       "open",
	},
	{
		LocalCode:    "LAX",
		Name:         "Los Angeles International Airport",
		Type:         "large_airport",
		TypeEmoji:    "🛫",
		Continent:    "NA",
		ISOCountry:   "US",
		ISORegion:    "US-CA",
		Municipality: "Los Angeles",
		Emoji:        "🇺🇸",
		Status:       "open",
	},
	{
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
	},
	{
		LocalCode:    "DFW",
		Name:         "Dallas Fort Worth International Airport",
		Type:         "large_airport",
		TypeEmoji:    "🛫",
		Continent:    "NA",
		ISOCountry:   "US",
		ISORegion:    "US-TX",
		Municipality: "Dallas-Fort Worth",
		Emoji:        "🇺🇸",
		Status:       "open",
	},
	{
		LocalCode:    "DEN",
		Name:         "Denver International Airport",
		Type:         "large_airport",
		TypeEmoji:    "🛫",
		Continent:    "NA",
		ISOCountry:   "US",
		ISORegion:    "US-CO",
		Municipality: "Denver",
		Emoji:        "🇺🇸",
		Status:       "open",
	},
	{
		LocalCode:    "JFK",
		Name:         "John F. Kennedy International Airport",
		Type:         "large_airport",
		TypeEmoji:    "🛫",
		Continent:    "NA",
		ISOCountry:   "US",
		ISORegion:    "US-NY",
		Municipality: "New York",
		Emoji:        "🇺🇸",
		Status:       "open",
	},
	{
		LocalCode:    "SFO",
		Name:         "San Francisco International Airport",
		Type:         "large_airport",
		TypeEmoji:    "🛫",
		Continent:    "NA",
		ISOCountry:   "US",
		ISORegion:    "US-CA",
		Municipality: "San Francisco",
		Emoji:        "🇺🇸",
		Status:       "open",
	},
	{
		LocalCode:    "SEA",
		Name:         "Seattle Tacoma International Airport",
		Type:         "large_airport",
		TypeEmoji:    "🛫",
		Continent:    "NA",
		ISOCountry:   "US",
		ISORegion:    "US-WA",
		Municipality: "Seattle",
		Emoji:        "🇺🇸",
		Status:       "open",
	},
	{
		LocalCode:    "LAS",
		Name:         "Harry Reid International Airport",
		Type:         "large_airport",
		TypeEmoji:    "🛫",
		Continent:    "NA",
		ISOCountry:   "US",
		ISORegion:    "US-NV",
		Municipality: "Las Vegas",
		Emoji:        "🇺🇸",
		Status:       "open",
	},
	{
		LocalCode:    "MCO",
		Name:         "Orlando International Airport",
		Type:         "large_airport",
		TypeEmoji:    "🛫",
		Continent:    "NA",
		ISOCountry:   "US",
		ISORegion:    "US-FL",
		Municipality: "Orlando",
		Emoji:        "🇺🇸",
		Status:       "open",
	},
	{
		LocalCode:    "CLT",
		Name:         "Charlotte Douglas International Airport",
		Type:         "large_airport",
		TypeEmoji:    "🛫",
		Continent:    "NA",
		ISOCountry:   "US",
		ISORegion:    "US-NC",
		Municipality: "Charlotte",
		Emoji:        "🇺🇸",
		Status:       "open",
	},
	{
		LocalCode:    "MIA",
		Name:         "Miami International Airport",
		Type:         "large_airport",
		TypeEmoji:    "🛫",
		Continent:    "NA",
		ISOCountry:   "US",
		ISORegion:    "US-FL",
		Municipality: "Miami",
		Emoji:        "🇺🇸",
		Status:       "open",
	},
	{
		LocalCode:    "PHX",
		Name:         "Phoenix Sky Harbor International Airport",
		Type:         "large_airport",
		TypeEmoji:    "🛫",
		Continent:    "NA",
		ISOCountry:   "US",
		ISORegion:    "US-AZ",
		Municipality: "Phoenix",
		Emoji:        "🇺🇸",
		Status:       "open",
	},
	{
		LocalCode:    "IAH",
		Name:         "George Bush Intercontinental Airport",
		Type:         "large_airport",
		TypeEmoji:    "🛫",
		Continent:    "NA",
		ISOCountry:   "US",
		ISORegion:    "US-TX",
		Municipality: "Houston",
		Emoji:        "🇺🇸",
		Status:       "open",
	},
	{
		LocalCode:    "BOS",
		Name:         "Logan International Airport",
		Type:         "large_airport",
		TypeEmoji:    "🛫",
		Continent:    "NA",
		ISOCountry:   "US",
		ISORegion:    "US-MA",
		Municipality: "Boston",
		Emoji:        "🇺🇸",
		Status:       "open",
	},
	{
		LocalCode:    "MSP",
		Name:         "Minneapolis Saint Paul International Airport",
		Type:         "large_airport",
		TypeEmoji:    "🛫",
		Continent:    "NA",
		ISOCountry:   "US",
		ISORegion:    "US-MN",
		Municipality: "Minneapolis",
		Emoji:        "🇺🇸",
		Status:       "open",
	},
	{
		LocalCode:    "DTW",
		Name:         "Detroit Metropolitan Wayne County Airport",
		Type:         "large_airport",
		TypeEmoji:    "🛫",
		Continent:    "NA",
		ISOCountry:   "US",
		ISORegion:    "US-MI",
		Municipality: "Detroit",
		Emoji:        "🇺🇸",
		Status:       "open",
	},
	{
		LocalCode:    "PHL",
		Name:         "Philadelphia International Airport",
		Type:         "large_airport",
		TypeEmoji:    "🛫",
		Continent:    "NA",
		ISOCountry:   "US",
		ISORegion:    "US-PA",
		Municipality: "Philadelphia",
		Emoji:        "🇺🇸",
		Status:       "open",
	},
	{
		LocalCode:    "LGA",
		Name:         "LaGuardia Airport",
		Type:         "large_airport",
		TypeEmoji:    "🛫",
		Continent:    "NA",
		ISOCountry:   "US",
		ISORegion:    "US-NY",
		Municipality: "New York",
		Emoji:        "🇺🇸",
		Status:       "open",
	},
	{
		LocalCode:    "BWI",
		Name:         "Baltimore Washington International Airport",
		Type:         "large_airport",
		TypeEmoji:    "🛫",
		Continent:    "NA",
		ISOCountry:   "US",
		ISORegion:    "US-MD",
		Municipality: "Baltimore",
		Emoji:        "🇺🇸",
		Status:       "open",
	},
}

// Function loads a static list of airports into SQL storage.
type Function struct {
	sdk     *sdk.SDK
	logging logging.Client
	sql     sdksql.Client
}

type seedSummary struct {
	SeededAirports int `json:"seeded_airports"`
	SuccessUpsert  int `json:"successful_upsert"`
	FailedUpsert   int `json:"failed_upsert"`
}

func escapeSQL(v string) string {
	return strings.ReplaceAll(v, "'", "''")
}

// Handler upserts the static seed list into the SQL database.
func (f *Function) Handler(_ []byte) ([]byte, error) {
	f.logging.Info("US primary airport seed starting")

	success := 0
	failure := 0
	for _, airport := range primaryUSAirports {
		query := buildSeedUpsertQuery(airport)

		if _, err := f.sql.Exec(query); err != nil {
			f.logging.Debug(fmt.Sprintf("Failed upsert for %s - %s", airport.LocalCode, err))
			failure++
			continue
		}
		success++
	}

	summary := marshalSeedSummary(seedSummary{
		SeededAirports: len(primaryUSAirports),
		SuccessUpsert:  success,
		FailedUpsert:   failure,
	})
	f.logging.Info(fmt.Sprintf("US primary airport seed complete: %s", string(summary)))
	if failure > 0 {
		return summary, fmt.Errorf(
			"airport seed completed with %d failed upserts: %s",
			failure,
			string(summary),
		)
	}

	return summary, nil
}

func buildSeedUpsertQuery(airport airportSeed) string {
	return fmt.Sprintf(`INSERT INTO airports (
      local_code,
      name,
      type,
      type_emoji,
      continent,
      iso_country,
      iso_region,
      municipality,
      emoji,
      status
    ) VALUES (
      '%s',
      '%s',
      '%s',
      '%s',
      '%s',
      '%s',
      '%s',
      '%s',
      '%s',
      '%s')
    ON DUPLICATE KEY UPDATE
      name = '%s',
      type = '%s',
      type_emoji = '%s',
      continent = '%s',
      iso_country = '%s',
      iso_region = '%s',
      municipality = '%s',
      emoji = '%s',
      status = '%s';`,
		escapeSQL(airport.LocalCode),
		escapeSQL(airport.Name),
		escapeSQL(airport.Type),
		escapeSQL(airport.TypeEmoji),
		escapeSQL(airport.Continent),
		escapeSQL(airport.ISOCountry),
		escapeSQL(airport.ISORegion),
		escapeSQL(airport.Municipality),
		escapeSQL(airport.Emoji),
		escapeSQL(airport.Status),
		escapeSQL(airport.Name),
		escapeSQL(airport.Type),
		escapeSQL(airport.TypeEmoji),
		escapeSQL(airport.Continent),
		escapeSQL(airport.ISOCountry),
		escapeSQL(airport.ISORegion),
		escapeSQL(airport.Municipality),
		escapeSQL(airport.Emoji),
		escapeSQL(airport.Status),
	)
}

func marshalSeedSummary(summary seedSummary) []byte {
	return []byte(fmt.Sprintf(
		`{"seeded_airports":%d,"successful_upsert":%d,"failed_upsert":%d}`,
		summary.SeededAirports,
		summary.SuccessUpsert,
		summary.FailedUpsert,
	))
}

// Initialize sets up SDK and SQL clients required by the seed function.
//
//go:wasmexport wapc_init
func Initialize() {
	var err error

	f := &Function{}

	f.sdk, err = sdk.New(sdk.Config{
		Namespace: "tarmac",
		Handler:   f.Handler,
	})
	if err != nil {
		return
	}

	cfg := f.sdk.Config()

	f.logging, err = logging.New(logging.Config{
		SDKConfig: cfg,
	})
	if err != nil {
		return
	}

	f.sql, err = sdksql.New(sdksql.Config{
		SDKConfig: cfg,
	})
	if err != nil {
		return
	}
}
