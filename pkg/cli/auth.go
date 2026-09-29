package cli

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/firehydrant/fhcli/pkg/auth"
	"github.com/urfave/cli/v3"
)

const (
	authInfoKey ctxKey = "authInfo"
	authErrKey  ctxKey = "authErr"
)

var errNoCredentials = errors.New("API token is required. Run `fhcli login`, or set the --token flag, the FIREHYDRANT_API_TOKEN env var, or a config file")

// authInfo is how this run authenticates.
type authInfo struct {
	source  string // "api_token" or "login"
	header  string // the Authorization header value
	apiHost string
	creds   *auth.Credentials // set when source is "login"
}

// resolveAuth picks credentials in precedence order: an API token (flag,
// env var, or config file) first, then the stored login for the current
// region. It returns nil, nil when there are none.
func resolveAuth(ctx context.Context, cmd *cli.Command) (*authInfo, error) {
	if token := cmd.String("token"); token != "" {
		return &authInfo{source: "api_token", header: token, apiHost: resolveAPIHost(cmd)}, nil
	}
	creds, err := auth.Token(ctx, auth.DefaultStore(), auth.CurrentRegion())
	if errors.Is(err, auth.ErrNotLoggedIn) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	host := creds.APIHost
	if cmd.IsSet("api-host") {
		host = resolveAPIHost(cmd)
	}
	return &authInfo{source: "login", header: "Bearer " + creds.AccessToken, apiHost: host, creds: creds}, nil
}

// sawUnauthorized records whether any API response in this run was a 401,
// so reportError can suggest `fhcli login` whichever client made the call.
var sawUnauthorized atomic.Bool

type unauthorizedRecorder struct {
	next http.RoundTripper
}

func (u unauthorizedRecorder) RoundTrip(r *http.Request) (*http.Response, error) {
	resp, err := u.next.RoundTrip(r)
	if err == nil && resp.StatusCode == http.StatusUnauthorized {
		sawUnauthorized.Store(true)
	}
	return resp, err
}

// apiHTTPClient is shared by the SDK and rawClient.
var apiHTTPClient = &http.Client{
	Timeout:   60 * time.Second,
	Transport: unauthorizedRecorder{next: http.DefaultTransport},
}

func withLoginHint(err error) error {
	if !sawUnauthorized.Load() {
		return err
	}
	return fmt.Errorf("%w (run `fhcli login`, or check your API token)", err)
}
