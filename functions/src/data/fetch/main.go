/*
This function fetches airport data from a remote CSV file via HTTP and returns it to the application. It is designed
for reuse throughout the application.
*/
package main

import (
	"fmt"
	"github.com/tarmac-project/tarmac/pkg/sdk"
)

type Function struct {
	tarmac *sdk.Tarmac
}

func (f *Function) Handler(_ []byte) ([]byte, error) {
	f.tarmac.Logger.Info("Downloading airports.csv")
	rsp, err := f.tarmac.HTTP.Get("https://raw.githubusercontent.com/davidmegginson/ourairports-data/main/airports.csv")
	if err != nil {
		return []byte(""), fmt.Errorf("failed to get airports.csv: %w", err)
	}
	f.tarmac.Logger.Info(fmt.Sprintf("airports.csv downloaded with return code: %d", rsp.StatusCode))

	if rsp.StatusCode >= 299 {
		f.tarmac.Logger.Error(fmt.Sprintf("airports.csv download failed with return code: %d", rsp.StatusCode))
		return []byte(""), fmt.Errorf("failed to get airports.csv: HTTP request returned %d", rsp.StatusCode)
	}

	return rsp.Body, nil
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
