package idgenerate

import (
	"context"
	"hercules/pkg/idgenerate/dep"
)

func IDX(generater dep.IIDGenerate) func(ctx context.Context, key string) uint64 {
	return func(ctx context.Context, key string) uint64 {
		return generater.NewIDX(ctx, key)
	}
}
