[![CI](https://github.com/infrasonar/go-libagent/workflows/CI/badge.svg)](https://github.com/infrasonar/go-libagent/actions)
[![Release Version](https://img.shields.io/github/release/infrasonar/go-libagent)](https://github.com/infrasonar/go-libagent/releases)

# Go library for building InfraSonar Probes

This library is created for building [InfraSonar](https://infrasonar.com) probes.

## Environment variables

Environment                 | Default                       | Description
----------------------------|-------------------------------|-------------------
`STORAGE_PATH`              | `HOME/.infrasonar/`           | Path where files are stored _(not used when `ASSET_ID` is set)_.
`TOKEN`                     | _required_                    | Token used for authentication _(This MUST be a container token)_.
`ASSET_NAME`                | _none_                        | Initial Asset Name. This will only be used at the announce. Once the asset is created, `ASSET_NAME` will be ignored.
`ASSET_ID`                  | _none_                        | Asset Id _(If not given, the asset Id will be stored and loaded from file)_.
`API_URI`                   | https://api.infrasonar.com    | InfraSonar API.
`CHECK_XXX_INTERVAL`        | `300`                         | Interval in seconds for the `xxx` check. _(should be one environment variable for each check)_


```golang
package main

import (
	"log"

	"github.com/infrasonar/go-libagent"
)

const version = "0.1.0"

func CheckSample() (map[string][]map[string]any, error) {
	state := map[string][]map[string]any{}

	// Here code to create a check state.

    // Returning with an error will result in an InfraSonar Notification with
    // a check error.

    // Both a state and an error may be returned.

	return state, nil
}


func main() {
	// Start collector
	log.Printf("Starting InfraSonar Example Agent Collector v%s\n", version)

	// Initialize random
	libagent.RandInit()

	// Initialize Helper
	libagent.GetHelper()

	// Set-up signal handler
	quit := make(chan bool)
	go libagent.SigHandler(quit)

	// Create Collector
	collector := libagent.NewCollector("example", version)

	// Create Asset
	asset := libagent.NewAsset(collector)
	//asset.Kind = "Asset" // Optionally, set the asset Kind
	asset.Announce()

	// Create and plan checks
	checkSample := libagent.Check{
		Key:          "sample",
		Collector:    collector,
		Asset:        asset,
		IntervalEnv:  "CHECK_SAMPLE_INTERVAL",
		NoCount:      false,
		SetTimestamp: false,
		Fn:           CheckSample,
	}
	go checkSample.Plan(quit)

	// Wait for quit
	<-quit
}
```
