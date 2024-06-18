package idsegment

import (
	"context"
)

type ICache interface {
	Eval(ctx context.Context, script string, keys []string, args ...interface{}) (interface{}, error)
	EvalSha(ctx context.Context, script string, keys []string, args ...interface{}) (interface{}, error)
}

type IRepository interface {
	// BlockNextSegment 利用事务以及行锁阻塞执行，同时在执行中加入尝试生成函数，如果成功执行，则表示不需要生成，执行回滚
	BlockNextSegment(ctx context.Context, key string, step uint64, tryGenerate func() (*Segment, error)) (uint64, uint64, *Segment, error)
}
