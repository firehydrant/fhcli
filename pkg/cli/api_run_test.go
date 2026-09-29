package cli

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

// runMini runs the CLI with the 7-operation fixture spec. An empty serverURL
// runs with no token and the default API host.
func runMini(t *testing.T, serverURL string, args ...string) (string, string, error) {
	t.Helper()
	argv := []string{"fhcli"}
	if serverURL != "" {
		argv = append(argv, "--token", "test-token", "--api-host", serverURL)
	}
	return runRoot(t, newRootCmd("test", "0.0.0-test", miniSpec(t), nil), append(argv, args...)...)
}

type dryRun struct {
	Method string              `json:"method"`
	URL    string              `json:"url"`
	Query  map[string][]string `json:"query"`
	Body   map[string]any      `json:"body"`
}

func parseDryRun(t *testing.T, out string) dryRun {
	t.Helper()
	var d dryRun
	if err := json.Unmarshal([]byte(out), &d); err != nil {
		t.Fatalf("dry-run output is not JSON: %v\n%s", err, out)
	}
	return d
}

func failIfCalled(t *testing.T) *httptest.Server {
	t.Helper()
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Errorf("no request should be sent, got %s %s", r.Method, r.URL.Path)
	}))
	t.Cleanup(ts.Close)
	return ts
}

func TestAPIDryRunPrintsRequestWithoutSending(t *testing.T) {
	ts := failIfCalled(t)

	out, _, err := runMini(t, ts.URL, "api", "incidents", "create-incident",
		"--name", "DB down", "--priority", "2", "--tag-list", "db", "--tag-list", "sev", "--dry-run")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	d := parseDryRun(t, out)
	if d.Method != "POST" || d.URL != ts.URL+"/v1/incidents" {
		t.Errorf("request = %s %s", d.Method, d.URL)
	}
	want := map[string]any{"name": "DB down", "priority": float64(2), "tag_list": []any{"db", "sev"}}
	if !reflect.DeepEqual(d.Body, want) {
		t.Errorf("body = %v, want %v", d.Body, want)
	}
}

func TestAPIDryRunIncludesQuery(t *testing.T) {
	ts := failIfCalled(t)

	out, _, err := runMini(t, ts.URL, "api", "incidents", "list-incidents",
		"--page", "2", "--status", "open", "--tags", "a", "--tags", "b", "--custom-fields-field-id", "cf1", "--dry-run")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	d := parseDryRun(t, out)
	want := map[string][]string{"page": {"2"}, "status": {"open"}, "tags": {"a", "b"}, "custom_fields[field_id]": {"cf1"}}
	if !reflect.DeepEqual(d.Query, want) {
		t.Errorf("query = %v, want %v", d.Query, want)
	}
	if d.Body != nil {
		t.Errorf("GET should have no body, got %v", d.Body)
	}
}

func TestAPIDryRunWorksWithoutToken(t *testing.T) {
	out, _, err := runMini(t, "", "api", "incidents", "list-incidents", "--dry-run")

	if err != nil {
		t.Fatalf("dry-run should not need a token: %v", err)
	}
	if d := parseDryRun(t, out); !strings.HasSuffix(d.URL, "/v1/incidents") {
		t.Errorf("url = %q", d.URL)
	}
}

func TestAPIAddsSchemeToAPIHost(t *testing.T) {
	root := newRootCmd("test", "0.0.0-test", miniSpec(t), nil)

	out, _, err := runRoot(t, root, "fhcli", "--api-host", "api.staging.firehydrant.io/", "api", "incidents", "list-incidents", "--dry-run")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if d := parseDryRun(t, out); d.URL != "https://api.staging.firehydrant.io/v1/incidents" {
		t.Errorf("url = %q", d.URL)
	}
}

func TestAPIEscapesPathArgs(t *testing.T) {
	ts := failIfCalled(t)

	out, _, err := runMini(t, ts.URL, "api", "incidents", "get-incident-role", "a/b", "r 1", "--dry-run")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if d := parseDryRun(t, out); d.URL != ts.URL+"/v1/incidents/a%2Fb/roles/r%201" {
		t.Errorf("url = %q", d.URL)
	}
}

func TestAPIRequiresPathArgs(t *testing.T) {
	ts := failIfCalled(t)

	_, stderr, err := runMini(t, ts.URL, "api", "incidents", "get-incident-role", "inc-1")

	if err == nil {
		t.Error("expected an error")
	}
	if !strings.Contains(stderr, "role_id is required") {
		t.Errorf("stderr = %q", stderr)
	}
}

func TestAPIRejectsExtraArgs(t *testing.T) {
	ts := failIfCalled(t)

	_, stderr, err := runMini(t, ts.URL, "api", "incidents", "list-incidents", "surprise")

	if err == nil {
		t.Error("expected an error")
	}
	if !strings.Contains(stderr, `unexpected argument \"surprise\"`) {
		t.Errorf("stderr = %q", stderr)
	}
}

func TestAPIDataMergesWithFlagOverrides(t *testing.T) {
	ts := failIfCalled(t)

	out, _, err := runMini(t, ts.URL, "api", "incidents", "create-incident",
		"--data", `{"name":"old","labels":{"team":"db"}}`, "--name", "new", "--dry-run")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := map[string]any{"name": "new", "labels": map[string]any{"team": "db"}}
	if d := parseDryRun(t, out); !reflect.DeepEqual(d.Body, want) {
		t.Errorf("body = %v, want %v", d.Body, want)
	}
}

func TestAPIDataFromFile(t *testing.T) {
	ts := failIfCalled(t)
	path := filepath.Join(t.TempDir(), "incident.json")
	if err := os.WriteFile(path, []byte(`{"name":"from file"}`), 0o644); err != nil {
		t.Fatal(err)
	}

	out, _, err := runMini(t, ts.URL, "api", "incidents", "create-incident", "--data", "@"+path, "--dry-run")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if d := parseDryRun(t, out); d.Body["name"] != "from file" {
		t.Errorf("body = %v", d.Body)
	}
}

func TestAPIDataFromStdin(t *testing.T) {
	ts := failIfCalled(t)
	old := stdinReader
	stdinReader = strings.NewReader(`{"name":"piped"}`)
	t.Cleanup(func() { stdinReader = old })

	out, _, err := runMini(t, ts.URL, "api", "incidents", "create-incident", "--data", "-", "--dry-run")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if d := parseDryRun(t, out); d.Body["name"] != "piped" {
		t.Errorf("body = %v", d.Body)
	}
}

func TestAPIDataMustBeJSONObject(t *testing.T) {
	ts := failIfCalled(t)
	old := stdinReader
	t.Cleanup(func() { stdinReader = old })

	for _, data := range []string{`[1]`, `null`, `not json`, `-`} {
		stdinReader = strings.NewReader("") // only read for "-": empty stdin

		_, stderr, err := runMini(t, ts.URL, "api", "incidents", "create-incident", "--data", data)

		if err == nil {
			t.Errorf("--data %s: expected an error", data)
		}
		if !strings.Contains(stderr, "--data must be a JSON object") {
			t.Errorf("--data %s: stderr = %q", data, stderr)
		}
	}
}

func TestAPIDataRejectedWhenOperationHasNoBody(t *testing.T) {
	ts := failIfCalled(t)

	_, stderr, err := runMini(t, ts.URL, "api", "incidents", "list-incidents", "--data", `{}`)

	if err == nil {
		t.Error("expected an error")
	}
	if !strings.Contains(stderr, "list_incidents takes no request body") {
		t.Errorf("stderr = %q", stderr)
	}
}

func TestAPISendsExplicitFalseBool(t *testing.T) {
	ts := failIfCalled(t)

	out, _, err := runMini(t, ts.URL, "api", "incidents", "close-incident", "inc-1", "--notify=false", "--dry-run")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if d := parseDryRun(t, out); d.Body["notify"] != false {
		t.Errorf("explicit --notify=false must be sent, body = %v", d.Body)
	}

	out, _, err = runMini(t, ts.URL, "api", "incidents", "close-incident", "inc-1", "--milestone", "resolved", "--dry-run")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, present := parseDryRun(t, out).Body["notify"]; present {
		t.Error("unset --notify must not be sent")
	}
}

func TestAPISendsRequest(t *testing.T) {
	var gotMethod, gotPath, gotAuth string
	var gotBody map[string]any
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotMethod, gotPath, gotAuth = r.Method, r.URL.Path, r.Header.Get("Authorization")
		json.NewDecoder(r.Body).Decode(&gotBody)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`{"id":"inc-new"}`))
	}))
	defer ts.Close()

	out, _, err := runMini(t, ts.URL, "api", "incidents", "create-incident", "--name", "DB down")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "inc-new") {
		t.Errorf("stdout = %q", out)
	}
	if gotMethod != "POST" || gotPath != "/v1/incidents" || gotAuth != "test-token" {
		t.Errorf("request = %s %s auth=%q", gotMethod, gotPath, gotAuth)
	}
	if gotBody["name"] != "DB down" {
		t.Errorf("sent body = %v", gotBody)
	}
}

func TestAPIDeleteRequiresYes(t *testing.T) {
	ts := failIfCalled(t)

	_, stderr, err := runMini(t, ts.URL, "api", "incidents", "delete-incident-role", "inc-1", "role-1")

	if err == nil {
		t.Error("expected an error")
	}
	if !strings.Contains(stderr, "delete_incident_role is destructive; re-run with --yes") {
		t.Errorf("stderr = %q", stderr)
	}
}

func TestAPIDryRunOfDeleteNeedsNoYes(t *testing.T) {
	ts := failIfCalled(t)

	out, _, err := runMini(t, ts.URL, "api", "incidents", "delete-incident-role", "inc-1", "role-1", "--dry-run")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if d := parseDryRun(t, out); d.Method != "DELETE" {
		t.Errorf("method = %q", d.Method)
	}
}

func TestAPIDeleteWithYesSends(t *testing.T) {
	var gotMethod, gotPath string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotMethod, gotPath = r.Method, r.URL.Path
		w.WriteHeader(http.StatusNoContent)
	}))
	defer ts.Close()

	out, _, err := runMini(t, ts.URL, "api", "incidents", "delete-incident-role", "inc-1", "role-1", "--yes")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotMethod != "DELETE" || gotPath != "/v1/incidents/inc-1/roles/role-1" {
		t.Errorf("request = %s %s", gotMethod, gotPath)
	}
	if strings.TrimSpace(out) != "{}" {
		t.Errorf("stdout = %q, want {}", out)
	}
}

func TestAPINon2xxExitsNonZero(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"detail":"incident not found"}`))
	}))
	defer ts.Close()

	_, stderr, err := runMini(t, ts.URL, "api", "incidents", "get-incident-role", "inc-1", "role-1")

	if err == nil {
		t.Error("expected an error so the CLI exits non-zero")
	}
	if !strings.Contains(stderr, "404") || !strings.Contains(stderr, "incident not found") {
		t.Errorf("stderr = %q", stderr)
	}
}

func TestAPIReportsInvalidSpecOverride(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", dir)
	os.MkdirAll(filepath.Join(dir, "fhcli"), 0o755)
	os.WriteFile(filepath.Join(dir, "fhcli", "spec.json"), []byte("garbage"), 0o644)

	_, stderr, err := runCmdFull(t, "http://127.0.0.1:1", "api", "incidents", "list-incidents")

	if err == nil {
		t.Error("expected an error")
	}
	if !strings.Contains(stderr, "fhcli spec update") {
		t.Errorf("stderr should explain how to recover, got %q", stderr)
	}
}

func TestEventWorksWithInvalidSpecOverride(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", dir)
	os.MkdirAll(filepath.Join(dir, "fhcli"), 0o755)
	os.WriteFile(filepath.Join(dir, "fhcli", "spec.json"), []byte("garbage"), 0o644)
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated) // CreateChangeEvent responds 201
		w.Write([]byte(`{"id":"evt-ok"}`))
	}))
	defer ts.Close()

	out, err := runCmd(t, ts.URL, "event", "deployed")

	if err != nil {
		t.Fatalf("event must not depend on the spec: %v", err)
	}
	if !strings.Contains(out, "evt-ok") {
		t.Errorf("stdout = %q", out)
	}
}
