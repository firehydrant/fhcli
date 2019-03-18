package cli

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/urfave/cli.v1"
	"gopkg.in/yaml.v2"
)

type configFile struct {
	APIKey       string `yaml:"apiKey"`
	APIHost      string `yaml:"apiHost"`
	Debug        bool   `yaml:"debug"`
	IgnoreErrors bool   `yaml:"ignoreErrors"`
	Identities   string `yaml:"identities"`
	Labels       string `yaml:"labels"`
	Environment  string `yaml:"environment"`
	Service      string `yaml:"service"`
}

func initCmd(c *cli.Context) error {

	out := c.GlobalString("config")
	if out == "" {
		out = "/tmp/firehydrant.cfg"
	}

	config := configFile{}
	if len(c.GlobalString("apiKey")) > 0 {
		config.APIKey = c.GlobalString("apiKey")
	}

	if len(c.GlobalString("apiHost")) > 0 {
		config.APIHost = c.GlobalString("apiHost")
	}

	config.Debug = c.GlobalBool("debug")
	config.IgnoreErrors = c.GlobalBool("ignoreErrors")

	if len(c.String("identities")) > 0 {
		config.Identities = c.String("identities")
	}

	if len(c.String("labels")) > 0 {
		config.Labels = c.String("labels")
	}

	if len(c.String("environment")) > 0 {
		config.Environment = c.String("environment")
	}

	if len(c.String("service")) > 0 {
		config.Service = c.String("service")
	}

	d, err := yaml.Marshal(config)

	if err != nil {
		return err
	}

	err = ioutil.WriteFile(out, d, 0644)
	if err != nil {
		return err
	}

	fmt.Println(fmt.Sprintf("Wrote config file %s", out))
	return nil
}
