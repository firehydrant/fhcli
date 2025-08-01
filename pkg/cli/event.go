package cli

import (
	"fmt"
	"strings"
	"time"

	"github.com/firehydrant/api-client-go/client/changes"
	"github.com/firehydrant/api-client-go/fhclient"
	"github.com/firehydrant/api-client-go/models"
	"github.com/go-openapi/strfmt"
	"gopkg.in/urfave/cli.v1"
)

func eventCmd(c *cli.Context) error {
	client, err := NewApiClient(c)
	if err != nil {
		return err
	}

	params := changes.NewPostV1ChangesEventsParams()
	summary := strings.Join(c.Args(), " ")

	identities := fhclient.APIKVtoChangeIdentities(fhclient.MapToAPIKV(fhclient.ParseKV(c.String("identities"))))

	params.PostV1ChangesEvents = &models.PostV1ChangesEvents{
		Environments:     fhclient.ParamToList(c.String("environment")),
		Services:         fhclient.ParamToList(c.String("service")),
		StartsAt:         strfmt.DateTime(time.Now()),
		EndsAt:           strfmt.DateTime(time.Now()),
		Summary:          &summary,
		ChangeIdentities: identities,
		Labels:           fhclient.ParseKV(c.String("labels")),
	}

	resp, err := client.Client.Changes.PostV1ChangesEvents(params, client.Auth)

	if err != nil {
		return handleError(c, err)
	}

	fmt.Println(fmt.Sprintf("Created change event %s", resp.Payload.ID))
	return nil
}
