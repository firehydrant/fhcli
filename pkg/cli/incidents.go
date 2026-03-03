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
			{
				Name:   "alerts",
				Usage:  "list alerts attached to an incident",
				Action: incidentsAlertsCmd,
			},
			{
				Name:   "summary",
				Usage:  "get AI-generated audience summaries for an incident",
				Action: incidentsSummaryCmd,
			},
			{
				Name:   "events",
				Usage:  "list incident timeline events",
				Action: incidentsEventsCmd,
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
						Name:  "types",
						Usage: "comma-separated event types to filter (e.g. note,chat_message,incident_status)",
					},
				},
			},
			{
				Name:   "transcript",
				Usage:  "list incident call transcript entries",
				Action: incidentsTranscriptCmd,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "sort",
						Usage: "sort order: asc or desc",
						Value: "asc",
					},
				},
			},
			{
				Name:   "update",
				Usage:  "update incident fields",
				Action: incidentsUpdateCmd,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "severity",
						Usage: "incident severity (e.g. SEV1)",
					},
					&cli.StringFlag{
						Name:  "summary",
						Usage: "incident summary text",
					},
					&cli.StringFlag{
						Name:  "priority",
						Usage: "incident priority (e.g. P1)",
					},
					&cli.StringFlag{
						Name:  "description",
						Usage: "incident description",
					},
				},
			},
			{
				Name:   "resolve",
				Usage:  "resolve an incident",
				Action: incidentsResolveCmd,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "milestone",
						Usage: "post-incident milestone slug to set",
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

func incidentsAlertsCmd(ctx context.Context, cmd *cli.Command) error {
	client, err := getClient(ctx)
	if err != nil {
		return handleError(cmd, err)
	}

	if cmd.NArg() < 1 {
		return handleError(cmd, errors.New("incident ID is required"))
	}
	incidentID := cmd.Args().First()

	resp, err := client.Alerts.ListIncidentAlerts(ctx, incidentID)
	if err != nil {
		return handleError(cmd, err)
	}

	return outputResult(cmd, resp)
}

func incidentsSummaryCmd(ctx context.Context, cmd *cli.Command) error {
	client, err := getClient(ctx)
	if err != nil {
		return handleError(cmd, err)
	}

	if cmd.NArg() < 1 {
		return handleError(cmd, errors.New("incident ID is required"))
	}
	incidentID := cmd.Args().First()

	resp, err := client.Audiences.ListAudienceSummaries(ctx, incidentID)
	if err != nil {
		return handleError(cmd, err)
	}

	return outputResult(cmd, resp)
}

func incidentsEventsCmd(ctx context.Context, cmd *cli.Command) error {
	client, err := getClient(ctx)
	if err != nil {
		return handleError(cmd, err)
	}

	if cmd.NArg() < 1 {
		return handleError(cmd, errors.New("incident ID is required"))
	}
	incidentID := cmd.Args().First()

	page := int(cmd.Int("page"))
	perPage := int(cmd.Int("per-page"))

	var types *string
	if t := cmd.String("types"); t != "" {
		types = firehydrantgosdk.String(t)
	}

	resp, err := client.Incidents.ListIncidentEvents(ctx, incidentID, types, firehydrantgosdk.Int(page), firehydrantgosdk.Int(perPage))
	if err != nil {
		return handleError(cmd, err)
	}

	return outputResult(cmd, resp)
}

func incidentsTranscriptCmd(ctx context.Context, cmd *cli.Command) error {
	client, err := getClient(ctx)
	if err != nil {
		return handleError(cmd, err)
	}

	if cmd.NArg() < 1 {
		return handleError(cmd, errors.New("incident ID is required"))
	}
	incidentID := cmd.Args().First()

	var sort *operations.ListTranscriptEntriesSort
	if s := cmd.String("sort"); s != "" {
		v := operations.ListTranscriptEntriesSort(s)
		sort = &v
	}

	resp, err := client.Incidents.ListTranscriptEntries(ctx, incidentID, nil, nil, sort)
	if err != nil {
		return handleError(cmd, err)
	}

	return outputResult(cmd, resp)
}

func incidentsUpdateCmd(ctx context.Context, cmd *cli.Command) error {
	client, err := getClient(ctx)
	if err != nil {
		return handleError(cmd, err)
	}

	if cmd.NArg() < 1 {
		return handleError(cmd, errors.New("incident ID is required"))
	}
	incidentID := cmd.Args().First()

	update := components.UpdateIncident{}

	if s := cmd.String("severity"); s != "" {
		update.Severity = firehydrantgosdk.String(s)
	}
	if s := cmd.String("summary"); s != "" {
		update.Summary = firehydrantgosdk.String(s)
	}
	if s := cmd.String("priority"); s != "" {
		update.Priority = firehydrantgosdk.String(s)
	}
	if s := cmd.String("description"); s != "" {
		update.Description = firehydrantgosdk.String(s)
	}

	resp, err := client.Incidents.UpdateIncident(ctx, incidentID, update)
	if err != nil {
		return handleError(cmd, err)
	}

	return outputResult(cmd, resp)
}

func incidentsResolveCmd(ctx context.Context, cmd *cli.Command) error {
	client, err := getClient(ctx)
	if err != nil {
		return handleError(cmd, err)
	}

	if cmd.NArg() < 1 {
		return handleError(cmd, errors.New("incident ID is required"))
	}
	incidentID := cmd.Args().First()

	var reqBody *operations.ResolveIncidentRequestBody
	if m := cmd.String("milestone"); m != "" {
		reqBody = &operations.ResolveIncidentRequestBody{
			Milestone: firehydrantgosdk.String(m),
		}
	}

	resp, err := client.Incidents.ResolveIncident(ctx, incidentID, reqBody)
	if err != nil {
		return handleError(cmd, err)
	}

	return outputResult(cmd, resp)
}
