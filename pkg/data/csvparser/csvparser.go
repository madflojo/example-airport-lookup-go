package csvparser

import (
  "github.com/tarmac-project/example-airport-lookup/pkg/data/airport"
)

type Airport struct {
  // IATACode is the IATA code of the airport. This is a three letter code and is unique to each airport.
  IATACode string

  // LocalCode is the local code of the airport. This is a three letter code and is unique to each airport.
  LocalCode string

  // Name is the name of the airport.
  Name string

  // Continent is the continent the airport is located in.
  Continent string

  // ISOCountry is the ISO code of the country the airport is located in.
  ISOCountry string

  // ISORegion is the ISO code of the region the airport is located in.
  ISORegion string

  // Municipality is the municipality the airport is located in.
  Municipality string

  // GPSCode is the GPS code of the airport.
  GPSCode string

  // AirportType is the type of airport (examples: small_airport, heliport, etc.).
  AirportType string
}

type Parser struct {}

func New() *Parser {
  return &Parser{}
}

func (p *Parser) ParseAirport(data []byte) ([]Airport, error) {
  return nil, nil
}
