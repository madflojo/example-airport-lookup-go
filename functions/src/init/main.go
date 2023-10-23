package main

import (
  "bytes"
	"fmt"
	"github.com/tarmac-project/tarmac/pkg/sdk"
  "github.com/tarmac-project/example-airport-lookup-go/pkg/airport/parsers/csv"
)

type Function struct {
	tarmac *sdk.Tarmac
}

func (f *Function) Handler(_ []byte) ([]byte, error) {
	f.tarmac.Logger.Info("Airport raw data download starting")

	// Fetch the airport data
	data, err := f.tarmac.Function.Call("fetch", []byte(""))
	if err != nil {
		f.tarmac.Logger.Error(fmt.Sprintf("Failed to fetch airport data", err))
		return []byte(""), fmt.Errorf("Failed to fetch airport data: %s", err)
	}

	// Parse the data
  parser, err := csv.New(bytes.NewReader(data))
  if err != nil {
    f.tarmac.Logger.Error(fmt.Sprintf("Failed to create CSV parser", err))
    return []byte(""), fmt.Errorf("Failed to create CSV parser: %s", err)
  }

  airports, err := parser.Parse()
  if err != nil {
    f.tarmac.Logger.Error(fmt.Sprintf("Failed to parse airport data", err))
    return []byte(""), fmt.Errorf("Failed to parse airport data: %s", err)
  }
  f.tarmac.Logger.Info(fmt.Sprintf("Fetched %d airports", len(airports)))

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
