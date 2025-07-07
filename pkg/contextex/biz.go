package contextex

import "context"

type CtxBizrKey struct{}

var ctxKeyBiz = CtxBizrKey{}

func WithBizId(ctx context.Context, v int64) context.Context {
	return context.WithValue(ctx, ctxKeyBiz, v)
}

func GetBizId(ctx context.Context) int64 {
	v := ctx.Value(ctxKeyBiz)
	if v == nil {
		return 0
	}
	p, ok := v.(int64)
	if !ok {
		return 0
	}
	return p
}
