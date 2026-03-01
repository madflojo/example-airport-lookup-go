/*
This function is a lookup request handler function. When a request is received, this function is called and will parse
the incoming request, look up the requested airport from the cache, on cache miss, find the requested data in the
database, and provide it back as a response to the request.
*/
package main

import (
	"encoding/base64"
	"fmt"
	"html"

	sdk "github.com/tarmac-project/sdk"
	"github.com/tarmac-project/sdk/kv"
	"github.com/tarmac-project/sdk/logging"
	sdksql "github.com/tarmac-project/sdk/sql"
	"github.com/valyala/fastjson"
)

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

	// Lookup the airport from cache and return if found
	cache, err := f.kv.Get(lc)
	if err == nil && len(cache) > 0 {
		return successResponse("cache", string(cache)), nil
	}
	if err != nil {
		f.logging.Debug(fmt.Sprintf("cache lookup miss or error for %s: %s", lc, err))
	}

	// If not in cache, lookup the airport from the database
	query := fmt.Sprintf(`SELECT * FROM airports WHERE local_code = "%s"`, html.EscapeString(lc))
	data, err := f.sql.Query(query)
	if err != nil {
		f.logging.Error(fmt.Sprintf("error querying database for %s: %s", lc, err))
		return errorResponse("sql_query", lc, err), nil
	}

	// Verify we got a result
	if len(data.Data) == 0 {
		return errorResponse("sql_query", lc, fmt.Errorf("airport not found")), nil
	}

	// Decode the data
	rsp, err := decodeData(data.Data)
	if err != nil {
		f.logging.Error(fmt.Sprintf("error decoding data for %s: %s", lc, err))
		return errorResponse("decode", lc, err), nil
	}

	// Create JSON response
	airport := fmt.Sprintf(`{"local_code": "%s", "name": "%s", "country": "%s", "emoji": "%s", "type": "%s", "type_emoji": "%s", "status": "%s"}`,
		rsp["local_code"], rsp["name"], rsp["country"], rsp["emoji"], rsp["type"], rsp["type_emoji"], rsp["status"])

	// If not previously in cache, add to cache
	err = f.kv.Set(lc, []byte(airport))
	if err != nil {
		f.logging.Error(fmt.Sprintf("error setting cache value: %s", err))
	}

	// Return the airport
	return successResponse("sql", airport), nil
}

func successResponse(source, airportJSON string) []byte {
	return []byte(fmt.Sprintf(`{"ok":true,"source":"%s","airport":%s}`, source, airportJSON))
}

func errorResponse(stage, localCode string, err error) []byte {
	return []byte(fmt.Sprintf(`{"ok":false,"stage":"%s","local_code":"%s","error":"%s"}`, stage, localCode, err.Error()))
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
