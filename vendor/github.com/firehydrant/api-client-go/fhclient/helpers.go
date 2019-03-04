package fhclient

import (
	"strings"

	"github.com/firehydrant/api-client-go/models"
)

type ApiKV struct {
	Type  *string
	Value *string
}

func MapToAPIKV(identities map[string]string) []ApiKV {
	items := []ApiKV{}

	for key, value := range identities {
		k := string(key)
		v := string(value)
		ci := ApiKV{Type: &k, Value: &v}
		items = append(items, ci)
	}

	return items
}

func ParseKV(raw string) map[string]string {
	t := make(map[string]string)
	for _, identity := range strings.Split(raw, ",") {
		iParts := strings.Split(identity, "=")
		if len(iParts) > 1 {
			t[iParts[0]] = iParts[1]
		}
	}
	return t
}

func ParamToList(raw string) []string {
	list := []string{}
	if len(raw) > 0 {
		list = append(list, raw)
	}
	return list
}

func APIKVtoChangeIdentities(akv []ApiKV) []*models.PostV1ChangesEventsChangeIdentitiesItems0 {
	cei := []*models.PostV1ChangesEventsChangeIdentitiesItems0{}

	for _, i := range akv {
		ci := models.PostV1ChangesEventsChangeIdentitiesItems0(i)
		cei = append(cei, &ci)
	}

	return cei
}
