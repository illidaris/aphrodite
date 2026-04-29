package imex

import (
	"context"

	"github.com/illidaris/aphrodite/pkg/dependency"
	group "github.com/illidaris/aphrodite/pkg/group/v2"
)

func BasePaged[Req dependency.IPage, Resp dependency.IPaginator, Item any](
	ctx context.Context,
	req Req,
	exportFunc func(context.Context, Req) (Resp, error),
	pagesFunc func(Req, Resp) []Req,
	getItemFunc func(Resp) []Item,
	opts ...ImExOptionFunc[Item],
) (<-chan Item, error) {
	opt := NewImExOption[Item]()
	for _, f := range opts {
		f(opt)
	}
	return basePaged(ctx, req, exportFunc, pagesFunc, getItemFunc, opt)
}

func basePaged[Req dependency.IPage, Resp dependency.IPaginator, Item any](
	ctx context.Context,
	req Req,
	exportFunc func(context.Context, Req) (Resp, error),
	pagesFunc func(Req, Resp) []Req,
	getItemFunc func(Resp) []Item,
	opt *ImExOption[Item],
) (<-chan Item, error) {
	inCh := make(chan Item, 10)
	firstResp, ex := exportFunc(ctx, req)
	if ex != nil {
		return inCh, ex
	}
	pages := pagesFunc(req, firstResp)
	go func() {
		defer func() {
			close(inCh)
		}()
		defer func() {
			if r := recover(); r != nil {
				println(r)
			}
		}()
		for _, v := range getItemFunc(firstResp) {
			inCh <- v
		}
		_, _ = group.GroupFunc(func(subReqs ...Req) (int64, error) {
			affect := 0
			for _, subReq := range subReqs {
				resp, ex := exportFunc(ctx, subReq)
				if ex != nil {
					continue
				}
				for _, v := range getItemFunc(resp) {
					inCh <- v
				}
				affect++
			}
			return int64(affect), nil
		}, pages, opt.GroupOptions...)
	}()
	return inCh, nil
}
