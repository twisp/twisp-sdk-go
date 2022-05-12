package client

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Khan/genqlient/graphql"
	"github.com/imdario/mergo"
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

func NewTwispHttp(authorization, customerAccount string) *http.Client {
	httpClient := http.Client{
		Transport: &authedTransport{
			jwt:             string(authorization),
			xtwispAccountID: customerAccount,
			wrapped:         http.DefaultTransport,
		},
	}

	return &httpClient
}

//NewTwispClient allows caller to override variables or merge objects.
func NewTwispClient(endpoint string, variables map[string]any, httpClient *http.Client) graphql.Client {
	return &twispClient{
		variables: variables,
		wrapped:   graphql.NewClient(endpoint, httpClient),
	}
}

type twispClient struct {
	variables map[string]any
	wrapped   graphql.Client
}

// MakeRequest implements graphql.Client
func (tc *twispClient) MakeRequest(ctx context.Context, req *graphql.Request, resp *graphql.Response) error {
	if len(tc.variables) == 0 {
		return tc.wrapped.MakeRequest(ctx, req, resp)
	}

	if req.Variables == nil {
		req.Variables = tc.variables
		return tc.wrapped.MakeRequest(ctx, req, resp)
	}

	err := tc.merge(req)
	if err != nil {
		return err
	}

	return tc.wrapped.MakeRequest(ctx, req, resp)
}

func (tc *twispClient) merge(req *graphql.Request) error {
	var variables any
	toMerge := []any{tc.variables, req.Variables}
	for _, vars := range toMerge {
		err := mergo.Merge(&variables, vars)
		if err != nil {
			return err
		}
	}
	req.Variables = variables
	return nil
}

var _ graphql.Client = (*twispClient)(nil)
