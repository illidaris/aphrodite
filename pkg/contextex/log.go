package contextex

import (
	"context"

	coreContextex "github.com/illidaris/core/contextex"
)

// WithLogIgnore marks ctx so logging helpers can skip writing a log entry.
func WithLogIgnore(ctx context.Context) context.Context {
	return coreContextex.WithLogIgnore(ctx)
}

// IsLogIgnoreFrmCtx reports whether ctx is marked to skip logging.
func IsLogIgnoreFrmCtx(ctx context.Context) bool {
	return coreContextex.IsLogIgnoreFrmCtx(ctx)
}
