package dependency

import (
	"context"
	"net/url"

	"github.com/illidaris/aphrodite/pkg/contextex"
)

// IBizRequest business tag reqquest
type IBizRequest interface {
	GetBizId() int64
	ToUrlQuery() url.Values
}

type IBindRequest interface {
	IBiz
	IIP
}
type IBiz interface {
	SetBiz(int64)
}

func BindBizFrmCtx(ctx context.Context, src any) any {
	if req, ok := src.(IBiz); ok {
		return BizFrmCtx(ctx, req)
	}
	return src
}
func BizFrmCtx(ctx context.Context, req IBiz) IBiz {
	req.SetBiz(contextex.GetBizId(ctx))
	return req
}

type IIP interface {
	SetIP(string)
}

func BindIPFrmCtx(ctx context.Context, src any) any {
	if req, ok := src.(IIP); ok {
		req.SetIP(contextex.GetIP(ctx))
	}
	return src
}

func IPFrmCtx(ctx context.Context, req IIP) IIP {
	req.SetIP(contextex.GetIP(ctx))
	return req
}
