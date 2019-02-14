package cli

import (
	"errors"

	apiclient "github.com/firehydrant/fhcli/pkg/api_client"
	"github.com/urfave/cli"
)

var sharedFlags = []cli.Flag{
	cli.StringFlag{
		Name:   "identities, i",
		Usage:  "identities for the event, comma separated. example: image=us.gcr.io/firehydrant/rails:abcdef0",
		EnvVar: "FH_IDENTITIES",
	},
	cli.StringFlag{
		Name:   "labels",
		Usage:  "labels for the event, comma separated. example: executor=jenkins",
		EnvVar: "FH_LABELS",
	},
	cli.StringFlag{
		Name:   "environment",
		Usage:  "environment for the event, example: production",
		EnvVar: "FH_ENVIRONMENT",
	},
	cli.StringFlag{
		Name:   "service",
		Usage:  "service for the event, example: rails-monolith",
		EnvVar: "FH_SERVICE",
	},
}

func NewApp() *cli.App {
	app := cli.NewApp()

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "api-key, k",
			Usage:  "firehydrant.io API Key",
			EnvVar: "FH_API_KEY",
		},
		cli.StringFlag{
			Name:   "api-host, a",
			Usage:  "firehydrant.io API hostname",
			Value:  "api.firehydrant.io",
			EnvVar: "FH_API_HOST",
		},
		cli.BoolFlag{
			Name:  "debug, d",
			Usage: "Enable debug output for API calls",
		},
		cli.BoolTFlag{
			Name:  "ignore-errors, x",
			Usage: "exit 0 on errors from FH API (default)",
		},
	}

	app.Commands = []cli.Command{
		{
			Name:   "event",
			Usage:  "submit an event to the FireHydrant API",
			Action: eventCmd,
			Flags:  sharedFlags,
		},
		{
			Name:   "execute",
			Usage:  "execute a command and submit metrics to the FireHydrant API",
			Action: executeCmd,
			Flags:  sharedFlags,
		},
	}

	return app
}

func NewApiClient(c *cli.Context) (apiclient.ApiClient, error) {
	config := apiclient.Config{
		ApiHost: c.GlobalString("api-host"),
		ApiKey:  c.GlobalString("api-key"),
		Debug:   c.GlobalBool("debug"),
	}

	if len(config.ApiHost) == 0 {
		return apiclient.ApiClient{}, errors.New("Invalid or no API host provided")
	}

	if len(config.ApiKey) == 0 {
		return apiclient.ApiClient{}, errors.New("Invalid or no API key provided")
	}

	return apiclient.NewApiClient(config), nil
}
