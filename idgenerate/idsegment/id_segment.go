package idsegment

import (
	"context"
	"errors"

	"github.com/illidaris/aphrodite/idgenerate/dep"
)

var _ = dep.IIDGenerate(&IdSegment{})

type IdSegment struct {
	Batch int64
	Cache ICache
	Repo  IRepository
}

func (i IdSegment) NewIDX(ctx context.Context, key string) int64 {
	id, err := i.NewID(ctx, key)
	if err != nil {
		return 0
	}
	return id
}

func (i IdSegment) NewID(ctx context.Context, key string) (int64, error) {
	seg, _, err := i.NewSegment(ctx, key)
	if err != nil {
		return 0, err
	}
	if seg == nil {
		return 0, errors.New("seg is nil")
	}
	if seg.Code != StatusCodeNil {
		return seg.Cursor, seg.Code.ToError()
	}
	return seg.Cursor, nil
}

func (i IdSegment) NewIDIterate(ctx context.Context, iterate func(int64), key string, opts ...dep.Option) error {
	seg, supseg, err := i.NewSegment(ctx, key, opts...)
	if err != nil {
		return err
	}
	if seg != nil && iterate != nil {
		for i := seg.MinId; i <= seg.Cursor; i++ {
			iterate(i)
		}
	}
	if supseg != nil && iterate != nil {
		for i := supseg.MinId; i <= supseg.Cursor; i++ {
			iterate(i)
		}
	}
	return nil
}

func (i IdSegment) NewSegment(ctx context.Context, key string, opts ...dep.Option) (*Segment, *Segment, error) {
	o := dep.NewOptions(opts...)
	if o.Num > i.Batch/2 {
		b, e, _, err := i.Repo.BlockNextSegment(ctx, key, o.Num, nil)
		if err != nil {
			return nil, nil, err
		}
		return &Segment{Code: StatusCodeNil, MinId: b, MaxId: e, Cursor: e}, nil, nil
	}
	seg, err := i.GenerateSegment(ctx, key, o.Num)
	if err != nil {
		return nil, nil, err
	}
	if seg.Code != StatusCodeNil {
		return seg, nil, seg.Code.ToError()
	}
	if hasNum := seg.Cursor - seg.MinId; hasNum < o.Num {
		supseg, err := i.GenerateSegment(ctx, key, o.Num-hasNum)
		if err != nil {
			return seg, supseg, err
		}
		if supseg.Code != StatusCodeNil {
			return seg, supseg, supseg.Code.ToError()
		}
		return seg, supseg, nil
	}
	return seg, nil, nil
}

func (i IdSegment) GenerateSegment(ctx context.Context, key string, num int64) (*Segment, error) {
	res, err := i.Cache.Eval(ctx, LUASCRIPT_HINCR, []string{key}, num)
	if err != nil {
		return nil, err
	}
	seg := parseLuaResult(res)
	switch seg.Code {
	case StatusCodeHUninit, StatusCodeOverflow:
		f := func() (*Segment, error) {
			res, err := i.Cache.Eval(ctx, LUASCRIPT_HINCR, []string{key}, num)
			if err != nil {
				return nil, err
			}
			seg := parseLuaResult(res)
			return seg, nil
		}
		return i.ReplGenerate(ctx, key, i.Batch, f)
	default:
		return seg, seg.Code.ToError()
	}
}

func (i IdSegment) ReplGenerate(ctx context.Context, key string, num int64, tryGenerate func() (*Segment, error)) (*Segment, error) {
	// 阻塞行锁获取新的号段
	b, e, res, err := i.Repo.BlockNextSegment(ctx, key, num, tryGenerate)
	if err != nil {
		return nil, err
	}
	if res != nil {
		return res, nil
	}
	result, err := i.Cache.Eval(ctx, LUASCRIPT_HREPL, []string{key}, b, e)
	if err != nil {
		return parseLuaResult(result), err
	}
	seg := parseLuaResult(result)
	if seg.Code != StatusCodeNil {
		return seg, seg.Code.ToError()
	}
	return tryGenerate()
}
