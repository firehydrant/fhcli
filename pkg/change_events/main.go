package events

import (
	"os"
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
	ce := ChangeEvent{
		Identities: make(map[string]string),
		Labels:     make(map[string]string),
		StartsAt:   time.Now(),
		EndsAt:     time.Now(),
	}

	ce.Labels["hostname"] = os.Getenv("HOSTNAME")

	return ce
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

func (ce *ChangeEvent) parseKV(raw string, target map[string]string) {
	for _, identity := range strings.Split(raw, ",") {
		iParts := strings.Split(identity, "=")
		if len(iParts) > 1 {
			target[iParts[0]] = iParts[1]
		}
	}
}

func (ce *ChangeEvent) Submit(client apiclient.ApiClient) (string, error) {
	c := changes.NewPostV1ChangesEventsParams()
	ce.parseKV(ce.RawIdentities, ce.Identities)
	ce.parseKV(ce.RawLabels, ce.Labels)

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
