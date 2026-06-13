package request

import (
	"context"
	"reflect"
)

type payloadContextKey string

func WithValidatedPayload[T any](ctx context.Context, payload *T) context.Context {
	return context.WithValue(ctx, payloadKey[T](), payload)
}

func ValidatedPayload[T any](ctx context.Context) (*T, bool) {
	payload, ok := ctx.Value(payloadKey[T]()).(*T)
	return payload, ok
}

func payloadKey[T any]() payloadContextKey {
	return payloadContextKey(reflect.TypeOf((*T)(nil)).String())
}
