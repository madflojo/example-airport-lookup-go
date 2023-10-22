package airport

type Airport struct {
	// Continent is the continent the airport is located in.
	Continent string `json:"continent"`

	// Emoji is the emoji of the airport location. (Example: 🇺🇸, 🇨🇦, 🇬🇧, etc.)
	Emoji string `json:"emoji"`

	// GPSCode is the GPS code of the airport.
	GPSCode int `json:"gps_code"`

	// IATACode is the IATA code of the airport. This is a three letter code and is unique to each airport.
	IATACode string `json:"iata_code"`

	// ISOCountry is the ISO code of the country the airport is located in.
	ISOCountry string `json:"iso_country"`

	// ISORegion is the ISO code of the region the airport is located in.
	ISORegion string `json:"iso_region"`

	// LocalCode is the local code of the airport. This is a three letter code and is unique to each airport.
	LocalCode string `json:"local_code"`

	// Municipality is the municipality the airport is located in.
	Municipality string `json:"municipality"`

	// Name is the name of the airport.
	Name string `json:"name"`

	// Type is the type of airport (examples: small_airport, heliport, closed, etc.).
	Type string `json:"type"`

	// TypeEmoji is the emoji of the airport type (examples: 🛬, 🚁, 🚧, etc.).
	TypeEmoji string `json:"type_emoji"`

	// Status is the status of airport (examples: open, closed, etc.).
	Status string `json:"status"`
}
