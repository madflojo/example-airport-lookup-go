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
	"html"

	"github.com/tarmac-project/example-airport-lookup-go/pkg/airport/parsers/csv"
	"github.com/tarmac-project/tarmac/pkg/sdk"
)

type Function struct {
	tarmac *sdk.Tarmac
}

func (f *Function) Handler(_ []byte) ([]byte, error) {
	f.tarmac.Logger.Info("Airport raw data download starting")

	// Fetch the airport data
	data, err := f.tarmac.Function.Call("fetch", []byte(""))
	if err != nil {
		f.tarmac.Logger.Error(fmt.Sprintf("Failed to fetch airport data - %s", err))
		return []byte(""), fmt.Errorf("Failed to fetch airport data: %s", err)
	}

	f.tarmac.Logger.Info("Airport raw data download complete, parsing data")

	// Parse the data
	parser, err := csv.New(bytes.NewReader(data))
	if err != nil {
		f.tarmac.Logger.Error(fmt.Sprintf("Failed to create CSV parser - %s", err))
		return []byte(""), fmt.Errorf("Failed to create CSV parser: %s", err)
	}

	airports, err := parser.Parse()
	if err != nil {
		f.tarmac.Logger.Error(fmt.Sprintf("Failed to parse airport data - %s", err))
		return []byte(""), fmt.Errorf("Failed to parse airport data: %s", err)
	}
	f.tarmac.Logger.Info(fmt.Sprintf("Fetched %d airports", len(airports)))

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
			html.EscapeString(airport.LocalCode),
			html.EscapeString(airport.Name),
			html.EscapeString(airport.Type),
			html.EscapeString(airport.TypeEmoji),
			html.EscapeString(airport.Continent),
			html.EscapeString(airport.ISOCountry),
			html.EscapeString(airport.ISORegion),
			html.EscapeString(airport.Municipality),
			html.EscapeString(airport.Emoji),
			html.EscapeString(airport.Status),
			html.EscapeString(airport.Name),
			html.EscapeString(airport.Type),
			html.EscapeString(airport.TypeEmoji),
			html.EscapeString(airport.Continent),
			html.EscapeString(airport.ISOCountry),
			html.EscapeString(airport.ISORegion),
			html.EscapeString(airport.Municipality),
			html.EscapeString(airport.Emoji),
			html.EscapeString(airport.Status),
		)
		f.tarmac.Logger.Trace(fmt.Sprintf("Executing query: %s", query))

		_, err := f.tarmac.SQL.Query(query)
		if err != nil {
			f.tarmac.Logger.Debug(fmt.Sprintf("Failed to execute query - %s", err))
			failure++
			continue
		}
		success++
	}
	f.tarmac.Logger.Info(fmt.Sprintf("Executed %d queries successfully, %d failures", success, failure))

	return []byte(""), nil
}

func main() {
	var err error

	// Initialize Function
	f := &Function{}

	// Initialize the Tarmac SDK
	f.tarmac, err = sdk.New(sdk.Config{
		Namespace: "tarmac",
		Handler:   f.Handler,
	})
	if err != nil {
		return
	}
}
