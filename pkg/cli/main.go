package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path"
	"strings"

	"github.com/firehydrant/fhcli/pkg/spec"
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

// newGlobalFlags returns fresh root flags. Each root command needs its own:
// urfave/cli keeps parsed state on flag values, so sharing them across roots
// (as tests do) would stop env var sources from being re-read.
func newGlobalFlags() []cli.Flag {
	return []cli.Flag{
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
			Value:   defaultAPIHost,
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
}

// NewRootCmd builds the CLI with the API spec from spec.Load.
func NewRootCmd(commit, version string) *cli.Command {
	sp, err := spec.Load()
	return newRootCmd(commit, version, sp, err)
}

func newRootCmd(commit, version string, sp *spec.Spec, specErr error) *cli.Command {
	cmd := &cli.Command{
		Name:    "fhcli",
		Usage:   "FireHydrant CLI",
		Version: fmt.Sprintf("%s (%s)", version, commit),
		Flags:   newGlobalFlags(),
		Before:  rootBefore,
		Commands: []*cli.Command{
			newInitCmd(),
			newEventCmd(),
			newExecuteCmd(),
			newIncidentsCmd(),
			newServicesCmd(),
			newTeamsCmd(),
			newAPICmd(sp, specErr, runOperation),
			newSpecCmd(),
			newLoginCmd(),
			newLogoutCmd(),
			newWhoamiCmd(),
			newInstallCmd(),
		},
	}

	return cmd
}

func rootBefore(ctx context.Context, cmd *cli.Command) (context.Context, error) {
	// Load config file values
	loadConfigFile(cmd)
	sawUnauthorized.Store(false)

	// Skip auth for commands that don't call the API with it
	cmdName := ""
	if cmd.Args().Present() {
		cmdName = cmd.Args().First()
	}
	switch cmdName {
	case "init", "install", "spec", "login", "logout":
		return ctx, nil
	}

	info, err := resolveAuth(ctx, cmd)
	if err != nil {
		// Reported by the command via getClient, so --ignore-errors still
		// applies to event/execute.
		return context.WithValue(ctx, authErrKey, err), nil
	}
	if info == nil {
		return ctx, nil
	}

	opts := []firehydrantgosdk.SDKOption{
		firehydrantgosdk.WithSecurity(components.Security{APIKey: info.header}),
		firehydrantgosdk.WithClient(apiHTTPClient),
	}
	if info.apiHost != defaultAPIHost {
		opts = append(opts, firehydrantgosdk.WithServerURL(info.apiHost))
	}

	ctx = context.WithValue(ctx, authInfoKey, info)
	ctx = context.WithValue(ctx, sdkClientKey, firehydrantgosdk.New(opts...))
	ctx = context.WithValue(ctx, rawClientKey, &rawClient{baseURL: info.apiHost, token: info.header, http: apiHTTPClient})
	return ctx, nil
}

// resolveAPIHost returns --api-host with a scheme and no trailing slash,
// defaulting to production.
func resolveAPIHost(cmd *cli.Command) string {
	host := strings.TrimRight(cmd.String("api-host"), "/")
	if host == "" {
		return defaultAPIHost
	}
	if !strings.HasPrefix(host, "http://") && !strings.HasPrefix(host, "https://") {
		host = "https://" + host
	}
	return host
}

func getClient(ctx context.Context) (*firehydrantgosdk.FireHydrant, error) {
	if err, ok := ctx.Value(authErrKey).(error); ok {
		return nil, err
	}
	client, ok := ctx.Value(sdkClientKey).(*firehydrantgosdk.FireHydrant)
	if !ok || client == nil {
		return nil, errNoCredentials
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
		// Flags are already parsed, so set root flags (token, api-host,
		// debug) directly; cmd.Set fails harmlessly for subcommand flags.
		// The env var export covers subcommand flags, which are parsed later.
		switch v := val.(type) {
		case string:
			if v != "" {
				_ = cmd.Set(flagName, v)
				os.Setenv(envForFlag(flagName), v)
			}
		case bool:
			if v {
				_ = cmd.Set(flagName, "true")
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

// outputJSON writes data as JSON: indented when w is a terminal, one compact
// line otherwise, so agents and pipes get terse output that costs fewer tokens.
func outputJSON(w io.Writer, data any) error {
	enc := json.NewEncoder(w)
	if isTerminal(w) {
		enc.SetIndent("", "  ")
	}
	return enc.Encode(data)
}

// isTerminal reports whether w is a character device (a terminal), as
// opposed to a pipe, file, or buffer.
func isTerminal(w io.Writer) bool {
	f, ok := w.(*os.File)
	if !ok {
		return false
	}
	info, err := f.Stat()
	return err == nil && info.Mode()&os.ModeCharDevice != 0
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

// handleError reports an error and respects --ignore-errors. It is for the
// change-event commands (event, execute, init), which default to exiting 0 so
// a FireHydrant outage never fails a CI pipeline.
func handleError(cmd *cli.Command, err error) error {
	reportError(cmd, err)

	if cmd.Bool("ignore-errors") {
		return nil
	}

	return err
}

// failCmd reports an error and always returns it so the CLI exits non-zero.
// Incident, service, and team commands use it: agents and scripts rely on the
// exit code to tell a failed lookup or write from a successful one.
func failCmd(cmd *cli.Command, err error) error {
	reportError(cmd, err)
	return err
}

func reportError(cmd *cli.Command, err error) {
	err = withLoginHint(err)
	if cmd.Bool("interactive") {
		fmt.Fprintln(os.Stderr, "Error:", err.Error())
	} else {
		outputJSON(os.Stderr, map[string]string{"error": err.Error()})
	}
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
