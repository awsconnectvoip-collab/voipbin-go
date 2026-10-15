package voipbin

import (
	"net/http"

	"github.com/voipbin/voipbin-go/gens/voipbin_client"
)

const (
	defaultServerAddress = "https://api.voipbin.net/v1.0"
)

// NewClient creates a new VoIPBIN API client using access key authentication.
func NewClient(accesskey string) (voipbin_client.ClientWithResponsesInterface, error) {
	res, err := voipbin_client.NewClientWithResponses(defaultServerAddress, withAccessKey(accesskey))
	if err != nil {
		return nil, err
	}

	return res, nil
}

type customTransport struct {
	accessKey string
	base      http.RoundTripper
}

// RoundTrip modifies the request to add the access key as an Authorization header
// and delegates to the wrapped RoundTripper. This preserves connection reuse
// and transport settings (timeouts, TLS, proxy) and avoids calling
// http.DefaultClient.Do from RoundTrip.
func (t *customTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	base := t.base
	if base == nil {
		base = http.DefaultTransport
	}

	newReq := req.Clone(req.Context())
	// Prefer Authorization header to query parameter to avoid leaking keys in URLs
	if t.accessKey != "" {
		newReq.Header.Set("Authorization", "Bearer "+t.accessKey)
	}

	return base.RoundTrip(newReq)
}

func withAccessKey(accessKey string) voipbin_client.ClientOption {
	return func(c *voipbin_client.Client) error {
		// Ensure we wrap an existing transport instead of replacing the entire http.Client
		var base http.RoundTripper
		if c.Client != nil && c.Client.Transport != nil {
			base = c.Client.Transport
		} else {
			base = http.DefaultTransport
		}

		c.Client = &http.Client{
			Transport: &customTransport{
				accessKey: accessKey,
				base:      base,
			},
		}
		return nil
	}
}

type basicAuthTransport struct {
	username string
	password string
	base     http.RoundTripper
}

// RoundTrip adds Basic Auth header to the request and delegates to the wrapped transport.
func (t *basicAuthTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	base := t.base
	if base == nil {
		base = http.DefaultTransport
	}

	newReq := req.Clone(req.Context())
	if t.username != "" || t.password != "" {
		newReq.SetBasicAuth(t.username, t.password)
	}

	return base.RoundTrip(newReq)
}

func withBasicAuth(username, password string) voipbin_client.ClientOption {
	return func(c *voipbin_client.Client) error {
		var base http.RoundTripper
		if c.Client != nil && c.Client.Transport != nil {
			base = c.Client.Transport
		} else {
			base = http.DefaultTransport
		}

		c.Client = &http.Client{
			Transport: &basicAuthTransport{
				username: username,
				password: password,
				base:     base,
			},
		}
		return nil
	}
}

// StringPtr returns a pointer to the given string value.
func StringPtr(s string) *string { return &s }

// IntPtr returns a pointer to the given int value.
func IntPtr(i int) *int { return &i }
