/*
This function fetches airport data from a remote CSV file via HTTP and returns it to the application. It is designed
for reuse throughout the application.
*/
package main

import (
	"fmt"
	"io"

	sdk "github.com/tarmac-project/sdk"
	"github.com/tarmac-project/sdk/httpclient"
	"github.com/tarmac-project/sdk/logging"
)

const airportsCSVURL = "https://raw.githubusercontent.com/davidmegginson/ourairports-data/main/airports.csv"

type Function struct {
	sdk     *sdk.SDK
	logging logging.Client
	http    httpclient.Client
}

func (f *Function) Handler(_ []byte) ([]byte, error) {
	f.logging.Info("downloading airports.csv")

	rsp, err := f.http.Get(airportsCSVURL)
	if err != nil {
		f.logging.Error(fmt.Sprintf("failed to get airports.csv: %v", err))
		return nil, fmt.Errorf("failed to get airports.csv: %w", err)
	}

	f.logging.Info(fmt.Sprintf("airports.csv downloaded with return code: %d", rsp.StatusCode))

	if rsp.StatusCode >= 299 {
		f.logging.Error(fmt.Sprintf("airports.csv download failed with return code: %d", rsp.StatusCode))
		return nil, fmt.Errorf("failed to get airports.csv: http request returned %d", rsp.StatusCode)
	}

	if rsp.Body == nil {
		return []byte(""), nil
	}
	defer func() {
		if closeErr := rsp.Body.Close(); closeErr != nil {
			f.logging.Warn(fmt.Sprintf("failed to close airports.csv response body: %v", closeErr))
		}
	}()

	body, err := io.ReadAll(rsp.Body)
	if err != nil {
		f.logging.Error(fmt.Sprintf("failed to read airports.csv body: %v", err))
		return nil, fmt.Errorf("failed to read airports.csv body: %w", err)
	}

	return body, nil
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

	// Initialize Logger client
	f.logging, err = logging.New(logging.Config{
		SDKConfig: cfg,
	})
	if err != nil {
		return
	}

	// Initialize HTTP client
	f.http, err = httpclient.New(httpclient.Config{
		SDKConfig: cfg,
	})
	if err != nil {
		f.logging.Error(fmt.Sprintf("failed to create HTTP client: %v", err))
		return
	}
}
