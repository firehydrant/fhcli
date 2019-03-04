package fhclient

import (
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"

	client "github.com/firehydrant/api-client-go/client"
	transport "github.com/go-openapi/runtime/client"
)

type Config struct {
	ApiHost string
	ApiKey  string
	Debug   bool
}

type ApiClient struct {
	transport *transport.Runtime
	Client    *client.FireHydrant
	Auth      runtime.ClientAuthInfoWriter
}

func NewApiClient(c Config) ApiClient {
	fhApiClient := ApiClient{}

	fhApiClient.transport = transport.New(c.ApiHost, "", []string{"https"})
	fhApiClient.transport.Debug = c.Debug

	fhApiClient.Client = client.New(fhApiClient.transport, strfmt.Default)
	fhApiClient.Auth = transport.BearerToken(c.ApiKey)

	return fhApiClient
}
