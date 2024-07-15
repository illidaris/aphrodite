package idgenerate

import (
	"context"

	"github.com/illidaris/aphrodite/idgenerate/dep"
)

func IDX(generater dep.IIDGenerate) func(ctx context.Context, key string) uint64 {
	return func(ctx context.Context, key string) uint64 {
		return generater.NewIDX(ctx, key)
	}
}
