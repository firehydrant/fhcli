package cli

import (
	"context"
	"errors"

	firehydrantgosdk "github.com/firehydrant/firehydrant-go-sdk"
	"github.com/firehydrant/firehydrant-go-sdk/models/components"
	"github.com/firehydrant/firehydrant-go-sdk/models/operations"
	"github.com/urfave/cli/v3"
)

func newIncidentsCmd() *cli.Command {
	return &cli.Command{
		Name:  "incidents",
		Usage: "manage FireHydrant incidents",
		Commands: []*cli.Command{
			{
				Name:   "list",
				Usage:  "list incidents",
				Action: incidentsListCmd,
				Flags: []cli.Flag{
					&cli.IntFlag{
						Name:  "page",
						Usage: "page number",
						Value: 1,
					},
					&cli.IntFlag{
						Name:  "per-page",
						Usage: "results per page",
						Value: 20,
					},
					&cli.StringFlag{
						Name:    "query",
						Aliases: []string{"q"},
						Usage:   "search query",
					},
					&cli.StringFlag{
						Name:  "status",
						Usage: "filter by status",
					},
					&cli.StringFlag{
						Name:  "severity",
						Usage: "filter by severity",
					},
				},
			},
			{
				Name:   "get",
				Usage:  "get incident details",
				Action: incidentsGetCmd,
			},
			{
				Name:   "note",
				Usage:  "post a status update note on an incident",
				Action: incidentsNoteCmd,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "body",
						Aliases:  []string{"b"},
						Usage:    "note body text",
						Required: true,
					},
				},
			},
			{
				Name:   "chat",
				Usage:  "post a chat message on an incident",
				Action: incidentsChatCmd,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "body",
						Aliases:  []string{"b"},
						Usage:    "message body text",
						Required: true,
					},
				},
			},
		},
	}
}

func incidentsListCmd(ctx context.Context, cmd *cli.Command) error {
	client, err := getClient(ctx)
	if err != nil {
		return handleError(cmd, err)
	}

	page := int(cmd.Int("page"))
	perPage := int(cmd.Int("per-page"))

	req := operations.ListIncidentsRequest{
		Page:    firehydrantgosdk.Int(page),
		PerPage: firehydrantgosdk.Int(perPage),
	}

	if q := cmd.String("query"); q != "" {
		req.Query = firehydrantgosdk.String(q)
	}
	if s := cmd.String("status"); s != "" {
		req.Status = firehydrantgosdk.String(s)
	}
	if s := cmd.String("severity"); s != "" {
		req.Severities = firehydrantgosdk.String(s)
	}

	resp, err := client.Incidents.ListIncidents(ctx, req)
	if err != nil {
		return handleError(cmd, err)
	}

	return outputResult(cmd, resp)
}

func incidentsGetCmd(ctx context.Context, cmd *cli.Command) error {
	client, err := getClient(ctx)
	if err != nil {
		return handleError(cmd, err)
	}

	if cmd.NArg() < 1 {
		return handleError(cmd, errors.New("incident ID is required"))
	}
	incidentID := cmd.Args().First()

	resp, err := client.Incidents.GetIncident(ctx, incidentID)
	if err != nil {
		return handleError(cmd, err)
	}

	return outputResult(cmd, resp)
}

func incidentsNoteCmd(ctx context.Context, cmd *cli.Command) error {
	client, err := getClient(ctx)
	if err != nil {
		return handleError(cmd, err)
	}

	if cmd.NArg() < 1 {
		return handleError(cmd, errors.New("incident ID is required"))
	}
	incidentID := cmd.Args().First()
	body := cmd.String("body")

	resp, err := client.Incidents.CreateIncidentNote(ctx, incidentID, components.CreateIncidentNote{
		Body: body,
	})
	if err != nil {
		return handleError(cmd, err)
	}

	return outputResult(cmd, resp)
}

func incidentsChatCmd(ctx context.Context, cmd *cli.Command) error {
	client, err := getClient(ctx)
	if err != nil {
		return handleError(cmd, err)
	}

	if cmd.NArg() < 1 {
		return handleError(cmd, errors.New("incident ID is required"))
	}
	incidentID := cmd.Args().First()
	body := cmd.String("body")

	resp, err := client.Incidents.CreateIncidentChatMessage(ctx, incidentID, components.CreateIncidentChatMessage{
		Body: body,
	})
	if err != nil {
		return handleError(cmd, err)
	}

	return outputResult(cmd, resp)
}
