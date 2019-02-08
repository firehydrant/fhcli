package cmd

import (
	"strings"
	"time"

	"github.com/firehydrant/api-client-go/client/changes"
	"github.com/firehydrant/api-client-go/models"
	"github.com/go-openapi/strfmt"
)

type ChangeEvent struct {
	Environment string
	Service     string
	Identities  string
	StartsAt    time.Time
	EndsAt      time.Time
	Summary     string
}

func NewChangeEvent() ChangeEvent {
	return ChangeEvent{
		StartsAt: time.Now(),
		EndsAt:   time.Now(),
	}
}

func (ce *ChangeEvent) parseIdentities() []*models.PostV1ChangesEventsChangeIdentitiesItems0 {
	ciItems := []*models.PostV1ChangesEventsChangeIdentitiesItems0{}

	for _, identity := range strings.Split(ce.Identities, ",") {
		iParts := strings.Split(identity, "=")
		if len(iParts) > 1 {
			ci := models.PostV1ChangesEventsChangeIdentitiesItems0{Type: &iParts[0], Value: &iParts[1]}
			ciItems = append(ciItems, &ci)
		}
	}

	return ciItems
}

func (ce *ChangeEvent) Submit() (string, error) {
	c := changes.NewPostV1ChangesEventsParams()

	c.V1ChangesEvents = &models.PostV1ChangesEvents{
		Environments:     []string{ce.Environment},
		Services:         []string{ce.Service},
		StartsAt:         strfmt.DateTime(ce.StartsAt),
		EndsAt:           strfmt.DateTime(ce.EndsAt),
		Summary:          &ce.Summary,
		ChangeIdentities: ce.parseIdentities(),
	}

	resp, err := client.client.Changes.PostV1ChangesEvents(c, client.auth)

	if err != nil {
		return "", err
	}

	return resp.Payload.ID, nil
}
