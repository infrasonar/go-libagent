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
