package v2

import (
	"fmt"
	"sync"
	"sync/atomic"
)

func Group[T any](srcs []T, opts ...Option) [][]T {
	return group(srcs, newOptions(opts...))
}

func group[T any](srcs []T, opt *options) [][]T {
	groups := [][]T{}

	total := len(srcs)
	if total == 0 {
		return groups
	}
	batch, gCount := batchByCount(srcs, opt)
	if gCount == 0 {
		return groups
	}
	for i := 0; i < gCount; i++ {
		beg := i * batch
		end := beg + batch
		if end > total {
			end = total
		}
		gp := srcs[beg:end]
		if len(gp) > 0 {
			groups = append(groups, gp)
		}
	}
	return groups
}

func GroupFunc[T any](f func(v ...T) (int64, error), srcs []T, opts ...Option) (int64, map[int]error) {
	var (
		wg          sync.WaitGroup
		errs        sync.Map
		affectTotal int64
		errM        = map[int]error{}
		opt         = newOptions(opts...)
	)
	total := len(srcs)
	if total == 0 {
		return 0, errM
	}
	groups := group(srcs, opt) // 调用Group函数，将参数p划分为多个批次，并返回每个批次的参数列表
	// 遍历每个批次的参数列表
	for index, g := range groups {
		wg.Add(1)
		go func(i int, param ...T) {
			var err error
			defer wg.Done()

			// 函数执行前的恢复函数，用于捕获panic并将其转换为error
			defer func() {
				if err != nil {
					errs.Store(i, err)
				}
			}()

			// 函数执行前的恢复函数，用于捕获panic并将其转换为error
			defer func() {
				if r := recover(); r != nil {
					err = fmt.Errorf("err:%v", r)
				}
			}()

			// 调用函数f，并获取返回的总影响量和错误信息
			affect, err := f(param...)
			atomic.AddInt64(&affectTotal, affect)
		}(index, g...)
	}
	wg.Wait()

	// 遍历错误信息的map，并将错误信息存储到errM中
	errs.Range(func(key, value any) bool {
		k := key.(int)
		v := value.(error)
		errM[k] = v
		return true
	})
	return affectTotal, errM
}
