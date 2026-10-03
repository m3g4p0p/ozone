package ctxval

import "context"

type ctxKey[T any] struct{}

func From[T any](ctx context.Context) (T, bool) {
	v, ok := ctx.Value(ctxKey[T]{}).(T)
	return v, ok
}

func With[T any](ctx context.Context, v T) context.Context {
	return context.WithValue(ctx, ctxKey[T]{}, v)
}
