package main

import (
	"fmt"
	"github.com/tarmac-project/tarmac/pkg/sdk"
)

type Function struct {
	tarmac *sdk.Tarmac
}

func (f *Function) Handler(_ []byte) ([]byte, error) {
	f.tarmac.Logger.Info("Airport raw data download starting")

	// Fetch the airport data
	_, err := f.tarmac.Function.Call("fetch", []byte(""))
	if err != nil {
		f.tarmac.Logger.Error(fmt.Sprintf("Failed to fetch airport data", err))
		return []byte(""), fmt.Errorf("Failed to fetch airport data: %s", err)
	}

	// Parse the data

	// Update the database

	return []byte(""), nil
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
