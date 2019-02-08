# FH CLI

This client allows you to submit free-form change events to the FireHydrant API from your infrastructure.

## Example Events
* Deploys
* Docker image builds
* Terraform runs
* Cronjobs

## Installation

## Usage

`fhcli [action] [parameters] [command]`

### Actions

* `event` Submits a change event to the FH API
* `init` Configures the client for subsequent invocations

### Parameters

* Environment Name `--environment`
* Service `--service`
* API Key `--api-key`
* Event name `--name`
* Payload `--payload` 
* Duration `--duration`
* Config file `--config`
* Identity `--identity`
* Git Revision `--revision`
  If the `fhcli` tool is invoked in a git repository this will automatically be set.
* Command
  This command is run in a subshell and is automatically instrumented with a duration/exit code. The exit code returned by fhcli is that of the command.

### Environment Variables

* `FH_ENVIRONMENT_NAME`
* `FH_COMPONENT`
* `FH_API_KEY`
* `FH_EVENT_NAME`

### Configuration file

We also accept a YAML configuration file, found at `~/.fhcli` or passed in as a param or environment variable.

### Examples

The following examples assume that `FH_API_KEY` is set in your environment.

    fhcli event --name "hourly reconciliation task" --environment "production" --duration=120ms
    fhcli event --name "Build image" --environment "staging" --git-revision=a0b0c0d0 docker build .
