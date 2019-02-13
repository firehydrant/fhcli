package api_client

import (
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"

	fhclient "github.com/firehydrant/api-client-go/client"
	transport "github.com/go-openapi/runtime/client"
)

type Config struct {
	ApiHost string
	ApiKey  string
	Debug   bool
}

type ApiClient struct {
	transport *transport.Runtime
	Client    *fhclient.FireHydrant
	Auth      runtime.ClientAuthInfoWriter
}

func NewApiClient(c Config) ApiClient {
	client := ApiClient{}

	client.transport = transport.New(c.ApiHost, "", []string{"http"})
	client.transport.Debug = c.Debug

	client.Client = fhclient.New(client.transport, strfmt.Default)
	client.Auth = transport.BearerToken(c.ApiKey)

	return client
}
