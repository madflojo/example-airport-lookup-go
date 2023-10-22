package csv

import (
	"encoding/csv"
	"fmt"
	"github.com/tarmac-project/example-airport-lookup-go/pkg/airport"
	"io"
)

var (
	ErrIsHeader        = fmt.Errorf("record is a header")
	ErrNotEnoughFields = fmt.Errorf("not enough fields")
)

// Parser is a CSV parser for Airport data.
type Parser struct {
	reader io.Reader
}

// New creates a new Parser.
func New(reader io.Reader) (*Parser, error) {
	p := &Parser{
		reader: reader,
	}
	return p, nil
}

// Parse parses a CSV file of Airports returning a slice of Airports.
func (p *Parser) Parse() ([]airport.Airport, error) {
	// Create a new CSV Reader
	reader := csv.NewReader(p.reader)

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
		Type:         rec[2],
	}

	// Validate the Airport Data
	a, err := airport.Validate(a)
	if err != nil {
		return airport.Airport{}, fmt.Errorf("unable to validate airport - %w", err)
	}

	return a, nil
}
