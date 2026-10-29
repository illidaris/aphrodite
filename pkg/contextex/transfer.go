package contextex

import (
	"context"

	"github.com/google/uuid"
	"github.com/illidaris/core"
	coreContextex "github.com/illidaris/core/contextex"
	"github.com/spf13/cast"
)

// TransferContext copies the trace ID and requested values from src into to.
// Unrequested values and cancellation behavior are inherited from to.
func TransferContext(src, to context.Context, keys ...any) context.Context {
	ctx := to
	if traceID := core.TraceID.Get(src); traceID != nil {
		ctx = core.TraceID.Set(ctx, traceID)
	}
	if traceID := coreContextex.TraceIdFrmCtx(src); traceID != "" {
		ctx = coreContextex.WithTraceId(ctx, traceID)
	}
	for _, key := range keys {
		if v := src.Value(key); v != nil {
			ctx = context.WithValue(ctx, key, v)
		}
	}
	return ctx
}

// TransferBackground creates a context based on context.Background, preserving
// trace metadata and the explicitly selected values from ctx. It assigns a new
// session ID and does not inherit cancellation or deadlines from ctx.
func TransferBackground(ctx context.Context, keys ...any) context.Context {
	var (
		newCtx    = context.Background() // new background ctx
		sessionId = uuid.NewString()     // new session id
	)
	// set trace id
	if traceId := coreContextex.TraceIdFrmCtx(ctx); traceId != "" {
		newCtx = coreContextex.WithTraceId(newCtx, traceId)
		newCtx = core.TraceID.Set(newCtx, traceId)
	} else if v := ctx.Value(core.TraceID); v != nil {
		newCtx = coreContextex.WithTraceId(newCtx, cast.ToString(v))
		newCtx = core.TraceID.Set(newCtx, v)
	} else {
		newCtx = coreContextex.WithTraceId(newCtx, sessionId)
		newCtx = core.TraceID.Set(newCtx, sessionId)
	}
	// set session id
	newCtx = coreContextex.WithSessionId(newCtx, sessionId)
	newCtx = core.SessionID.Set(newCtx, sessionId)

	for _, key := range keys {
		if v := ctx.Value(key); v != nil {
			newCtx = context.WithValue(newCtx, key, v)
		}
	}
	return newCtx
}
