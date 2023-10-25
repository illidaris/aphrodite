package logex

import (
	"context"

	"github.com/illidaris/core"
	"go.uber.org/zap"
)

// FieldsFromCtx write context meta data
func FieldsFromCtx(ctx context.Context) []zap.Field {
	return []zap.Field{
		buildZapField(ctx, core.TraceID),
		buildZapField(ctx, core.SessionID),
		buildZapField(ctx, core.Action),
		buildZapField(ctx, core.Step),
	}
}

// buildZapField use core meta build field
func buildZapField(ctx context.Context, key core.MetaData) zap.Field {
	return zap.String(key.String(), key.GetString(ctx))
}
