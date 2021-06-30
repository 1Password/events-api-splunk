# Events API for Splunk

This is the Splunk add-on which is comprized of various splunk configuration files and go binaries, used to ingest 1Password events. This Add-on can be run locally in docker or built and installed on Splunk. For building the add-on, see the root README.md file.

## Requirements

### Go

If you do not have `go` locally installed, you can find installation steps [here](https://golang.org/doc/install).

### Docker

If you do not have `Docker` installed, you can find the installation steps [here](https://docs.docker.com/get-docker).

### Directory Structure

The top level files contain the commands and configuration to build and run the application in docker. All the necessary dependencies that comprise the Add-on can be found in `onepassword_events_api`.

## Running

Start by running `make`, a docker container should start with a Splunk image. Navigate to `http://localhost:8000/` to login. Use the following credentials to continue:

- username: admin
- password: hey1234567890

The 1Password icon should be shown under the Applications pane on the left. That viewable icon, the associated setup pages, and the ingestion script are all configurations/scripts included in the `onepassword_events_api` app.

## Development

- When making changes to the `appserver` files, cURL `http://localhost:8000/en-US/_bump` to get the updated files.
- When making changes to `.conf` files or dashboards, cURL `http://localhost:8000/en-US/debug/refresh` to get the updated files.
- When updating the binaries found in `onepassword_events_api/bin/`, you must restart Splunk.

## Debugging

The Splunk logs, as well as error logs triggered from the "scripted input" can be found in `splunkd.log`
