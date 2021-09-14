# Events Reporting for Splunk

This repository contains code to integrate Splunk with 1Password's Events API. It includes a Splunk Add-on that meets Splunk's AppInspect requirements, binary source code, and make commands to build the project.

## Directory Structure

The top level directory only contains two files, this `README.md` and a `Makefile` which has all the commands to build the Splunk Add-on for various distributions as well as build the linux specific version for running the application locally in docker.

### src/

This folder contains the go source and dependency code used in the Splunk Add-on. Changing this source code will not be reflected in the Splunk Add-on until you recompile the source. Use the `make compile_app_binary` to accomplish this.

### app/

This folder contains the necessary Splunk configuration files and compiled go source code. See this folder's `README.md` to learn about running the Splunk add-on locally.

### builds/

This folder will contain the OS specific Add-ons, compressed to Splunk's distribution requirements as well as installation steps.

## Requirements

### Go

If you do not have `go` locally installed, you can find installation steps [here](https://golang.org/doc/install).

## Commands

- `make compile_app_binary`
  This command will update the Splunk Add-on, found in `app`, with any changes made from `src`.

- `make new_version`
  This command will update the JS portion of the Splunk Add-on to `Makefile VERSION` and build a release bundle for the web app.

- `make build_all_binaries build_all_apps` will first compile the `src` code to various Operating System distributions, and then bundle them with the Splunk specific configurations (found in `app`). The output will be found in `builds/bin`.
