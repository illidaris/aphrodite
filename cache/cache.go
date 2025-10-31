package cache

import (
	"context"
	"time"
)

type ICache interface {
	EvalContext(ctx context.Context, script string, keys []string, args ...any) (any, error)
	SetContext(ctx context.Context, key string, val any, timeout time.Duration) (string, error)
	SetNXContext(ctx context.Context, key string, val any, timeout time.Duration) (bool, error)
	TTLContext(ctx context.Context, key string) (time.Duration, error)
	GetContext(ctx context.Context, key string) (any, error)
	DelContext(ctx context.Context, key string) (int64, error)
	LPushContext(ctx context.Context, key string, vals ...any) (int64, error)
	RPushContext(ctx context.Context, key string, vals ...any) (int64, error)
	LRangeContext(ctx context.Context, key string, start, stop int64) ([]string, error)
	TTL(key string) time.Duration
}
