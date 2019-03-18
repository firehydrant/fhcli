package cli

import (
	"errors"
	"fmt"
	"os"
	"path"

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

var flags = []cli.Flag{
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

func NewApp(commit string, version string) *cli.App {
	app := cli.NewApp()

	app.Commands = []cli.Command{
		{
			Name:   "init",
			Usage:  "write a config file to be used in future invocations of the FireHydrant CLI",
			Action: initCmd,
			Flags:  sharedFlags,
		},
		{
			Name:   "event",
			Usage:  "submit an event to the FireHydrant API",
			Action: eventCmd,
			Flags:  sharedFlags,
			Before: parseConfigFile,
		},
		{
			Name:   "execute",
			Usage:  "execute a command and submit metrics to the FireHydrant API",
			Action: executeCmd,
			Flags:  sharedFlags,
			Before: parseConfigFile,
		},
	}

	app.Before = parseConfigFile
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

	if len(config.ApiHost) == 0 {
		return fhclient.ApiClient{}, errors.New("Invalid or no API host provided")
	}

	if len(config.ApiKey) == 0 {
		return fhclient.ApiClient{}, errors.New("Invalid or no API key provided")
	}

	return fhclient.NewApiClient(config), nil
}

func parseConfigFile(c *cli.Context) error {
	// Ignore the empty or missing configuration file if we're creating one
	if len(c.Args()) == 0 || c.Args()[0] == "init" {
		return nil
	}

	paths := []string{
		c.GlobalString("config"),
		path.Join("/etc", "firehydrant.cfg"),
		path.Join(os.Getenv("HOME"), "firehydrant.cfg"),
		path.Join("/tmp", "firehydrant.cfg"),
	}

	p := ""
	for _, pa := range paths {
		if _, err := os.Stat(pa); err == nil {
			p = pa
			break
		}
	}

	if p != "" {
		s, err := altsrc.NewYamlSourceFromFile(p)
		if err != nil {
			return err
		}

		err = altsrc.ApplyInputSourceValues(c, s, flags)
		if err != nil {
			return nil
		}

		return altsrc.ApplyInputSourceValues(c, s, sharedFlags)
	}

	return nil
}

func handleError(c *cli.Context, err error) error {
	fmt.Println(err.Error())
	if c.GlobalBool("ignoreErrors") {
		return nil
	}
	return err
}
