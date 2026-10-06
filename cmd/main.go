package main

import (
	"context"
	"os"

	"github.com/firehydrant/fhcli/pkg/cli"
)

var commit string
var version string

func main() {
	cmd := cli.NewRootCmd(commit, version)
	// Commands report their own errors (JSON on stderr) and urfave/cli prints
	// usage errors, so just set the exit code here.
	if err := cmd.Run(context.Background(), os.Args); err != nil {
		os.Exit(1)
	}
}
