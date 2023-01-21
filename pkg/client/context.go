package client

import "context"

type TwispContextKeyType string

var TwispContextKey = TwispContextKeyType("ctx")

// WithVariables puts the `variables` on the context and will mix in with the all other
// variables leaving twisp.
func WithVariables(ctx context.Context, variables map[string]interface{}) context.Context {
	return context.WithValue(ctx, TwispContextKey, variables)
}
