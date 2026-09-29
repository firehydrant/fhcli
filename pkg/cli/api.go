package cli

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"unicode"

	"github.com/firehydrant/fhcli/pkg/spec"
	"github.com/urfave/cli/v3"
)

// flagBinding ties a generated CLI flag to the spec parameter it sets.
type flagBinding struct {
	flag   string     // CLI flag name, e.g. "tag-list"
	param  spec.Param // wire name and type, e.g. tag_list / array
	inBody bool       // body field (true) or query param (false)
}

// operationFunc executes a generated operation command.
type operationFunc func(ctx context.Context, cmd *cli.Command, op spec.Operation, bindings []flagBinding) error

// newAPICmd builds `fhcli api <tag> <operation>` from the spec. If the spec
// failed to load, the command reports that error instead, so a bad spec never
// breaks the other commands.
func newAPICmd(sp *spec.Spec, loadErr error, run operationFunc) *cli.Command {
	cmd := &cli.Command{
		Name:  "api",
		Usage: "call any FireHydrant API operation (generated from the OpenAPI spec)",
		Description: "Operations are grouped by API area: fhcli api <area> <operation> [args] [flags].\n" +
			"Run `fhcli api <area> --help` to list operations. Every operation supports --dry-run;\n" +
			"DELETE operations require --yes.",
	}

	if loadErr != nil {
		cmd.Action = func(ctx context.Context, c *cli.Command) error {
			return failCmd(c, loadErr)
		}
		return cmd
	}

	reserved := reservedFlagNames()
	groups := map[string]*cli.Command{}
	for _, op := range sp.Operations {
		name := commandName(op.Tag)
		group, ok := groups[name]
		if !ok {
			group = &cli.Command{Name: name, Usage: op.Tag + " operations"}
			groups[name] = group
			cmd.Commands = append(cmd.Commands, group)
		}
		group.Commands = append(group.Commands, newOperationCmd(op, reserved, run))
	}

	return cmd
}

func newOperationCmd(op spec.Operation, reserved map[string]bool, run operationFunc) *cli.Command {
	flags, bindings, dataOnly := operationFlags(op, reserved)

	args := make([]string, len(op.PathParams))
	for i, p := range op.PathParams {
		args[i] = "<" + p.Name + ">"
	}

	var help []string
	if op.Description != "" {
		help = append(help, op.Description)
	}
	help = append(help, op.Method+" "+op.Path)
	if len(dataOnly) > 0 {
		fields := make([]string, len(dataOnly))
		for i, p := range dataOnly {
			fields[i] = describeField(p)
		}
		help = append(help, "Set these body fields with --data: "+strings.Join(fields, ", "))
	}

	return &cli.Command{
		Name:        commandName(op.ID),
		Usage:       op.Summary,
		Description: strings.Join(help, "\n\n"),
		ArgsUsage:   strings.Join(args, " "),
		Flags:       flags,
		Action: func(ctx context.Context, cmd *cli.Command) error {
			return run(ctx, cmd, op, bindings)
		},
	}
}

// operationFlags returns the generated flags for an operation, the bindings
// that map them back to parameters, and the body fields that can only be set
// via --data (listed in help so agents can discover them). Query params win
// over body fields; names that would shadow a reserved flag are skipped.
func operationFlags(op spec.Operation, reserved map[string]bool) ([]cli.Flag, []flagBinding, []spec.Param) {
	used := map[string]bool{}
	var flags []cli.Flag
	var bindings []flagBinding
	var dataOnly []spec.Param
	listed := map[string]bool{}

	addDataOnly := func(p spec.Param) {
		if p.Name == "" || listed[p.Name] {
			return
		}
		listed[p.Name] = true
		dataOnly = append(dataOnly, p)
	}

	add := func(p spec.Param, inBody bool) {
		// Bracketed names (fields[label]) are Grape nested-hash notation:
		// meaningless as a JSON body key, so list the root key instead.
		if inBody && strings.ContainsAny(p.Name, "[]") {
			addDataOnly(spec.Param{Name: p.Name[:strings.IndexAny(p.Name, "[]")], Type: "object"})
			return
		}
		name := commandName(p.Name)
		if name == "" || reserved[name] || used[name] || !p.Scalar() {
			if inBody {
				addDataOnly(p)
			}
			return
		}
		used[name] = true
		flags = append(flags, newParamFlag(name, p))
		bindings = append(bindings, flagBinding{flag: name, param: p, inBody: inBody})
	}

	for _, p := range op.QueryParams {
		add(p, false)
	}
	for _, p := range op.BodyFields {
		add(p, true)
	}

	flags = append(flags,
		&cli.StringFlag{Name: "data", Usage: "request body as JSON: inline, @file.json, or - for stdin"},
		&cli.BoolFlag{Name: "dry-run", Usage: "print the request instead of sending it"},
		&cli.BoolFlag{Name: "yes", Usage: "confirm a destructive (DELETE) operation"},
	)

	return flags, bindings, dataOnly
}

// describeField renders a body field for help, e.g. "impacts (array of object)".
func describeField(p spec.Param) string {
	kind := p.Type
	if p.Type == "array" && p.ItemType != "" {
		kind = "array of " + p.ItemType
	}
	if p.Required {
		kind += ", required"
	}
	return p.Name + " (" + kind + ")"
}

func newParamFlag(name string, p spec.Param) cli.Flag {
	usage := p.Description
	if len(p.Enum) > 0 {
		usage += " (one of: " + strings.Join(p.Enum, ", ") + ")"
	}
	if p.Required {
		usage += " (required)"
	}
	// urfave/cli treats `text` in a usage string as the value placeholder.
	usage = strings.TrimSpace(strings.ReplaceAll(usage, "`", "'"))

	switch p.Type {
	case "integer":
		return &cli.IntFlag{Name: name, Usage: usage}
	case "number":
		return &cli.FloatFlag{Name: name, Usage: usage}
	case "boolean":
		return &cli.BoolFlag{Name: name, Usage: usage}
	case "array":
		return &cli.StringSliceFlag{Name: name, Usage: usage}
	default:
		return &cli.StringFlag{Name: name, Usage: usage}
	}
}

// reservedFlagNames are names generated flags may not use: root flags are
// persistent in urfave/cli v3 (they apply to every subcommand), and data,
// dry-run, and yes are added to every generated command.
func reservedFlagNames() map[string]bool {
	reserved := map[string]bool{
		"help": true, "h": true, "version": true, "v": true,
		"data": true, "dry-run": true, "yes": true,
	}
	for _, f := range globalFlags {
		for _, n := range f.Names() {
			reserved[n] = true
		}
	}
	return reserved
}

// commandName turns a tag, operationId, or param name into a kebab-case CLI
// name: "Incident Settings" → "incident-settings", "list_incidents" →
// "list-incidents", "startIndex" → "start-index", "SCIM" → "scim".
func commandName(s string) string {
	var b strings.Builder
	pendingDash := false
	prevLowerOrDigit := false

	for _, r := range s {
		switch {
		case unicode.IsUpper(r):
			if prevLowerOrDigit {
				pendingDash = true
			}
			if pendingDash && b.Len() > 0 {
				b.WriteByte('-')
			}
			b.WriteRune(unicode.ToLower(r))
			pendingDash, prevLowerOrDigit = false, false
		case unicode.IsLetter(r) || unicode.IsDigit(r):
			if pendingDash && b.Len() > 0 {
				b.WriteByte('-')
			}
			b.WriteRune(r)
			pendingDash, prevLowerOrDigit = false, true
		default:
			pendingDash, prevLowerOrDigit = true, false
		}
	}
	return b.String()
}

// stdinReader is where `--data -` reads from; tests replace it.
var stdinReader io.Reader = os.Stdin

// apiRequest is a fully resolved request for one operation.
type apiRequest struct {
	Method string
	Path   string // path with args substituted, e.g. /v1/incidents/abc
	Query  url.Values
	Body   map[string]any // nil when no body is sent
}

type dryRunOutput struct {
	Method string         `json:"method"`
	URL    string         `json:"url"`
	Query  url.Values     `json:"query,omitempty"`
	Body   map[string]any `json:"body,omitempty"`
}

// runOperation builds the request for op and either prints it (--dry-run)
// or sends it. DELETE operations require --yes.
func runOperation(ctx context.Context, cmd *cli.Command, op spec.Operation, bindings []flagBinding) error {
	req, err := buildRequest(cmd, op, bindings)
	if err != nil {
		return failCmd(cmd, err)
	}

	if cmd.Bool("dry-run") {
		return outputJSON(os.Stdout, dryRunOutput{
			Method: req.Method,
			URL:    resolveAPIHost(cmd) + req.Path,
			Query:  req.Query,
			Body:   req.Body,
		})
	}

	if req.Method == http.MethodDelete && !cmd.Bool("yes") {
		return failCmd(cmd, fmt.Errorf("%s is destructive; re-run with --yes", op.ID))
	}

	client, err := getRawClient(ctx)
	if err != nil {
		return failCmd(cmd, err)
	}

	var body []byte
	if req.Body != nil {
		if body, err = json.Marshal(req.Body); err != nil {
			return failCmd(cmd, err)
		}
	}

	resp, err := client.do(ctx, req.Method, req.Path, req.Query, body)
	if err != nil {
		return failCmd(cmd, err)
	}
	return outputResult(cmd, resp)
}

// buildRequest substitutes path args, collects set query flags, and merges
// --data with set body flags (flags win).
func buildRequest(cmd *cli.Command, op spec.Operation, bindings []flagBinding) (*apiRequest, error) {
	args := cmd.Args().Slice()
	if len(args) < len(op.PathParams) {
		return nil, fmt.Errorf("%s is required", op.PathParams[len(args)].Name)
	}
	if len(args) > len(op.PathParams) {
		return nil, fmt.Errorf("unexpected argument %q", args[len(op.PathParams)])
	}

	req := &apiRequest{Method: op.Method, Path: op.Path, Query: url.Values{}}
	for i, p := range op.PathParams {
		req.Path = strings.ReplaceAll(req.Path, "{"+p.Name+"}", url.PathEscape(args[i]))
	}

	body, err := readData(cmd.String("data"))
	if err != nil {
		return nil, err
	}
	if body != nil && !op.HasBody {
		return nil, fmt.Errorf("%s takes no request body; remove --data", op.ID)
	}

	for _, b := range bindings {
		if !cmd.IsSet(b.flag) {
			continue
		}
		value := flagValue(cmd, b)
		if b.inBody {
			if body == nil {
				body = map[string]any{}
			}
			body[b.param.Name] = value
			continue
		}
		if values, ok := value.([]string); ok {
			key := arrayQueryKey(b.param.Name)
			for _, v := range values {
				req.Query.Add(key, v)
			}
		} else {
			req.Query.Set(b.param.Name, fmt.Sprint(value))
		}
	}

	req.Body = body
	return req, nil
}

// arrayQueryKey returns the Rack array form of a query param name, which is
// what Grape parses as an Array: "include" → "include[]",
// "custom_fields[field_id]" → "custom_fields[][field_id]". A repeated bare
// key keeps only the last value, and parent[child] parses as a Hash.
func arrayQueryKey(name string) string {
	if i := strings.Index(name, "["); i >= 0 {
		return name[:i] + "[]" + name[i:]
	}
	return name + "[]"
}

func flagValue(cmd *cli.Command, b flagBinding) any {
	switch b.param.Type {
	case "integer":
		return cmd.Int(b.flag)
	case "number":
		return cmd.Float(b.flag)
	case "boolean":
		return cmd.Bool(b.flag)
	case "array":
		return cmd.StringSlice(b.flag)
	default:
		return cmd.String(b.flag)
	}
}

// readData parses --data: inline JSON, @path to a file, or - for stdin. The
// value must be a JSON object; "" means no --data.
func readData(value string) (map[string]any, error) {
	if value == "" {
		return nil, nil
	}

	var raw []byte
	var err error
	switch {
	case value == "-":
		raw, err = io.ReadAll(stdinReader)
	case strings.HasPrefix(value, "@"):
		raw, err = os.ReadFile(strings.TrimPrefix(value, "@"))
	default:
		raw = []byte(value)
	}
	if err != nil {
		return nil, fmt.Errorf("read --data: %w", err)
	}

	var body map[string]any
	if err := json.Unmarshal(raw, &body); err != nil || body == nil {
		return nil, errors.New("--data must be a JSON object")
	}
	return body, nil
}
