package csvparser

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"github.com/enescakir/emoji"
	"github.com/tarmac-project/example-airport-lookup-go/pkg/data/airport"
)

type Parser struct{}

var ErrIsHeader = fmt.Errorf("record is a header")
var ErrNotEnoughFields = fmt.Errorf("not enough fields")

// ParseAirport parses a CSV file of Airports returning a slice of Airports.
func ParseAirport(data []byte) ([]airport.Airport, error) {
	// Create a new CSV Reader
	reader := csv.NewReader(bytes.NewReader(data))

	// Read the CSV file
	rec, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("unable to read csv data - %w", err)
	}

	// Create a slice of Airports
	airports := make([]airport.Airport, 0, len(rec))

	// Iterate over the CSV records
	for _, r := range rec {
		// Convert the record to an Airport
		a, err := RecordToAirport(r)
		if err != nil {
			// Skip headers and records with not enough fields
			if err == ErrIsHeader || err == ErrNotEnoughFields {
				continue
			}
			return nil, fmt.Errorf("unable to convert record to airport - %w", err)
		}
		// Append the Airport to the slice
		airports = append(airports, a)
	}

	// Return the slice of Airports
	return airports, nil
}

// RecordToAirport converts a CSV record to an Airport.
func RecordToAirport(rec []string) (airport.Airport, error) {
	// Check minimum length
	if len(rec) < 15 {
		return airport.Airport{}, ErrNotEnoughFields
	}

	// Check if the record is a header
	if rec[0] == "id" {
		return airport.Airport{}, ErrIsHeader
	}

	// Create a new Airport and set basic fields
	a := airport.Airport{
		Continent:    rec[7],
		ISOCountry:   rec[8],
		ISORegion:    rec[9],
		LocalCode:    rec[14],
		Municipality: rec[10],
		Name:         rec[3],
		Status:       "open",
	}

	// Set the Airport type and emoji
	switch rec[2] {
	case "heliport":
		a.Type = "heliport"
		a.TypeEmoji = emoji.Helicopter.String()
	case "small_airport":
		a.Type = "small_airport"
		a.TypeEmoji = emoji.SmallAirplane.String()
	case "medium_airport", "large_airport":
		a.Type = rec[2]
		a.TypeEmoji = emoji.Airplane.String()
	case "seaplane_base":
		a.Type = "seaplane_base"
		a.TypeEmoji = emoji.Anchor.String()
	case "balloonport":
		a.Type = "balloonport"
		a.TypeEmoji = emoji.Balloon.String()
	case "closed":
		a.Type = "unknown"
		a.TypeEmoji = emoji.Construction.String()
		a.Status = "closed"
	default:
		a.Type = rec[2]
		a.TypeEmoji = emoji.QuestionMark.String()
		a.Status = "unknown"
	}

	// Set the Country Emoji
	flag, err := emoji.CountryFlag(a.ISOCountry)
	if err != nil {
		return a, fmt.Errorf("unable to lookup country emoji - %w", err)
	}
	a.Emoji = flag.String()

	// Verify required data is present
	if a.LocalCode == "" || a.Name == "" || a.ISOCountry == "" {
		return airport.Airport{}, ErrNotEnoughFields
	}

	return a, nil
}
