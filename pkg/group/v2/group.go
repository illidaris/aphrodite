package v2

import (
	"fmt"
	"sync"
	"sync/atomic"
)

// Group 将源切片按指定选项分组为多个子切片
// srcs: 待分组的源数据切片
// opts: 可选的分组配置参数，可设置批次大小/并发数等
// 返回值: 分组后的二维切片，每个子切片长度不超过配置的批次大小
func Group[T any](srcs []T, opts ...Option) [][]T {
	return group(srcs, newOptions(opts...))
}

// group 实现核心分组逻辑，根据配置选项将源切片划分为多个子切片
// srcs: 待分组的源数据切片
// opt: 已初始化的分组配置选项
// 返回值: 分组后的二维切片
func group[T any](srcs []T, opt *options) [][]T {
	groups := [][]T{}

	total := len(srcs)
	if total == 0 {
		return groups
	}

	// 根据配置计算单批数量和总分组数
	batch, gCount := batchByCount(srcs, opt)
	if gCount == 0 {
		return groups
	}

	// 按批次切割原切片并收集非空分组
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

// GroupFunc 并发执行分组函数并汇总处理结果
// f: 待执行的处理函数，接收可变参数返回影响数和错误
// srcs: 待处理的源数据切片
// opts: 可选的分组配置参数
// 返回值1: 所有分组处理的总影响数
// 返回值2: 包含错误信息的映射表，key为分组索引
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

	// 获取分组后的参数列表
	groups := group(srcs, opt)

	// 为每个分组启动协程执行处理函数
	for index, g := range groups {
		wg.Add(1)
		go func(i int, param ...T) {
			var err error
			defer wg.Done()

			// 双defer保证panic恢复和错误存储顺序
			defer func() {
				if err != nil {
					errs.Store(i, err)
				}
			}()

			defer func() {
				if r := recover(); r != nil {
					err = fmt.Errorf("err:%v", r)
				}
			}()

			// 执行实际处理函数并累加影响数
			affect, err := f(param...)
			atomic.AddInt64(&affectTotal, affect)
		}(index, g...)
	}
	wg.Wait()

	// 将并发安全的sync.Map转换为标准map
	errs.Range(func(key, value any) bool {
		k := key.(int)
		v := value.(error)
		errM[k] = v
		return true
	})
	return affectTotal, errM
}
