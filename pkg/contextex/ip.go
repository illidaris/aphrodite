package contextex

import (
	"context"

	"github.com/illidaris/aphrodite/pkg/netex"
)

type CtxIPKey struct{}

var ctxKeyIP = CtxIPKey{}

func WithIP(ctx context.Context, v string) context.Context {
	return context.WithValue(ctx, ctxKeyIP, v)
}

func GetIP(ctx context.Context) string {
	v := ctx.Value(ctxKeyIP)
	if v == nil {
		return ""
	}
	p, ok := v.(string)
	if !ok {
		return ""
	}
	return p
}

func GetIPInt(ctx context.Context) uint32 {
	v, err := netex.IPv4ToInt(GetIP(ctx))
	if err != nil {
		return 0
	}
	return v
}
