# Events API - Splunk

This is the source code to retrieve events from the 1Password Events API and is meant to be used with the 1Password Splunk app.

## Requirements

### Go

If you do not have `go` locally installed, you can find installation steps [here](https://golang.org/doc/install).

## Installation

This application assumes it is running with the associated Splunk configuration files. These are found in `events-api-splunk/app`.
<br/><br/>
Follow the compilation steps found in the this projects root `README.md` to build the source code and move the binaries to their appropriate location.
<br/><br/>
If you do not wish to run this application from within Splunk, you will need to create the necessary configuration and storage files as well as update the source code to point to these newly created files. This is currently _not_ supported by default.

## Program Flow

On startup the `main.go` file will read the configuration file for various values, most importantly your 1Password Events API token. It will then check which scopes your token provides and from there continue down an event type path. Each event type path follows the same logic.
<br/><br/>
First it checks if there is a cursor file which is a temporary history of the last requests to the API. If one exists, it will use this cursor to get the next set of events, and on success update the history. The program will also roll the history file to control the file's length. If there is no history, a new file is created and the limit and start date values are used to create the first request to the 1Password Events API. The cursor file is also updated after this request.

## Directory Structure

The entry point to the application is found in `main.go`. There are various packages that the main file depends on:

### actions/

This folder contains all the business logic explained in the Program Flow above for each event type.

### api/

This folder contains the client methods to interact with the 1Password Events API.

### store/

This folder contains source code for interacting with the cursor stores such as opening the file, closing the file, getting the last cursor, updating the history and rolling the cursor file.

### vendor/

This folder contains all third party dependencies.
