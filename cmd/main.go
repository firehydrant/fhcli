package main

import (
	"context"
	"log"
	"os"

	"github.com/firehydrant/fhcli/pkg/cli"
)

var commit string
var version string

func main() {
	cmd := cli.NewRootCmd(commit, version)
	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
