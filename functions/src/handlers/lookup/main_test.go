package main

import (
	"testing"
)

type DecodeDataTestCase struct {
	name     string
	input    []byte
	expected map[string]string
	err      bool
}

func TestDecodeFields(t *testing.T) {
	tt := []DecodeDataTestCase{
		{
			name:  "Valid",
			input: []byte(`[{"local_code":"bG9jYWxfY29kZQ==", "name":"bmFtZQ==", "country":"Y291bnRyeQ==", "emoji":"ZW1vamk=", "type":"dHlwZQ==", "type_emoji":"dHlwZV9lbW9qaQ==", "status":"c3RhdHVz"}]`),
			expected: map[string]string{
				"local_code": "local_code",
				"name":       "name",
				"country":    "country",
				"emoji":      "emoji",
				"type":       "type",
				"type_emoji": "type_emoji",
				"status":     "status",
			},
		},
		{
			name:     "Invalid",
			input:    []byte(`[{"local_code":"bG9jYWxfY29kZQ==", "name":"INVALID", "country":"Y291bnRyeQ==", "emoji":"ZW1vamk=", "type":"dHlwZQ==", "type_emoji":"dHlwZV9lbW9qaQ==", "status":"c3RhdHVz"}]`),
			expected: nil,
			err:      true,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			result, err := decodeData(tc.input)
			if err != nil {
				if !tc.err {
					t.Errorf("Expected no error, got %v", err)
				}
				return
			}
			if err == nil && tc.err {
				t.Fatalf("Expected error, got none")
			}

			for k, v := range tc.expected {
				if result[k] != v {
					t.Errorf("Expected %s, got %s", v, result[k])
				}
			}
		})
	}
}
