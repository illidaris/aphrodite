package contextex

import (
	"context"

	coreContextex "github.com/illidaris/core/contextex"
)

// WithBizId stores the business identifier in ctx.
func WithBizId(ctx context.Context, v int64) context.Context {
	return coreContextex.WithBizId(ctx, v)
}

// GetBizId returns the business identifier stored in ctx, or 0 if absent.
func GetBizId(ctx context.Context) int64 {
	return coreContextex.GetBizId(ctx)
}
