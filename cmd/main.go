package main

import (
	"log"
	"os"

	"github.com/firehydrant/fhcli/pkg/cli"
)

var commit string
var version string

func main() {
	app := cli.NewApp(commit, version)

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}

}
