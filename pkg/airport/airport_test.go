package airport

import (
	"reflect"
	"testing"
)

type AirportTestCase struct {
	name     string
	input    Airport
	expected Airport
	err      error
}

func TestAirport(t *testing.T) {
	tt := []AirportTestCase{
		{
			name: "Basic Test",
			input: Airport{
				Name:       "Test Airport",
				LocalCode:  "TST",
				ISOCountry: "US",
			},
			expected: Airport{
				Name:       "Test Airport",
				LocalCode:  "TST",
				ISOCountry: "US",
				Emoji:      "🇺🇸",
				TypeEmoji:  "❓",
				Status:     "unknown",
			},
		},
		{
			name: "Missing name",
			input: Airport{
				LocalCode:  "TST",
				ISOCountry: "US",
			},
			err: ErrMissingFields,
		},
		{
			name: "Unknown country",
			input: Airport{
				Name:       "Test Airport",
				LocalCode:  "TST",
				ISOCountry: "TST",
			},
			err: ErrUnknownCountry,
		},
		{
			name: "Heliport",
			input: Airport{
				Name:       "Test Heliport",
				LocalCode:  "TST",
				ISOCountry: "US",
				Type:       "heliport",
			},
			expected: Airport{
				Name:       "Test Heliport",
				LocalCode:  "TST",
				ISOCountry: "US",
				Emoji:      "🇺🇸",
				TypeEmoji:  "🚁",
				Status:     "open",
				Type:       "heliport",
			},
		},
		{
			name: "Small Airport",
			input: Airport{
				Name:       "Test Small Airport",
				LocalCode:  "TST",
				ISOCountry: "US",
				Type:       "small_airport",
			},
			expected: Airport{
				Name:       "Test Small Airport",
				LocalCode:  "TST",
				ISOCountry: "US",
				Emoji:      "🇺🇸",
				TypeEmoji:  "🛩️",
				Status:     "open",
				Type:       "small_airport",
			},
		},
		{
			name: "Medium Airport",
			input: Airport{
				Name:       "Test Medium Airport",
				LocalCode:  "TST",
				ISOCountry: "US",
				Type:       "medium_airport",
			},
			expected: Airport{
				Name:       "Test Medium Airport",
				LocalCode:  "TST",
				ISOCountry: "US",
				Emoji:      "🇺🇸",
				TypeEmoji:  "✈️",
				Status:     "open",
				Type:       "medium_airport",
			},
		},
		{
			name: "Large Airport",
			input: Airport{
				Name:       "Test Large Airport",
				LocalCode:  "TST",
				ISOCountry: "US",
				Type:       "large_airport",
			},
			expected: Airport{
				Name:       "Test Large Airport",
				LocalCode:  "TST",
				ISOCountry: "US",
				Emoji:      "🇺🇸",
				TypeEmoji:  "✈️",
				Status:     "open",
				Type:       "large_airport",
			},
		},
		{
			name: "Balloonport",
			input: Airport{
				Name:       "Test Balloonport",
				LocalCode:  "TST",
				ISOCountry: "US",
				Type:       "balloonport",
			},
			expected: Airport{
				Name:       "Test Balloonport",
				LocalCode:  "TST",
				ISOCountry: "US",
				Emoji:      "🇺🇸",
				TypeEmoji:  "🎈",
				Status:     "open",
				Type:       "balloonport",
			},
		},
		{
			name: "Seaplane Base",
			input: Airport{
				Name:       "Test Seaplane Base",
				LocalCode:  "TST",
				ISOCountry: "US",
				Type:       "seaplane_base",
			},
			expected: Airport{
				Name:       "Test Seaplane Base",
				LocalCode:  "TST",
				ISOCountry: "US",
				Emoji:      "🇺🇸",
				Status:     "open",
				Type:       "seaplane_base",
				TypeEmoji:  "⚓",
			},
		},
		{
			name: "Closed Airport",
			input: Airport{
				Name:       "Test Closed Airport",
				LocalCode:  "TST",
				ISOCountry: "US",
				Type:       "closed",
			},
			expected: Airport{
				Name:       "Test Closed Airport",
				LocalCode:  "TST",
				ISOCountry: "US",
				Emoji:      "🇺🇸",
				Status:     "closed",
				Type:       "unknown",
				TypeEmoji:  "🚧",
			},
		},
	}

	for _, tc := range tt {
		t.Run("Airport Validation: "+tc.name, func(t *testing.T) {
			result, err := Validate(tc.input)
			if err != nil {
				if tc.err == nil {
					t.Fatalf("Unexpected error: %v", err)
				}
				if tc.err != err {
					t.Fatalf("Expected error %v, got %v", tc.err, err)
				}
				return
			}

			if !reflect.DeepEqual(result, tc.expected) {
				t.Errorf("Expected %v, got %v", tc.expected, result)
			}
		})
	}
}
