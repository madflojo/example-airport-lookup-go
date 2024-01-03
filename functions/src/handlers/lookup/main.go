/*
This function is a lookup request handler function. When a request is received, this function is called and will parse 
the incoming request, look up the requested airport from the cache, on cache miss, find the requested data in the 
database, and provide it back as a response to the request.
*/
package main

import (
	"encoding/base64"
	"fmt"
	"github.com/tarmac-project/tarmac/pkg/sdk"
	"github.com/valyala/fastjson"
	"html"
)

// Function is the main function object that will be initialized
// and called by the Tarmac SDK.
type Function struct {
	tarmac *sdk.Tarmac
}

// Handler is the entry point for the function and will be called
// when a request is received.
func (f *Function) Handler(payload []byte) ([]byte, error) {
	// Parse the incoming request
	lc := fastjson.GetString(payload, "local_code")
	if lc == "" {
		return []byte(`{"error": "local_code is required"}`), fmt.Errorf("local_code is required")
	}

	// Lookup the airport from cache and return if found
	cache, err := f.tarmac.KV.Get(lc)
	if err == nil && len(cache) > 0 {
		return cache, nil
	}

	// If not in cache, lookup the airport from the database
	query := fmt.Sprintf(`SELECT * FROM airports WHERE local_code = "%s"`, html.EscapeString(lc))
	data, err := f.tarmac.SQL.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error querying database: %w", err)
	}

	// Verify we got a result
	if len(data) == 0 {
		return []byte(`{"error": "airport not found"}`), fmt.Errorf("airport not found")
	}

	// Decode the data
	rsp, err := decodeData(data)
	if err != nil {
		return nil, fmt.Errorf("error decoding data: %w", err)
	}

	// Create JSON response
	j := fmt.Sprintf(`{"local_code": "%s", "name": "%s", "country": "%s", "emoji": "%s", "type": "%s", "type_emoji": "%s", "status": "%s"}`,
		rsp["local_code"], rsp["name"], rsp["country"], rsp["emoji"], rsp["type"], rsp["type_emoji"], rsp["status"])

	// If not previously in cache, add to cache
	err = f.tarmac.KV.Set(lc, []byte(j))
	if err != nil {
		f.tarmac.Logger.Error(fmt.Sprintf("error setting cache value: %s", err))
	}

	// Return the airport
	return []byte(j), nil
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

func main() {
	var err error

	// Initialize Function
	f := &Function{}

	// Initialize the Tarmac SDK
	f.tarmac, err = sdk.New(sdk.Config{
		Namespace: "tarmac",
		Handler:   f.Handler,
	})
	if err != nil {
		return
	}
}
