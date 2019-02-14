package cli

import (
	"fmt"
	"strings"

	changeevents "github.com/firehydrant/fhcli/pkg/change_events"
	"github.com/urfave/cli"
)

func eventCmd(c *cli.Context) error {
	client, err := NewApiClient(c)
	if err != nil {
		return err
	}

	ce := changeevents.NewChangeEvent()

	ce.Summary = strings.Join(c.Args(), " ")
	ce.Environment = c.String("environment")
	ce.Service = c.String("service")
	ce.RawIdentities = c.String("identities")
	ce.RawLabels = c.String("labels")

	id, err := ce.Submit(client)
	if err != nil {
		return err
	}

	fmt.Println(fmt.Sprintf("Created event %s", id))

	return nil
}
