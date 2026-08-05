package contextex

import "context"

type CtxLogIgnoreKey struct{}

var ctxLogIgnoreKey CtxLogIgnoreKey

// 根据上下文，不打日志
func WithLogIgnore(ctx context.Context) context.Context {
	v := ctx.Value(ctxLogIgnoreKey)
	if v != nil {
		return ctx
	}
	return context.WithValue(ctx, ctxLogIgnoreKey, true)
}

func IsLogIgnoreFrmCtx(ctx context.Context) bool {
	v := ctx.Value(ctxLogIgnoreKey)
	if v == nil {
		return false
	}
	b, ok := v.(bool)
	if !ok {
		return false
	}
	return b
}
