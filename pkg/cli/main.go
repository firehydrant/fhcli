package cli

import (
	"errors"
	"fmt"

	"github.com/davecgh/go-spew/spew"
	"github.com/firehydrant/api-client-go/fhclient"
	"gopkg.in/urfave/cli.v1"
	"gopkg.in/urfave/cli.v1/altsrc"
)

var sharedFlags = []cli.Flag{
	altsrc.NewStringFlag(cli.StringFlag{
		Name:   "identities",
		Usage:  "identities for the event, comma separated. example: image=us.gcr.io/firehydrant/rails:abcdef0",
		EnvVar: "FH_IDENTITIES",
	}),
	altsrc.NewStringFlag(cli.StringFlag{
		Name:   "labels",
		Usage:  "labels for the event, comma separated. example: executor=jenkins",
		EnvVar: "FH_LABELS",
	}),
	altsrc.NewStringFlag(cli.StringFlag{
		Name:   "environment",
		Usage:  "environment for the event, example: production",
		EnvVar: "FH_ENVIRONMENT",
	}),
	altsrc.NewStringFlag(cli.StringFlag{
		Name:   "service",
		Usage:  "service for the event, example: rails-monolith",
		EnvVar: "FH_SERVICE",
	}),
}

func NewApp(commit string, version string) *cli.App {
	app := cli.NewApp()

	flags := []cli.Flag{
		altsrc.NewStringFlag(cli.StringFlag{
			Name:   "apiKey",
			Usage:  "firehydrant.io API Key",
			EnvVar: "FH_API_KEY",
		}),
		altsrc.NewStringFlag(cli.StringFlag{
			Name:   "apiHost",
			Usage:  "firehydrant.io API hostname",
			Value:  "api.firehydrant.io",
			EnvVar: "FH_API_HOST",
		}),
		altsrc.NewBoolFlag(cli.BoolFlag{
			Name:  "debug",
			Usage: "Enable debug output for API calls",
		}),
		altsrc.NewBoolTFlag(cli.BoolTFlag{
			Name:  "ignoreErrors",
			Usage: "exit 0 on errors from FH API (default)",
		}),
		cli.StringFlag{
			Name:   "config",
			Usage:  "config file",
			EnvVar: "FH_CONFIG_FILE",
		},
	}

	app.Commands = []cli.Command{
		{
			Name:   "event",
			Usage:  "submit an event to the FireHydrant API",
			Action: eventCmd,
			Flags:  sharedFlags,
			Before: func(c *cli.Context) error {
				s, err := altsrc.NewYamlSourceFromFile(c.GlobalString("config"))
				if err != nil {
					return err
				}
				return altsrc.ApplyInputSourceValues(c, s, sharedFlags)
			},
		},
		{
			Name:   "execute",
			Usage:  "execute a command and submit metrics to the FireHydrant API",
			Action: executeCmd,
			Flags:  sharedFlags,
			Before: func(c *cli.Context) error {
				s, err := altsrc.NewYamlSourceFromFile(c.GlobalString("config"))
				if err != nil {
					return err
				}
				return altsrc.ApplyInputSourceValues(c, s, sharedFlags)
			},
		},
	}

	app.Before = altsrc.InitInputSourceWithContext(flags, altsrc.NewYamlSourceFromFlagFunc("config"))
	app.Flags = flags
	app.Version = fmt.Sprintf("%s (%s)", version, commit)

	return app
}

func NewApiClient(c *cli.Context) (fhclient.ApiClient, error) {
	config := fhclient.Config{
		ApiHost: c.GlobalString("apiHost"),
		ApiKey:  c.GlobalString("apiKey"),
		Debug:   c.GlobalBool("debug"),
	}
	spew.Dump(config)

	if len(config.ApiHost) == 0 {
		return fhclient.ApiClient{}, errors.New("valid or no API host provided")
	}

	if len(config.ApiKey) == 0 {
		return fhclient.ApiClient{}, errors.New("Invalid or no API key provided")
	}

	return fhclient.NewApiClient(config), nil
}
