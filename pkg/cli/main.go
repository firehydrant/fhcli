package cli

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"strings"

	firehydrantgosdk "github.com/firehydrant/firehydrant-go-sdk"
	"github.com/firehydrant/firehydrant-go-sdk/models/components"
	"github.com/urfave/cli/v3"
	"gopkg.in/yaml.v2"
)

type ctxKey string

const sdkClientKey ctxKey = "sdkClient"

var sharedFlags = []cli.Flag{
	&cli.StringFlag{
		Name:    "identities",
		Usage:   "identities for the event, comma separated. example: image=quay.io/firehydrant/rails:abcdef0",
		Sources: cli.EnvVars("FH_IDENTITIES"),
	},
	&cli.StringFlag{
		Name:    "labels",
		Usage:   "labels for the event, comma separated. example: executor=jenkins",
		Sources: cli.EnvVars("FH_LABELS"),
	},
	&cli.StringFlag{
		Name:    "environment",
		Usage:   "environment for the event, example: production",
		Sources: cli.EnvVars("FH_ENVIRONMENT"),
	},
	&cli.StringFlag{
		Name:    "service",
		Usage:   "service for the event, example: rails-monolith",
		Sources: cli.EnvVars("FH_SERVICE"),
	},
}

var globalFlags = []cli.Flag{
	&cli.StringFlag{
		Name:    "token",
		Aliases: []string{"apiKey"},
		Usage:   "FireHydrant API token",
		Sources: cli.EnvVars("FIREHYDRANT_API_TOKEN", "FH_API_KEY"),
	},
	&cli.StringFlag{
		Name:    "api-host",
		Aliases: []string{"apiHost"},
		Usage:   "FireHydrant API hostname",
		Value:   "https://api.firehydrant.io",
		Sources: cli.EnvVars("FH_API_HOST"),
	},
	&cli.BoolFlag{
		Name:  "debug",
		Usage: "enable debug output for API calls",
	},
	&cli.BoolFlag{
		Name:    "interactive",
		Aliases: []string{"i"},
		Usage:   "enable human-readable output instead of JSON",
	},
	&cli.BoolFlag{
		Name:    "ignore-errors",
		Aliases: []string{"ignoreErrors"},
		Usage:   "exit 0 on errors from FH API",
		Value:   true,
	},
	&cli.StringFlag{
		Name:    "config",
		Usage:   "config file path",
		Sources: cli.EnvVars("FH_CONFIG_FILE"),
	},
}

func NewRootCmd(commit, version string) *cli.Command {
	cmd := &cli.Command{
		Name:    "fhcli",
		Usage:   "FireHydrant CLI",
		Version: fmt.Sprintf("%s (%s)", version, commit),
		Flags:   globalFlags,
		Before:  rootBefore,
		Commands: []*cli.Command{
			newInitCmd(),
			newEventCmd(),
			newExecuteCmd(),
			newIncidentsCmd(),
			newServicesCmd(),
			newTeamsCmd(),
			newInstallCmd(),
		},
	}

	return cmd
}

func rootBefore(ctx context.Context, cmd *cli.Command) (context.Context, error) {
	// Load config file values
	loadConfigFile(cmd)

	// Skip SDK client creation for commands that don't need it
	cmdName := ""
	if cmd.Args().Present() {
		cmdName = cmd.Args().First()
	}
	if cmdName == "init" || cmdName == "install" {
		return ctx, nil
	}

	token := cmd.String("token")
	if token == "" {
		return ctx, nil
	}

	opts := []firehydrantgosdk.SDKOption{
		firehydrantgosdk.WithSecurity(components.Security{
			APIKey: token,
		}),
	}

	apiHost := cmd.String("api-host")
	if apiHost != "" && apiHost != "https://api.firehydrant.io" {
		// Ensure the host has a scheme
		if !strings.HasPrefix(apiHost, "http://") && !strings.HasPrefix(apiHost, "https://") {
			apiHost = "https://" + apiHost
		}
		opts = append(opts, firehydrantgosdk.WithServerURL(apiHost))
	}

	client := firehydrantgosdk.New(opts...)
	ctx = context.WithValue(ctx, sdkClientKey, client)

	return ctx, nil
}

func getClient(ctx context.Context) (*firehydrantgosdk.FireHydrant, error) {
	client, ok := ctx.Value(sdkClientKey).(*firehydrantgosdk.FireHydrant)
	if !ok || client == nil {
		return nil, errors.New("API token is required. Set --token flag, FIREHYDRANT_API_TOKEN env var, or configure a config file")
	}
	return client, nil
}

// loadConfigFile reads the YAML config file and applies values for flags that are at their zero value.
func loadConfigFile(cmd *cli.Command) {
	paths := []string{
		cmd.String("config"),
		path.Join("/etc", "firehydrant.cfg"),
		path.Join(os.Getenv("HOME"), "firehydrant.cfg"),
		path.Join("/tmp", "firehydrant.cfg"),
	}

	var cfgPath string
	for _, p := range paths {
		if p == "" {
			continue
		}
		if _, err := os.Stat(p); err == nil {
			cfgPath = p
			break
		}
	}

	if cfgPath == "" {
		return
	}

	data, err := os.ReadFile(cfgPath)
	if err != nil {
		return
	}

	var cfg map[string]interface{}
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return
	}

	// Map config file keys to flag names
	keyMap := map[string]string{
		"apiKey":      "token",
		"apiHost":     "api-host",
		"debug":       "debug",
		"identities":  "identities",
		"labels":      "labels",
		"environment": "environment",
		"service":     "service",
	}

	for cfgKey, flagName := range keyMap {
		val, ok := cfg[cfgKey]
		if !ok || val == nil {
			continue
		}
		// Only apply config values if the flag hasn't been set via CLI or env
		if cmd.IsSet(flagName) {
			continue
		}
		// We can't easily set flag values after parsing in v3,
		// so we set environment variables as a fallback mechanism
		switch v := val.(type) {
		case string:
			if v != "" {
				os.Setenv(envForFlag(flagName), v)
			}
		case bool:
			if v {
				os.Setenv(envForFlag(flagName), "true")
			}
		}
	}
}

// envForFlag returns the primary environment variable name for a flag.
func envForFlag(flagName string) string {
	switch flagName {
	case "token":
		return "FH_API_KEY"
	case "api-host":
		return "FH_API_HOST"
	case "identities":
		return "FH_IDENTITIES"
	case "labels":
		return "FH_LABELS"
	case "environment":
		return "FH_ENVIRONMENT"
	case "service":
		return "FH_SERVICE"
	default:
		return ""
	}
}

// outputJSON writes JSON-encoded data to the given writer.
func outputJSON(w io.Writer, data any) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(data)
}

// outputResult writes data as JSON to stdout, or human-readable if --interactive is set.
func outputResult(cmd *cli.Command, data any) error {
	if cmd.Bool("interactive") {
		return outputInteractive(os.Stdout, data)
	}
	return outputJSON(os.Stdout, data)
}

// outputInteractive writes a simple human-readable representation.
func outputInteractive(w io.Writer, data any) error {
	// Marshal to JSON then unmarshal to a generic structure for pretty printing
	b, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}
	_, err = fmt.Fprintln(w, string(b))
	return err
}

// handleError formats an error appropriately and respects --ignore-errors.
func handleError(cmd *cli.Command, err error) error {
	if cmd.Bool("interactive") {
		fmt.Fprintln(os.Stderr, "Error:", err.Error())
	} else {
		outputJSON(os.Stderr, map[string]string{"error": err.Error()})
	}

	if cmd.Bool("ignore-errors") {
		return nil
	}

	return err
}

// parseKV parses a "k=v,k2=v2" string into a map.
func parseKV(s string) map[string]string {
	result := make(map[string]string)
	if s == "" {
		return result
	}
	pairs := strings.Split(s, ",")
	for _, pair := range pairs {
		kv := strings.SplitN(pair, "=", 2)
		if len(kv) == 2 {
			result[strings.TrimSpace(kv[0])] = strings.TrimSpace(kv[1])
		}
	}
	return result
}

// paramToList splits a comma-separated string into a list of trimmed strings.
func paramToList(s string) []string {
	if s == "" {
		return nil
	}
	parts := strings.Split(s, ",")
	result := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			result = append(result, p)
		}
	}
	return result
}
