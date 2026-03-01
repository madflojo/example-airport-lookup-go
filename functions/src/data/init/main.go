/*
This function executes on startup to verify the database setup and load initial airport data.

The service will not start until the database is setup and the initial data is loaded.
*/
package main

import (
	"fmt"

	sdk "github.com/tarmac-project/sdk"
	functionsdk "github.com/tarmac-project/sdk/function"
	"github.com/tarmac-project/sdk/logging"
	sdksql "github.com/tarmac-project/sdk/sql"
)

type Function struct {
	sdk      *sdk.SDK
	logging  logging.Client
	function functionsdk.Client
	sql      sdksql.Client
}

func (f *Function) Handler(_ []byte) ([]byte, error) {
	f.logging.Info("Initializing Airport Lookup Service")

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
	_, err := f.sql.Exec(query)
	if err != nil {
		f.logging.Error(fmt.Sprintf("Failed to create table - %s", err))
		return []byte(""), fmt.Errorf("Failed to create table: %s", err)
	}
	f.logging.Info("Created database table")

	// Seed baseline airport data
	loadRsp, err := f.function.Call("seed", []byte(""))
	if err != nil {
		f.logging.Error(fmt.Sprintf("Failed to seed airport data - %s", err))
		return []byte(""), fmt.Errorf("Failed to seed airport data: %s", err)
	}
	if len(loadRsp) == 0 {
		f.logging.Warn("Seed function returned empty summary payload")
	} else {
		f.logging.Info(fmt.Sprintf("Seed function summary: %s", string(loadRsp)))
	}
	f.logging.Info("Seeded airport data")

	return []byte(""), nil
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

	f.function, err = functionsdk.New(functionsdk.Config{
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
