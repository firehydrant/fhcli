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
* `exec` Executes a command and submits a corresponding change event to the FH API
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

We also accept a YAML configuration file, specified with `-config` or passed in as an environment variable. If none is specified, the following paths will be checked in order for a configuration file.

* /etc/firehydrant.cfg
* $HOME/firehydrant.cfg
* /tmp/firehydrant.cfg

### Examples

The following examples assume that `FH_API_KEY` is set in your environment.

    fhcli event --name "hourly reconciliation task" --environment "production" --duration=120ms
    fhcli event --name "Job run" --environment "staging" --identities git=git@github.com/firehydrant/myrepo:a0b0c0d0
    fhcli exec --name "Job run" --environment "staging" --identities git=git@github.com/firehydrant/myrepo:a0b0c0d0 docker build .

To configure firehydrant and write a config file (defaults to /tmp/firehydrant.cfg) for subsequent invocations of the tool:
    fhcli --config /etc/firehydrant.cfg --apiKey fhb-a0b010 --identities git=git@github.com/firehydrant/fhcli:a0b0c0d0 --service myservice init
