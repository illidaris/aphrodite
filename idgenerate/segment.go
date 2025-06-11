package idgenerate

import (
	"context"

	"github.com/illidaris/aphrodite/idgenerate/dep"
)

func IDX(generater dep.IIDGenerate) func(ctx context.Context, key string) int64 {
	return func(ctx context.Context, key string) int64 {
		return generater.NewIDX(ctx, key)
	}
}
