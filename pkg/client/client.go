package client

import (
	"fmt"
	"net/http"
)

type authedTransport struct {
	jwt             string
	xtwispAccountID string
	wrapped         http.RoundTripper
}

func (t *authedTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", t.jwt))
	if t.xtwispAccountID != "" {
		req.Header.Set("X-Twisp-Account-Id", t.xtwispAccountID)
	}
	return t.wrapped.RoundTrip(req)
}

func NewTwispClient(authorization, customerAccount string) *http.Client {
	httpClient := http.Client{
		Transport: &authedTransport{
			jwt:             string(authorization),
			xtwispAccountID: customerAccount,
			wrapped:         http.DefaultTransport,
		},
	}

	return &httpClient
}
