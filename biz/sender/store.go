package sender

import (
	"context"
	"time"
)

type IStore interface {
	EvalContext(ctx context.Context, script string, keys []string, args ...any) (any, error)
	SetContext(ctx context.Context, key string, val any, timeout time.Duration) (string, error)
	SetNXContext(ctx context.Context, key string, val any, timeout time.Duration) (bool, error)
	TTLContext(ctx context.Context, key string) (time.Duration, error)
	GetContext(ctx context.Context, key string) (any, error)
	DelContext(ctx context.Context, key string) (int64, error)
}
