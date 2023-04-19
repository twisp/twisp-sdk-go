package client

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/Khan/genqlient/graphql"
	"github.com/imdario/mergo"
	"github.com/twisp/twisp-sdk-go/pkg/token"
)

var (
	Expired = time.Date(1900, time.January, 1, 0, 0, 0, 0, time.UTC)
)

func NewTwispRoundTripper(awsAccount, twispEnvironment, region string, now Now) http.RoundTripper {
	return &roundTripper{
		awsAccount:       awsAccount,
		twispEnvironment: twispEnvironment,
		region:           region,
		now:              now,
		m:                &sync.RWMutex{},
		expire:           Expired,
		auth:             []byte{},
		wrapped:          http.DefaultTransport,
	}
}

type roundTripper struct {
	awsAccount       string
	twispEnvironment string
	region           string

	now    func() time.Time
	auth   []byte
	expire time.Time
	m      *sync.RWMutex

	wrapped http.RoundTripper
}

func (r *roundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	jwt, err := r.authorization()
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return nil, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", string(jwt)))
	req.Header.Set("X-Twisp-Account-Id", r.awsAccount)

	return r.wrapped.RoundTrip(req)
}

// authorization returns an OIDC token from Twisp by exchanging the current IAM credentials.
// The OIDC token is cached until it expires.
func (t *roundTripper) authorization() ([]byte, error) {
	// Read cached version
	t.m.RLock()
	if len(t.auth) > 0 && t.now().Before(t.expire) {
		return t.auth, nil
	}
	t.m.RUnlock()

	// We need to get a new token. Take write lock.
	t.m.Lock()
	defer t.m.Unlock()

	b, err := token.Exchange(t.twispEnvironment, t.region)
	if err != nil {
		return nil, err
	}
	exp, err := extractExpire(string(b))
	if err != nil {
		return nil, err
	}

	t.auth = b
	t.expire = exp

	return t.auth, nil
}

// extractExpire extracts the exp claim from the jwt and returns it.
func extractExpire(jwt string) (time.Time, error) {
	parts := strings.Split(jwt, ".")
	if len(parts) < 2 {
		return Expired, fmt.Errorf("verify: malformed jwt, expected 3 parts got %d", len(parts))
	}

	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return Expired, fmt.Errorf("verify: malformed jwt payload: %v", err)
	}

	type idToken struct {
		Expire int64 `json:"exp"`
	}

	var token idToken
	if err := json.Unmarshal(payload, &token); err != nil {
		return Expired, fmt.Errorf("verify: failed to unmarshal claims: %v", err)
	}

	return time.Unix(token.Expire, 0), nil
}

type Now func() time.Time

// NewTwispHttp returns an *http.Client that sets authorization and x-twisp-account-id headers.
// example: NewTwispHttp("Twisp1234", "cloud", "us-east-1")
func NewTwispHttp(customerAccount, twispEnvironment, region string) *http.Client {
	httpClient := http.Client{
		Transport: &roundTripper{
			awsAccount:       customerAccount,
			twispEnvironment: twispEnvironment,
			region:           region,
			now:              time.Now,
			m:                &sync.RWMutex{},
			expire:           Expired,
			auth:             []byte{},
			wrapped:          http.DefaultTransport,
		},
	}

	return &httpClient
}

// NewTwispClient implements a graphql.Client that allows override/merging of Variables sent by
// the client.  This allows Twisp to use graphql variables without having to have every query
// have typed inputs.
func NewTwispClient(endpoint string, httpClient *http.Client) graphql.Client {
	return &twispClient{
		wrapped: graphql.NewClient(endpoint, httpClient),
	}
}

type twispClient struct {
	wrapped graphql.Client
}

// MakeRequest will wraps up the standard genqlient graphql.Client but:
// 1. Adds variables to standard request if none are set.
// 2. If variables and req.Variables are set, merges together favoring req.Variables.
func (tc *twispClient) MakeRequest(ctx context.Context, req *graphql.Request, resp *graphql.Response) error {
	vars := ctx.Value(TwispContextKey)
	if vars == nil {
		return tc.wrapped.MakeRequest(ctx, req, resp)
	}

	varsMap, ok := vars.(map[string]interface{})
	if !ok {
		return tc.wrapped.MakeRequest(ctx, req, resp)
	}

	if req.Variables == nil {
		req.Variables = varsMap
		return tc.wrapped.MakeRequest(ctx, req, resp)
	}

	err := merge(varsMap, req)
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
