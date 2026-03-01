/*
This function is a lookup request handler function. When a request is received, this function is called and will parse
the incoming request, look up the requested airport from the cache, on cache miss, find the requested data in the
database, and provide it back as a response to the request.
*/
package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"regexp"

	sdk "github.com/tarmac-project/sdk"
	"github.com/tarmac-project/sdk/kv"
	"github.com/tarmac-project/sdk/logging"
	sdksql "github.com/tarmac-project/sdk/sql"
	"github.com/valyala/fastjson"
)

var localCodePattern = regexp.MustCompile(`^[A-Za-z0-9]{1,8}$`)

// Function is the main function object that will be initialized
// and called by the Tarmac SDK.
type Function struct {
	sdk     *sdk.SDK
	logging logging.Client
	kv      kv.Client
	sql     sdksql.Client
}

// Handler is the entry point for the function and will be called
// when a request is received.
func (f *Function) Handler(payload []byte) ([]byte, error) {
	// Parse the incoming request
	lc := fastjson.GetString(payload, "local_code")
	if lc == "" {
		return errorResponse("validation", "", fmt.Errorf("local_code is required")), nil
	}

	validatedLocalCode, err := validateLocalCode(lc)
	if err != nil {
		return errorResponse("validation", lc, err), nil
	}

	// Lookup the airport from cache and return if found
	cache, err := f.kv.Get(validatedLocalCode)
	if err == nil && len(cache) > 0 {
		return successResponse("cache", string(cache)), nil
	}
	if err != nil {
		f.logging.Debug(fmt.Sprintf("cache lookup miss or error for %s: %s", validatedLocalCode, err))
	}

	// If not in cache, lookup the airport from the database
	query := fmt.Sprintf(`SELECT * FROM airports WHERE local_code = "%s"`, validatedLocalCode)
	data, err := f.sql.Query(query)
	if err != nil {
		f.logging.Error(fmt.Sprintf("error querying database for %s: %s", validatedLocalCode, err))
		return errorResponse("sql_query", validatedLocalCode, err), nil
	}

	// Verify we got a result
	if len(data.Data) == 0 {
		return errorResponse("sql_query", validatedLocalCode, fmt.Errorf("airport not found")), nil
	}

	// Decode the data
	rsp, err := decodeData(data.Data)
	if err != nil {
		f.logging.Error(fmt.Sprintf("error decoding data for %s: %s", validatedLocalCode, err))
		return errorResponse("decode", validatedLocalCode, err), nil
	}

	// Create JSON response
	airport := struct {
		LocalCode string `json:"local_code"`
		Name      string `json:"name"`
		Country   string `json:"country"`
		Emoji     string `json:"emoji"`
		Type      string `json:"type"`
		TypeEmoji string `json:"type_emoji"`
		Status    string `json:"status"`
	}{
		LocalCode: rsp["local_code"],
		Name:      rsp["name"],
		Country:   rsp["country"],
		Emoji:     rsp["emoji"],
		Type:      rsp["type"],
		TypeEmoji: rsp["type_emoji"],
		Status:    rsp["status"],
	}

	airportJSON, err := json.Marshal(airport)
	if err != nil {
		f.logging.Error(fmt.Sprintf("error encoding airport response for %s: %s", validatedLocalCode, err))
		return errorResponse("encode", validatedLocalCode, err), nil
	}

	// If not previously in cache, add to cache
	err = f.kv.Set(validatedLocalCode, airportJSON)
	if err != nil {
		f.logging.Error(fmt.Sprintf("error setting cache value: %s", err))
	}

	// Return the airport
	return successResponse("sql", string(airportJSON)), nil
}

func successResponse(source, airportJSON string) []byte {
	airport := json.RawMessage(`{}`)
	if json.Valid([]byte(airportJSON)) {
		airport = json.RawMessage(airportJSON)
	} else {
		quotedAirportJSON, err := json.Marshal(airportJSON)
		if err == nil {
			airport = json.RawMessage(quotedAirportJSON)
		}
	}

	response := struct {
		Ok      bool            `json:"ok"`
		Source  string          `json:"source"`
		Airport json.RawMessage `json:"airport"`
	}{
		Ok:      true,
		Source:  source,
		Airport: airport,
	}

	data, err := json.Marshal(response)
	if err != nil {
		return []byte(`{"ok":true,"source":"unknown","airport":{}}`)
	}

	return data
}

func errorResponse(stage, localCode string, err error) []byte {
	errMessage := ""
	if err != nil {
		errMessage = err.Error()
	}

	response := struct {
		Ok        bool   `json:"ok"`
		Stage     string `json:"stage"`
		LocalCode string `json:"local_code"`
		Error     string `json:"error"`
	}{
		Ok:        false,
		Stage:     stage,
		LocalCode: localCode,
		Error:     errMessage,
	}

	data, marshalErr := json.Marshal(response)
	if marshalErr != nil {
		return []byte(`{"ok":false,"stage":"internal","local_code":"","error":"failed to encode error response"}`)
	}

	return data
}

func validateLocalCode(localCode string) (string, error) {
	if !localCodePattern.MatchString(localCode) {
		return "", fmt.Errorf("local_code must match [A-Za-z0-9] and be 1-8 characters")
	}

	return localCode, nil
}

// decodeData decodes the data from the database and returns a map of the fields
func decodeData(data []byte) (map[string]string, error) {
	fields := []string{"local_code", "name", "country", "emoji", "type", "type_emoji", "status"}
	rsp := make(map[string]string)

	for _, field := range fields {
		v, err := base64.StdEncoding.DecodeString(fastjson.GetString(data, "0", field))
		if err != nil {
			return nil, fmt.Errorf("error decoding %s: %w", field, err)
		}
		rsp[field] = string(v)
	}
	return rsp, nil
}

//go:wasmexport wapc_init
func Initialize() {
	var err error

	// Initialize Function
	f := &Function{}

	// Initialize the Tarmac SDK
	f.sdk, err = sdk.New(sdk.Config{
		Namespace: "tarmac",
		Handler:   f.Handler,
	})
	if err != nil {
		return
	}

	cfg := f.sdk.Config()

	f.logging, err = logging.New(logging.Config{
		SDKConfig: cfg,
	})
	if err != nil {
		return
	}

	f.kv, err = kv.New(kv.Config{
		SDKConfig: cfg,
	})
	if err != nil {
		return
	}

	f.sql, err = sdksql.New(sdksql.Config{
		SDKConfig: cfg,
	})
	if err != nil {
		return
	}
}
