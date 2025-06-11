package dep

import "context"

type IIDGenerate interface {
	NewIDX(ctx context.Context, key string) int64
	NewID(ctx context.Context, key string) (int64, error)
	NewIDIterate(ctx context.Context, iterate func(int64), key string, opts ...Option) error
}
