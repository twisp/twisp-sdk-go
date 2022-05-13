package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Khan/genqlient/graphql"
	"github.com/imdario/mergo"
)

//authedTransport puts the authorization headers in the correct
//spots on the client.
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

//NewTwispHttp returns an *http.Client that sets authorization and x-twisp-account-id headers.
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

//NewTwispClient implements a graphql.Client that allows override/merging of Variables sent by
//the client.  This allows Twisp to use graphql variables without having to have every query
//have typed inputs.
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

// MakeRequest will wrapps up the standard genqlient graphql.Client but:
// 1. Adds variables to standard request if none are set.
// 2. If variables and req.Variables are set, merges together favoring req.Variables.
func (tc *twispClient) MakeRequest(ctx context.Context, req *graphql.Request, resp *graphql.Response) error {
	if len(tc.variables) == 0 {
		return tc.wrapped.MakeRequest(ctx, req, resp)
	}

	if req.Variables == nil {
		req.Variables = tc.variables
		return tc.wrapped.MakeRequest(ctx, req, resp)
	}

	err := merge(tc.variables, req)
	if err != nil {
		return err
	}

	return tc.wrapped.MakeRequest(ctx, req, resp)
}

func merge(variables map[string]any, req *graphql.Request) error {
	var finalVariables map[string]any
	toMerge := []any{variables, req.Variables}
	for _, vars := range toMerge {
		var asMap map[string]any
		if varsMap, ok := vars.(map[string]any); ok {
			asMap = varsMap
		} else {
			b, err := json.Marshal(vars)
			if err != nil {
				return err
			}

			err = json.Unmarshal(b, &asMap)
			if err != nil {
				return err
			}
		}

		err := mergo.Merge(&finalVariables, asMap, mergo.WithOverride)
		if err != nil {
			return err
		}
	}
	req.Variables = finalVariables
	return nil
}

var _ graphql.Client = (*twispClient)(nil)
