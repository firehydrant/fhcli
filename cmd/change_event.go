package cmd

import (
	"strings"
	"time"

	"github.com/firehydrant/api-client-go/client/changes"
	"github.com/firehydrant/api-client-go/models"
	"github.com/go-openapi/strfmt"
)

type ChangeEvent struct {
	Environment   string
	Service       string
	RawIdentities string
	Identities    map[string]string
	StartsAt      time.Time
	EndsAt        time.Time
	Summary       string
}

func NewChangeEvent() ChangeEvent {
	return ChangeEvent{
		Identities: make(map[string]string),
		StartsAt:   time.Now(),
		EndsAt:     time.Now(),
	}
}

func (ce *ChangeEvent) identitiesToAPI() []*models.PostV1ChangesEventsChangeIdentitiesItems0 {
	ciItems := []*models.PostV1ChangesEventsChangeIdentitiesItems0{}

	for key, value := range ce.Identities {
		k := string(key)
		v := string(value)
		ci := models.PostV1ChangesEventsChangeIdentitiesItems0{Type: &k, Value: &v}
		ciItems = append(ciItems, &ci)
	}

	return ciItems
}

func (ce *ChangeEvent) parseIdentities() {
	for _, identity := range strings.Split(ce.RawIdentities, ",") {
		iParts := strings.Split(identity, "=")
		if len(iParts) > 1 {
			ce.Identities[iParts[0]] = iParts[1]
		}
	}
}

func (ce *ChangeEvent) Submit() (string, error) {
	c := changes.NewPostV1ChangesEventsParams()
	ce.parseIdentities()

	c.V1ChangesEvents = &models.PostV1ChangesEvents{
		Environments:     []string{ce.Environment},
		Services:         []string{ce.Service},
		StartsAt:         strfmt.DateTime(ce.StartsAt),
		EndsAt:           strfmt.DateTime(ce.EndsAt),
		Summary:          &ce.Summary,
		ChangeIdentities: ce.identitiesToAPI(),
	}

	resp, err := client.client.Changes.PostV1ChangesEvents(c, client.auth)

	if err != nil {
		return "", err
	}

	return resp.Payload.ID, nil
}
