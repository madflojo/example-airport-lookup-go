package airport

type Airport struct {
	// Continent is the continent the airport is located in.
	Continent string

	// Emoji is the emoji of the airport location. (Example: 🇺🇸, 🇨🇦, 🇬🇧, etc.)
	Emoji string

	// GPSCode is the GPS code of the airport.
	GPSCode int

	// IATACode is the IATA code of the airport. This is a three letter code and is unique to each airport.
	IATACode string

	// ISOCountry is the ISO code of the country the airport is located in.
	ISOCountry string

	// ISORegion is the ISO code of the region the airport is located in.
	ISORegion string

	// LocalCode is the local code of the airport. This is a three letter code and is unique to each airport.
	LocalCode string

	// Municipality is the municipality the airport is located in.
	Municipality string

	// Name is the name of the airport.
	Name string

	// Status is the type and status of airport (examples: small_airport, heliport, closed, etc.).
	Status string
}
