package main

import (
	"fmt"
	"github.com/tarmac-project/tarmac/pkg/sdk"
)

type Function struct {
	tarmac *sdk.Tarmac
}

func (f *Function) Handler(payload []byte) ([]byte, error) {
	return []byte("Hello World"), fmt.Errorf("not implemented http handler")
	// Parse the incoming request

	// Lookup the airport from cache

	// If not in cache, lookup the airport from the database

	// If not previously in cache, add to cache

	// If nothing found, return an error

	// Return the airport
	return nil, nil
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
