package contextex

import (
	"context"

	"github.com/google/uuid"
	"github.com/illidaris/core"
)

// TransferContext transfer kv from src context to target context
func TransferContext(src, to context.Context, keys ...any) context.Context {
	v := core.TraceID.Get(src)
	ctx := core.TraceID.Set(to, v)
	for _, key := range keys {
		if v := src.Value(key); v != nil {
			ctx = context.WithValue(ctx, key, v)
		}
	}
	return ctx
}

// TransferBackground transfer kv from src context to background context
func TransferBackground(ctx context.Context, keys ...any) context.Context {
	var (
		newCtx    = context.Background() // new background ctx
		sessionId = uuid.NewString()     // new session id
	)
	// set trace id
	if v := ctx.Value(core.TraceID); v != nil {
		newCtx = core.TraceID.Set(newCtx, v)
	} else {
		newCtx = core.TraceID.Set(newCtx, sessionId)
	}
	// set session id
	newCtx = core.SessionID.Set(newCtx, sessionId)
	return TransferContext(ctx, newCtx, keys...)
}
