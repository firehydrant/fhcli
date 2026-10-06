package cli

import (
	"context"
	"strings"
	"time"

	"github.com/firehydrant/firehydrant-go-sdk/models/components"
	"github.com/urfave/cli/v3"
)

func newEventCmd() *cli.Command {
	return &cli.Command{
		Name:   "event",
		Usage:  "submit a change event to the FireHydrant API",
		Action: eventCmd,
		Flags:  sharedFlags,
	}
}

func eventCmd(ctx context.Context, cmd *cli.Command) error {
	client, err := getClient(ctx)
	if err != nil {
		return handleError(cmd, err)
	}

	summary := strings.Join(cmd.Args().Slice(), " ")

	identities := parseKV(cmd.String("identities"))
	changeIdentities := make([]components.ChangeIdentity, 0, len(identities))
	for k, v := range identities {
		changeIdentities = append(changeIdentities, components.ChangeIdentity{
			Type:  k,
			Value: v,
		})
	}

	now := time.Now()
	req := components.CreateChangeEvent{
		Summary:          summary,
		Environments:     paramToList(cmd.String("environment")),
		Services:         paramToList(cmd.String("service")),
		StartsAt:         &now,
		EndsAt:           &now,
		Labels:           parseKV(cmd.String("labels")),
		ChangeIdentities: changeIdentities,
	}

	resp, err := client.Changes.CreateChangeEvent(ctx, req)
	if err != nil {
		return handleError(cmd, err)
	}

	return outputResult(cmd, resp)
}
