package when

import (
	"context"
	"time"
)

type clockOffsetCtxKey struct{}

func WithClockOffset(ctx context.Context, v time.Duration) context.Context {
	return context.WithValue(ctx, clockOffsetCtxKey{}, v)
}

func Now(ctx context.Context) time.Time {
	v := ctx.Value(clockOffsetCtxKey{})
	if v == nil {
		return time.Now()
	}
	offset, ok := v.(time.Duration)
	if ok {
		return time.Now().Add(offset)
	}
	return time.Now()
}
