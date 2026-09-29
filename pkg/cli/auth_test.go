package cli

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/firehydrant/fhcli/pkg/auth"
)

// storeLogin saves c as the "us" login in a fresh temp XDG_CONFIG_HOME.
func storeLogin(t *testing.T, c *auth.Credentials) {
	t.Helper()
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	if err := auth.DefaultStore().Save(c); err != nil {
		t.Fatal(err)
	}
}

func loginCreds(host string, expiry time.Time) *auth.Credentials {
	return &auth.Credentials{
		AccessToken: "login-at", RefreshToken: "login-rt", Expiry: expiry,
		Region: "us", AppHost: host, APIHost: host, ClientID: "cid",
		User: "Dylan", Email: "dylan@example.com", Organization: "FireHydrant", OrganizationID: "org-1",
	}
}

// authServer records the Authorization header of each API request and
// serves an incident, a token refresh, and a 401 route.
type authServer struct {
	*httptest.Server
	headers        []string
	refreshInvalid bool
}

func newAuthServer(t *testing.T) *authServer {
	t.Helper()
	s := &authServer{}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /v1/incidents/{id}", func(w http.ResponseWriter, r *http.Request) {
		s.headers = append(s.headers, r.Header.Get("Authorization"))
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{"id": r.PathValue("id"), "name": "DB down"})
	})
	mux.HandleFunc("GET /v1/unauthorized", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"detail":"The bearer token you provided is invalid or expired."}`))
	})
	mux.HandleFunc("POST /oauth/token", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		w.Header().Set("Content-Type", "application/json")
		if s.refreshInvalid || r.Form.Get("refresh_token") != "login-rt" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"error":"invalid_grant"}`))
			return
		}
		json.NewEncoder(w).Encode(map[string]any{"access_token": "refreshed-at", "refresh_token": "refreshed-rt", "token_type": "Bearer", "expires_in": 7200})
	})
	s.Server = httptest.NewServer(mux)
	t.Cleanup(s.Close)
	return s
}

func (s *authServer) lastHeader() string {
	if len(s.headers) == 0 {
		return ""
	}
	return s.headers[len(s.headers)-1]
}

func TestStoredLoginSendsBearerToLoginHost(t *testing.T) {
	srv := newAuthServer(t)
	storeLogin(t, loginCreds(srv.URL, time.Now().Add(time.Hour)))

	_, _, err := runRoot(t, NewRootCmd("test", "0.0.0-test"), "fhcli", "incidents", "get", "inc-1")

	if err != nil {
		t.Fatal(err)
	}
	if srv.lastHeader() != "Bearer login-at" {
		t.Errorf("Authorization = %q, want the login's bearer token", srv.lastHeader())
	}
}

func TestStoredLoginWorksForAPICommands(t *testing.T) {
	srv := newAuthServer(t)
	storeLogin(t, loginCreds(srv.URL, time.Now().Add(time.Hour)))

	_, _, err := runRoot(t, NewRootCmd("test", "0.0.0-test"), "fhcli", "api", "incidents", "get-incident", "inc-1")

	if err != nil {
		t.Fatal(err)
	}
	if srv.lastHeader() != "Bearer login-at" {
		t.Errorf("Authorization = %q", srv.lastHeader())
	}
}

func TestTokenFlagBeatsLogin(t *testing.T) {
	srv := newAuthServer(t)
	storeLogin(t, loginCreds(srv.URL, time.Now().Add(time.Hour)))

	_, err := runCmd(t, srv.URL, "incidents", "get", "inc-1")

	if err != nil || srv.lastHeader() != "test-token" {
		t.Errorf("Authorization = %q, err = %v; want --token to win", srv.lastHeader(), err)
	}
}

func TestEnvTokenBeatsLogin(t *testing.T) {
	srv := newAuthServer(t)
	storeLogin(t, loginCreds(srv.URL, time.Now().Add(time.Hour)))
	t.Setenv("FIREHYDRANT_API_TOKEN", "env-token")

	_, _, err := runRoot(t, NewRootCmd("test", "0.0.0-test"), "fhcli", "--api-host", srv.URL, "incidents", "get", "inc-1")

	if err != nil || srv.lastHeader() != "env-token" {
		t.Errorf("Authorization = %q, err = %v; want the env token to win", srv.lastHeader(), err)
	}
}

func TestEnvTokenIsReadByEveryRootCommand(t *testing.T) {
	srv := newAuthServer(t)
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	if _, err := runCmd(t, srv.URL, "incidents", "get", "inc-1"); err != nil {
		t.Fatal(err)
	}
	t.Setenv("FIREHYDRANT_API_TOKEN", "env-token")

	_, _, err := runRoot(t, NewRootCmd("test", "0.0.0-test"), "fhcli", "--api-host", srv.URL, "incidents", "get", "inc-1")

	if err != nil || srv.lastHeader() != "env-token" {
		t.Errorf("Authorization = %q, err = %v; a second root must still read FIREHYDRANT_API_TOKEN", srv.lastHeader(), err)
	}
}

func TestConfigFileTokenBeatsLogin(t *testing.T) {
	srv := newAuthServer(t)
	storeLogin(t, loginCreds(srv.URL, time.Now().Add(time.Hour)))
	cfg := filepath.Join(t.TempDir(), "firehydrant.cfg")
	os.WriteFile(cfg, []byte("apiKey: config-token\n"), 0o600)
	// loadConfigFile exports config values as env vars; t.Setenv makes the
	// test restore FH_API_KEY to unset afterwards.
	t.Setenv("FH_API_KEY", "")
	os.Unsetenv("FH_API_KEY")

	_, _, err := runRoot(t, NewRootCmd("test", "0.0.0-test"), "fhcli", "--config", cfg, "--api-host", srv.URL, "incidents", "get", "inc-1")

	if err != nil || srv.lastHeader() != "config-token" {
		t.Errorf("Authorization = %q, err = %v; want the config token to win", srv.lastHeader(), err)
	}
}

func TestAPIHostFlagOverridesLoginHost(t *testing.T) {
	srv := newAuthServer(t)
	storeLogin(t, loginCreds("http://127.0.0.1:1", time.Now().Add(time.Hour)))

	_, _, err := runRoot(t, NewRootCmd("test", "0.0.0-test"), "fhcli", "--api-host", srv.URL, "incidents", "get", "inc-1")

	if err != nil || srv.lastHeader() != "Bearer login-at" {
		t.Errorf("Authorization = %q, err = %v; want the login sent to --api-host", srv.lastHeader(), err)
	}
}

func TestAPIHostEnvOverridesLoginHost(t *testing.T) {
	srv := newAuthServer(t)
	storeLogin(t, loginCreds("http://127.0.0.1:1", time.Now().Add(time.Hour)))
	t.Setenv("FH_API_HOST", srv.URL)

	_, _, err := runRoot(t, NewRootCmd("test", "0.0.0-test"), "fhcli", "incidents", "get", "inc-1")

	if err != nil || srv.lastHeader() != "Bearer login-at" {
		t.Errorf("Authorization = %q, err = %v; want the login sent to FH_API_HOST", srv.lastHeader(), err)
	}
}

func TestExpiredLoginRefreshesBeforeCall(t *testing.T) {
	srv := newAuthServer(t)
	storeLogin(t, loginCreds(srv.URL, time.Now().Add(-time.Minute)))

	_, _, err := runRoot(t, NewRootCmd("test", "0.0.0-test"), "fhcli", "incidents", "get", "inc-1")

	if err != nil || srv.lastHeader() != "Bearer refreshed-at" {
		t.Errorf("Authorization = %q, err = %v; want the refreshed token", srv.lastHeader(), err)
	}
	stored, _ := auth.DefaultStore().Load("us")
	if stored.RefreshToken != "refreshed-rt" {
		t.Errorf("rotated refresh token not saved: %q", stored.RefreshToken)
	}
}

func TestLoginExpiredIsReported(t *testing.T) {
	srv := newAuthServer(t)
	srv.refreshInvalid = true
	storeLogin(t, loginCreds(srv.URL, time.Now().Add(-time.Minute)))

	_, stderr, err := runRoot(t, NewRootCmd("test", "0.0.0-test"), "fhcli", "incidents", "get", "inc-1")

	if err == nil || !strings.Contains(stderr, "login expired") || !strings.Contains(stderr, "fhcli login") {
		t.Errorf("err = %v, stderr = %q; want a login-expired error", err, stderr)
	}
	if len(srv.headers) != 0 {
		t.Error("no API request should be sent with an expired login")
	}
}

func TestExpiredLoginDoesNotFailEvent(t *testing.T) {
	srv := newAuthServer(t)
	srv.refreshInvalid = true
	storeLogin(t, loginCreds(srv.URL, time.Now().Add(-time.Minute)))

	_, _, err := runRoot(t, NewRootCmd("test", "0.0.0-test"), "fhcli", "event", "deployed")

	if err != nil {
		t.Errorf("event must exit 0 by default even when the login is expired, got %v", err)
	}
}

func TestCorruptLoginIsReported(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	dir := filepath.Join(auth.ConfigDir(), "credentials")
	os.MkdirAll(dir, 0o700)
	os.WriteFile(filepath.Join(dir, "us.json"), []byte("garbage"), 0o600)

	_, stderr, err := runRoot(t, NewRootCmd("test", "0.0.0-test"), "fhcli", "incidents", "get", "inc-1")

	if err == nil || !strings.Contains(stderr, "fhcli login") {
		t.Errorf("err = %v, stderr = %q; want an error pointing at fhcli login", err, stderr)
	}
}

func TestNoCredentialsMentionsLogin(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())

	_, stderr, err := runRoot(t, NewRootCmd("test", "0.0.0-test"), "fhcli", "incidents", "get", "inc-1")

	if err == nil || !strings.Contains(stderr, "API token is required") || !strings.Contains(stderr, "fhcli login") {
		t.Errorf("err = %v, stderr = %q", err, stderr)
	}
}

func TestUnauthorizedAddsLoginHint(t *testing.T) {
	var calls atomic.Int32
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls.Add(1)
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"detail":"The bearer token you provided is invalid or expired."}`))
	}))
	defer ts.Close()

	for _, args := range [][]string{
		{"incidents", "get", "inc-1"},                 // SDK client
		{"api", "incidents", "get-incident", "inc-1"}, // raw client
	} {
		_, stderr, err := runCmdFull(t, ts.URL, args...)
		if err == nil || !strings.Contains(stderr, "fhcli login") {
			t.Errorf("%v: stderr = %q; want the fhcli login hint on a 401", args, stderr)
		}
	}
}

func TestOtherErrorsHaveNoLoginHint(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"detail":"not found"}`))
	}))
	defer ts.Close()

	_, stderr, _ := runCmdFull(t, ts.URL, "api", "incidents", "get-incident", "inc-1")

	if strings.Contains(stderr, "fhcli login") {
		t.Errorf("a 404 should not suggest logging in: %q", stderr)
	}
}
