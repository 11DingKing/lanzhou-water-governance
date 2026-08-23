package audit

import "context"

type key struct{}

func WithRequestID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, key{}, id)
}
func RequestID(ctx context.Context) string { value, _ := ctx.Value(key{}).(string); return value }
