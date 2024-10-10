[![CI](https://github.com/infrasonar/linux-agent/workflows/CI/badge.svg)](https://github.com/infrasonar/linux-agent/actions)
[![Release Version](https://img.shields.io/github/release/infrasonar/linux-agent)](https://github.com/infrasonar/linux-agent/releases)

# InfraSonar Linux Agent

Documentation: https://docs.infrasonar.com/collectors/agents/linux/

## Environment variables

Environment                 | Default                       | Description
----------------------------|-------------------------------|-------------------
`STORAGE_PATH`              | `HOME/.infrasonar/`           | Path where files are stored
`TOKEN`                     | _required_                    | Token to connect to.
`ASSET_NAME`                | _none_                        | Initial Asset Name. This will only be used at the announce. Once the asset is created, `ASSET_NAME` will be ignored.
`ASSET_ID`                  | `/data/.asset.json`           | Asset Id _or_ file where the Agent asset Id is stored _(must be a volume mount)_.
`API_URI`                   | https://api.infrasonar.com    | InfraSonar API.
`CHECK_XXX_INTERVAL`        | `300`                         | Interval for the docker containers check in seconds.
`VERIFY_SSL`                | `1`                           | Verify SSL certificate, 0 _(=disabled)_ or 1 _(=enabled)_.

