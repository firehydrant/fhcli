# FireHydrant CLI (fhcli)

This client allows you to submit free-form change events to the FireHydrant API from your infrastructure. When incidents are opened on FireHydrant, recent change events will be automatically suggested to the responding team.

## Example Events
* Deploys
* Docker image builds
* Terraform runs
* Cronjobs

## Installation

![GitHub release (latest by date)](https://img.shields.io/github/v/release/firehydrant/fhcli?display_name=release)

New tags automatically create a release using [goreleaser](https://goreleaser.com/) and can be found on the [releases page](https://github.com/firehydrant/fhcli/releases) of this repository.

Our releases follow the standard format of `fhcli_{version}_{os}_{arch}.tar.gz` - To grab the latest on your CI environment (as an example), you can run something like:

```
$ cd /usr/local/bin
$ wget -c https://github.com/firehydrant/fhcli/releases/latest/download/fhcli_0.0.14_linux_arm64.tar.gz -O - | tar -xz
$ fhcli --version
```

## Usage

`fhcli [action] [parameters] [command]`

### Actions

* `event` Submits a change event to the FH API
* `execute` Executes a command and submits a corresponding change event to the FH API
* `init` Configures the client for subsequent invocations

### Parameters

* Environment Name `--environment` `FH_ENVIRONMENT`
* Service `--service` `FH_SERVICE`
* API Key `--apiKey` `FH_API_KEY`
* Config file `--config` `FH_CONFIG_FILE`
* Identities `--identities` `FH_IDENTITIES`
* Event name
_or_
* Command
  This command is run in a subshell, automatically instrumented with a duration and the result is submitted to the FH API

### Environment Variables

### Configuration file

We also accept a YAML configuration file, specified with `-config` or passed in as an environment variable. If none is specified, the following paths will be checked in order for a configuration file.

* /etc/firehydrant.cfg
* $HOME/firehydrant.cfg
* /tmp/firehydrant.cfg

### Examples

The following examples assume that `FH_API_KEY` is set in your environment.

    fhcli event --environment "production" "hourly reconciliation task"
    fhcli event --environment "staging" --identities git=git@github.com/firehydrant/myrepo:a0b0c0d0 "CI build succeeded"
    fhcli execute --environment "staging" --identities revision=a0a0a0 -- docker build .

To configure firehydrant and write a config file (defaults to /tmp/firehydrant.cfg) for subsequent invocations of the tool:

    fhcli --config /etc/firehydrant.cfg --apiKey fhb-a0b010 --environment=production --service=monolith init
