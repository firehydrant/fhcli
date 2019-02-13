package events

import (
	"strings"
	"time"

	apiclient "github.com/firehydrant/fhcli/pkg/api_client"

	"github.com/firehydrant/api-client-go/client/changes"
	"github.com/firehydrant/api-client-go/models"
	"github.com/go-openapi/strfmt"
)

type ChangeEvent struct {
	Environment   string
	Service       string
	RawIdentities string
	RawLabels     string
	Identities    map[string]string
	Labels        map[string]string
	StartsAt      time.Time
	EndsAt        time.Time
	Summary       string
}

func NewChangeEvent() ChangeEvent {
	return ChangeEvent{
		Identities: make(map[string]string),
		Labels:     make(map[string]string),
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

func (ce *ChangeEvent) parseLabels() {
	for _, label := range strings.Split(ce.RawLabels, ",") {
		lParts := strings.Split(label, "=")
		if len(lParts) > 1 {
			ce.Labels[lParts[0]] = lParts[1]
		}
	}
}

func (ce *ChangeEvent) Submit(client apiclient.ApiClient) (string, error) {
	c := changes.NewPostV1ChangesEventsParams()
	ce.parseIdentities()
	ce.parseLabels()

	envList := []string{}
	if len(ce.Environment) > 0 {
		envList = append(envList, ce.Environment)
	}

	svcList := []string{}
	if len(ce.Service) > 0 {
		svcList = append(svcList, ce.Service)
	}

	c.V1ChangesEvents = &models.PostV1ChangesEvents{
		Environments:     envList,
		Services:         svcList,
		StartsAt:         strfmt.DateTime(ce.StartsAt),
		EndsAt:           strfmt.DateTime(ce.EndsAt),
		Summary:          &ce.Summary,
		ChangeIdentities: ce.identitiesToAPI(),
		Labels:           ce.Labels,
	}

	resp, err := client.Client.Changes.PostV1ChangesEvents(c, client.Auth)

	if err != nil {
		return "", err
	}

	return resp.Payload.ID, nil
}
