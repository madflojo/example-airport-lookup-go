package main

import (
	"encoding/base64"
	"fmt"
	"github.com/tarmac-project/tarmac/pkg/sdk"
	"github.com/valyala/fastjson"
	"html"
)

type Function struct {
	tarmac *sdk.Tarmac
}

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

	// Build a response
	rsp := make(map[string]string)

	// Grab fields from the database and base64 decode them (this would be much easier with encoding/json)
	v, err := base64.StdEncoding.DecodeString(fastjson.GetString(data, "0", "local_code"))
	if err != nil {
		return nil, fmt.Errorf("error decoding local_code: %w", err)
	}
	rsp["local_code"] = string(v)

	v, err = base64.StdEncoding.DecodeString(fastjson.GetString(data, "0", "name"))
	if err != nil {
		return nil, fmt.Errorf("error decoding name: %w", err)
	}
	rsp["name"] = string(v)

	v, err = base64.StdEncoding.DecodeString(fastjson.GetString(data, "0", "iso_country"))
	if err != nil {
		return nil, fmt.Errorf("error decoding country: %w", err)
	}
	rsp["country"] = string(v)

	v, err = base64.StdEncoding.DecodeString(fastjson.GetString(data, "0", "emoji"))
	if err != nil {
		return nil, fmt.Errorf("error decoding emoji: %w", err)
	}
	rsp["emoji"] = string(v)

	v, err = base64.StdEncoding.DecodeString(fastjson.GetString(data, "0", "type"))
	if err != nil {
		return nil, fmt.Errorf("error decoding type: %w", err)
	}
	rsp["type"] = string(v)

	v, err = base64.StdEncoding.DecodeString(fastjson.GetString(data, "0", "type_emoji"))
	if err != nil {
		return nil, fmt.Errorf("error decoding type_emoji: %w", err)
	}
	rsp["type_emoji"] = string(v)

	v, err = base64.StdEncoding.DecodeString(fastjson.GetString(data, "0", "status"))
	if err != nil {
		return nil, fmt.Errorf("error decoding status: %w", err)
	}
	rsp["status"] = string(v)

	// Create JSON response
	j := fmt.Sprintf(`{"local_code": "%s", "name": "%s", "country": "%s", "emoji": "%s", "type": "%s", "type_emoji": "%s", "status": "%s"}`,
		rsp["local_code"], rsp["name"], rsp["country"], rsp["emoji"], rsp["type"], rsp["type_emoji"], rsp["status"])

	// If not previously in cache, add to cache
	err = f.tarmac.KV.Set(lc, []byte(j))
	if err != nil {
		f.tarmac.Logger.Error(fmt.Sprintf("error setting cache value: %w", err))
	}

	// Return the airport
	return []byte(j), nil
}

func main() {
	var err error

	// Initialize Function
	f := &Function{}

	// Initialize the Tarmac SDK
	f.tarmac, err = sdk.New(sdk.Config{
		Namespace: "airport-lookup",
		Handler:   f.Handler,
	})
	if err != nil {
		return
	}
}
