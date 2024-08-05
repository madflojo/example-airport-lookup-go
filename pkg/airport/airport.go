package airport

import (
	"fmt"

	"github.com/enescakir/emoji"
)

// Airport is a struct that contains information about an airport.
type Airport struct {
	// Continent is the continent the airport is located in.
	Continent string `json:"continent"`

	// Emoji is the emoji of the airport location. (Example: 🇺🇸, 🇨🇦, 🇬🇧, etc.)
	Emoji string `json:"emoji"`

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

var (
	ErrMissingFields  = fmt.Errorf("Missing required fields")
	ErrUnknownCountry = fmt.Errorf("Unable to lookup country emoji")
)

// Validate validates the airport struct and returns the airport with the emoji and country flag set.
func Validate(a Airport) (Airport, error) {
	var err error
	// Validate Minimum Fields Exist
	if a.LocalCode == "" || a.Name == "" || a.ISOCountry == "" {
		return a, ErrMissingFields
	}

	// Set Airport Emoji
	a, err = setTypeEmoji(a)
	if err != nil {
		return a, err
	}

	// Set Country Flag
	a, err = setCountryFlag(a)
	if err != nil {
		return a, err
	}

	return a, nil
}

// setTypeEmoji sets the emoji for the airport type.
func setTypeEmoji(a Airport) (Airport, error) {
	switch a.Type {
	case "heliport":
		a.TypeEmoji = emoji.Helicopter.String()
		a.Status = "open"
	case "small_airport":
		a.TypeEmoji = emoji.SmallAirplane.String()
		a.Status = "open"
	case "medium_airport", "large_airport":
		a.TypeEmoji = emoji.Airplane.String()
		a.Status = "open"
	case "seaplane_base":
		a.TypeEmoji = emoji.Anchor.String()
		a.Status = "open"
	case "balloonport":
		a.TypeEmoji = emoji.Balloon.String()
		a.Status = "open"
	case "closed", "unknown":
		a.Type = "unknown"
		a.TypeEmoji = emoji.Construction.String()
		a.Status = "closed"
	default:
		a.TypeEmoji = emoji.QuestionMark.String()
		a.Status = "unknown"
	}

	return a, nil
}

// setCountryFlag sets the emoji for the country the airport is located in.
func setCountryFlag(a Airport) (Airport, error) {
	// Set the Country Emoji
	flag, err := emoji.CountryFlag(a.ISOCountry)
	if err != nil {
		return a, ErrUnknownCountry
	}
	a.Emoji = flag.String()
	return a, nil
}
