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

// Function initializes database objects and startup seed data.
type Function struct {
	sdk      *sdk.SDK
	logging  logging.Client
	function functionsdk.Client
	sql      sdksql.Client
}

// Handler initializes SQL schema and seeds initial data.
func (f *Function) Handler(_ []byte) ([]byte, error) {
	f.logging.Info("initializing airport lookup service")

	if err := f.createTable(); err != nil {
		f.logging.Error(fmt.Sprintf("failed to create table: %v", err))
		return []byte(""), fmt.Errorf("failed to create table: %w", err)
	}
	f.logging.Info("created database table")

	loadRsp, err := f.seedBaselineData()
	if err != nil {
		f.logging.Error(fmt.Sprintf("failed to seed airport data: %v", err))
		return []byte(""), fmt.Errorf("failed to seed airport data: %w", err)
	}
	if len(loadRsp) == 0 {
		f.logging.Error("seed function returned empty summary payload")
		return []byte(""), fmt.Errorf("seed function returned empty summary payload")
	}
	f.logging.Info(fmt.Sprintf("seed function summary: %s", string(loadRsp)))
	f.logging.Info("seeded airport data")

	return []byte(""), nil
}

func (f *Function) createTable() error {
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
		return err
	}

	return nil
}

func (f *Function) seedBaselineData() ([]byte, error) {
	return f.function.Call("seed", []byte(""))
}

// Initialize sets up SDK and clients required by the init function.
//
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
