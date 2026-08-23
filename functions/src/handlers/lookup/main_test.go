package main

import (
	"encoding/json"
	"errors"
	"strings"
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
			name: "Valid",
			input: []byte(
				`[{"local_code":"bG9jYWxfY29kZQ==", "name":"bmFtZQ==", "country":"Y291bnRyeQ==", "emoji":"ZW1vamk=", "type":"dHlwZQ==", "type_emoji":"dHlwZV9lbW9qaQ==", "status":"c3RhdHVz"}]`,
			),
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
			name: "Invalid",
			input: []byte(
				`[{"local_code":"bG9jYWxfY29kZQ==", "name":"INVALID", "country":"Y291bnRyeQ==", "emoji":"ZW1vamk=", "type":"dHlwZQ==", "type_emoji":"dHlwZV9lbW9qaQ==", "status":"c3RhdHVz"}]`,
			),
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

type ResponseHelperTestCase struct {
	name          string
	payload       string
	expectedParts []string
}

func TestSuccessResponse(t *testing.T) {
	tt := []ResponseHelperTestCase{
		{
			name:    "WrapsAirportPayload",
			payload: `{"local_code":"PHX","name":"Phoenix Sky Harbor International Airport"}`,
			expectedParts: []string{
				`"ok":true`,
				`"source":"sql"`,
				`"airport":{"local_code":"PHX","name":"Phoenix Sky Harbor International Airport"}`,
			},
		},
		{
			name:    "QuotesInvalidPayload",
			payload: `{"local_code":"PHX"`,
			expectedParts: []string{
				`"ok":true`,
				`"source":"sql"`,
				`"airport":"{\"local_code\":\"PHX\""`,
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			result := string(successResponse("sql", tc.payload))
			for _, part := range tc.expectedParts {
				if !strings.Contains(result, part) {
					t.Fatalf("expected response to contain %q, got %q", part, result)
				}
			}
		})
	}
}

func TestErrorResponse(t *testing.T) {
	tt := []ResponseHelperTestCase{
		{
			name:    "IncludesStageAndError",
			payload: "local_code is required",
			expectedParts: []string{
				`"ok":false`,
				`"stage":"validation"`,
				`"local_code":""`,
				`"error":"local_code is required"`,
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			result := string(errorResponse("validation", "", errors.New(tc.payload)))
			for _, part := range tc.expectedParts {
				if !strings.Contains(result, part) {
					t.Fatalf("expected response to contain %q, got %q", part, result)
				}
			}
		})
	}
}

type ValidateLocalCodeTestCase struct {
	name      string
	localCode string
	err       bool
}

func TestValidateLocalCode(t *testing.T) {
	tt := []ValidateLocalCodeTestCase{
		{name: "ValidUpper", localCode: "PHX"},
		{name: "ValidMixed", localCode: "ab12Cd"},
		{name: "InvalidChars", localCode: `A"OR1=1`, err: true},
		{name: "InvalidLength", localCode: "ABCDEFGHI", err: true},
		{name: "InvalidEmpty", localCode: "", err: true},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			_, err := validateLocalCode(tc.localCode)
			if tc.err && err == nil {
				t.Fatalf("expected error for local_code=%q, got nil", tc.localCode)
			}
			if !tc.err && err != nil {
				t.Fatalf("expected no error for local_code=%q, got %v", tc.localCode, err)
			}
		})
	}
}

func TestErrorResponseEscapesJSON(t *testing.T) {
	rawErr := errors.New("bad \"quote\" and \\ slash\nnext line")
	data := errorResponse("validation", "PHX", rawErr)

	var payload map[string]any
	if err := json.Unmarshal(data, &payload); err != nil {
		t.Fatalf("expected valid json, got error: %v", err)
	}

	if payload["ok"] != false {
		t.Fatalf("expected ok=false, got %v", payload["ok"])
	}
	if payload["stage"] != "validation" {
		t.Fatalf("expected stage=validation, got %v", payload["stage"])
	}
	if payload["local_code"] != "PHX" {
		t.Fatalf("expected local_code=PHX, got %v", payload["local_code"])
	}
	if payload["error"] != rawErr.Error() {
		t.Fatalf("expected error message to round-trip, got %v", payload["error"])
	}
}

func TestErrorResponseEscapesControlBytesAsJSON(t *testing.T) {
	rawErr := errors.New("bad \x00 byte")
	data := errorResponse("validation", "PHX", rawErr)

	var payload map[string]any
	if err := json.Unmarshal(data, &payload); err != nil {
		t.Fatalf("expected valid json, got error: %v", err)
	}

	if payload["error"] != rawErr.Error() {
		t.Fatalf("expected error message to round-trip, got %v", payload["error"])
	}
}

func TestParseLocalCode(t *testing.T) {
	_, err := parseLocalCode([]byte(`{}`))
	if err == nil {
		t.Fatalf("expected error when local_code is missing")
	}

	code, err := parseLocalCode([]byte(`{"local_code":"PHX"}`))
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if code != "PHX" {
		t.Fatalf("expected PHX, got %q", code)
	}
}

func TestTrimStagePrefix(t *testing.T) {
	err := trimStagePrefix(errors.New("sql_query: airport not found"))
	if err.Error() != "airport not found" {
		t.Fatalf("expected trimmed error, got %q", err.Error())
	}

	err = trimStagePrefix(errors.New("decode: invalid payload"))
	if err.Error() != "invalid payload" {
		t.Fatalf("expected trimmed error, got %q", err.Error())
	}
}
