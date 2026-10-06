package cli

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
	"strings"
	"testing"
)

func TestOutputJSONIsOneLineWhenNotATerminal(t *testing.T) {
	var buf bytes.Buffer

	if err := outputJSON(&buf, map[string]any{"id": "inc-1", "labels": map[string]string{"team": "db"}}); err != nil {
		t.Fatal(err)
	}

	if got := buf.String(); got != `{"id":"inc-1","labels":{"team":"db"}}`+"\n" {
		t.Errorf("output = %q, want compact single-line JSON", got)
	}
}

func TestOutputJSONCompactsRawAPIResponses(t *testing.T) {
	var buf bytes.Buffer
	raw := json.RawMessage("{\n  \"data\": [\n    {\"id\": \"inc-1\"}\n  ]\n}")

	if err := outputJSON(&buf, raw); err != nil {
		t.Fatal(err)
	}

	if got := buf.String(); got != `{"data":[{"id":"inc-1"}]}`+"\n" {
		t.Errorf("output = %q, want the API body compacted", got)
	}
}

func TestIsTerminal(t *testing.T) {
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	defer r.Close()
	defer w.Close()
	devNull, err := os.Open(os.DevNull)
	if err != nil {
		t.Fatal(err)
	}
	defer devNull.Close()

	if isTerminal(w) {
		t.Error("a pipe is not a terminal")
	}
	if isTerminal(&bytes.Buffer{}) {
		t.Error("a buffer is not a terminal")
	}
	if !isTerminal(devNull) {
		t.Error("a character device should be treated as a terminal")
	}
}

func TestPipedCommandOutputIsCompact(t *testing.T) {
	ts := newTestServer(t, map[string]http.HandlerFunc{
		"GET /v1/incidents/": jsonHandler(200, map[string]any{"id": "inc-abc", "name": "DB down"}),
	})
	defer ts.Close()

	out, err := runCmd(t, ts.URL, "incidents", "get", "inc-abc")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if strings.Count(strings.TrimSpace(out), "\n") != 0 || !strings.Contains(out, `"id":"inc-abc"`) {
		t.Errorf("piped output should be one compact line, got:\n%s", out)
	}
}

func TestPipedErrorsAreCompact(t *testing.T) {
	ts := newTestServer(t, map[string]http.HandlerFunc{})
	defer ts.Close()

	_, stderr, _ := runCmdFull(t, ts.URL, "incidents", "get")

	if got := strings.TrimSpace(stderr); got != `{"error":"incident ID is required"}` {
		t.Errorf("stderr = %q, want one compact JSON line", got)
	}
}
