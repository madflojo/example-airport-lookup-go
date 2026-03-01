/*
The purpose of this function is to download the CSV data (by calling another function), parse it, enrich the data, and
load the contents within the SQL database.

This function is called multiple times throughout the application. It's called by an "init" function and also set
as a scheduled task by itself.
*/
package main

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/tarmac-project/example-airport-lookup-go/pkg/airport/parsers/csv"
	sdk "github.com/tarmac-project/sdk"
	functionsdk "github.com/tarmac-project/sdk/function"
	"github.com/tarmac-project/sdk/logging"
	sdksql "github.com/tarmac-project/sdk/sql"
)

type Function struct {
	sdk      *sdk.SDK
	logging  logging.Client
	function functionsdk.Client
	sql      sdksql.Client
}

func escapeSQL(v string) string {
	return strings.ReplaceAll(v, "'", "''")
}

func (f *Function) Handler(_ []byte) ([]byte, error) {
	f.logging.Info("Airport raw data download starting")

	// Fetch the airport data
	data, err := f.function.Call("fetch", []byte(""))
	if err != nil {
		f.logging.Error(fmt.Sprintf("failed to fetch airport data: %v", err))
		return []byte(""), fmt.Errorf("failed to fetch airport data: %w", err)
	}
	f.logging.Info(fmt.Sprintf("Airport raw data download complete - %d bytes", len(data)))

	f.logging.Info("Airport raw data download complete, parsing data")

	// Parse the data
	parser, err := csv.New(bytes.NewReader(data))
	if err != nil {
		f.logging.Error(fmt.Sprintf("failed to create csv parser: %v", err))
		return []byte(""), fmt.Errorf("failed to create csv parser: %w", err)
	}

	airports, err := parser.Parse()
	if err != nil {
		f.logging.Error(fmt.Sprintf("failed to parse airport data: %v", err))
		return []byte(""), fmt.Errorf("failed to parse airport data: %w", err)
	}
	f.logging.Info(fmt.Sprintf("Fetched %d airports", len(airports)))
	if len(airports) == 0 {
		f.logging.Warn("Parsed zero airport records from CSV data")
	}

	// Update the database
	success := 0
	failure := 0
	for _, airport := range airports {
		query := fmt.Sprintf(`INSERT INTO airports (
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
		f.logging.Trace(fmt.Sprintf("Executing query: %s", query))

		_, err := f.sql.Exec(query)
		if err != nil {
			f.logging.Debug(fmt.Sprintf("Failed to execute query - %s", err))
			failure++
			continue
		}
		success++
	}
	f.logging.Info(fmt.Sprintf("Executed %d queries successfully, %d failures", success, failure))
	summary := fmt.Sprintf(
		`{"fetched_bytes":%d,"parsed_airports":%d,"successful_upsert":%d,"failed_upsert":%d}`,
		len(data),
		len(airports),
		success,
		failure,
	)
	f.logging.Info(fmt.Sprintf("Load summary: %s", summary))

	return []byte(summary), nil
}

//go:wasmexport wapc_init
func Initialize() {
	var err error

	// Initialize Function
	f := &Function{}

	// Initialize the Tarmac SDK
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

	f.function, err = functionsdk.New(functionsdk.Config{
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
