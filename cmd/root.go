package cmd

import (
	"fmt"
	"os"

	fhclient "github.com/firehydrant/api-client-go/client"
	runtime "github.com/go-openapi/runtime"
	transport "github.com/go-openapi/runtime/client"

	"github.com/go-openapi/strfmt"
	"github.com/spf13/cobra"
)

type apiClient struct {
	transport *transport.Runtime
	client    *fhclient.FireHydrant
	auth      runtime.ClientAuthInfoWriter
}

var apiKey string
var apiHost string
var debug bool
var client apiClient

func init() {
	rootCmd.PersistentFlags().StringVar(&apiKey, "api-key", "", "API key")
	rootCmd.PersistentFlags().StringVar(&apiHost, "api-host", "", "API hostname (default: api.firehydrant.io")
	rootCmd.PersistentFlags().BoolVarP(&debug, "debug", "d", false, "Debug HTTP requests")
}

var rootCmd = &cobra.Command{
	Use:   "fhcli",
	Short: "fhcli is the command line client for the FireHydrant API",
	Long: `Submit change events and other metadata to the FireHydrant API.

API Documentation: https://apidocs.firehydrant.io`,
	Run: func(cmd *cobra.Command, args []string) {},
}

// Execute the base command
func Execute() {

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func setupClient() error {
	client.transport = transport.New(apiHost, "", []string{"https"})
	client.transport.Debug = debug

	client.client = fhclient.New(client.transport, strfmt.Default)
	client.auth = transport.BearerToken(apiKey)

	return nil
}
