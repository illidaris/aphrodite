package dep

import "context"

type IIDGenerate interface {
	NewIDX(ctx context.Context, key string) uint64
	NewID(ctx context.Context, key string) (uint64, error)
	NewIDIterate(ctx context.Context, iterate func(uint64), key string, opts ...Option) error
}
