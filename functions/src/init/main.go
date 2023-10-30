package main

import (
	"fmt"
	"github.com/tarmac-project/tarmac/pkg/sdk"
)

type Function struct {
	tarmac *sdk.Tarmac
}

func (f *Function) Handler(_ []byte) ([]byte, error) {
	f.tarmac.Logger.Info("Initializing Airport Lookup Service")

	// Create MySQL Database structure
	query := `CREATE TABLE IF NOT EXISTS airports (
    local_code VARCHAR(25) NOT NULL UNIQUE,
    name VARCHAR(255) NOT NULL,
    type VARCHAR(255) NOT NULL,
    type_emoji VARCHAR(255),
    continent VARCHAR(255),
    iso_country VARCHAR(255) NOT NULL,
    iso_region VARCHAR(255),
    municipality VARCHAR(255),
    emoji VARCHAR(255),
    status VARCHAR(255),
    PRIMARY KEY (local_code)
  );`
	_, err := f.tarmac.SQL.Query(query)
	if err != nil {
		f.tarmac.Logger.Error(fmt.Sprintf("Failed to create table - %s", err))
		return []byte(""), fmt.Errorf("Failed to create table: %s", err)
	}
	f.tarmac.Logger.Info("Created database table")

	// Load Airport Data
	_, err = f.tarmac.Function.Call("load", []byte(""))
	if err != nil {
		f.tarmac.Logger.Error(fmt.Sprintf("Failed to load airport data - %s", err))
		return []byte(""), fmt.Errorf("Failed to load airport data: %s", err)
	}
	f.tarmac.Logger.Info("Loaded airport data")

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
