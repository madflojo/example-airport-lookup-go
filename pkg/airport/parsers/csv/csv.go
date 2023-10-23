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

	// Create a slice of Airports
	var airports []airport.Airport

	// Read each line of the CSV file
	for {
		// Read the CSV file
		rec, err := reader.Read()
		if err != nil {
			// Breakout if we've reached the end of the file
			if err == io.EOF {
				break
			}
			return nil, fmt.Errorf("unable to read csv data - %w", err)
		}

		// Convert the record to an Airport
		a, err := RecordToAirport(rec)
		if err != nil {
			// Skip any record level errors
			continue
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
