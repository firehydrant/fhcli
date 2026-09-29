package cli

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

const defaultAPIHost = "https://api.firehydrant.io"

const rawClientKey ctxKey = "rawClient"

// rawClient calls the FireHydrant API directly and returns the response body
// untouched. Use it where the generated SDK's response types don't match what
// the API actually returns (the SDK silently drops unknown fields, so the
// command would print "{}").
type rawClient struct {
	baseURL string
	token   string
	http    *http.Client
}

func getRawClient(ctx context.Context) (*rawClient, error) {
	client, ok := ctx.Value(rawClientKey).(*rawClient)
	if !ok || client == nil {
		return nil, errors.New("API token is required. Set --token flag, FIREHYDRANT_API_TOKEN env var, or configure a config file")
	}
	return client, nil
}

// get issues a GET to path (e.g. "/v1/incidents/abc/transcript") and returns
// the JSON body. Non-2xx responses become errors that include the API's body.
func (c *rawClient) get(ctx context.Context, path string, query url.Values) (json.RawMessage, error) {
	u := strings.TrimRight(c.baseURL, "/") + path
	if len(query) > 0 {
		u += "?" + query.Encode()
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", c.token)
	req.Header.Set("Accept", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return nil, fmt.Errorf("GET %s: API returned %d: %s", path, resp.StatusCode, strings.TrimSpace(string(body)))
	}

	if !json.Valid(body) {
		return nil, fmt.Errorf("GET %s: API returned non-JSON body", path)
	}

	return body, nil
}
