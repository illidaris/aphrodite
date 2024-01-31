package group

import (
	"fmt"
	"sync"
	"sync/atomic"
)

// Count函数用于计算切片p中元素的批次数量。
// 参数batch表示每个批次的元素数量。
// 参数p为切片指针，切片中的元素类型为any。
// 返回值为int类型，表示切片p中的元素能够组成的批次数量。
func Count[T any](batch int, p ...*T) int {
	total := len(p) // 获取切片p的长度
	if total == 0 { // 判断切片p是否为空
		return 0 // 返回0表示无法组成批次
	}
	if int(total)%batch == 0 { // 判断切片p的元素数量能否整除batch
		return int(total) / batch // 返回能够整除的结果，即批次数量
	}
	return int(total)/batch + 1 // 返回无法整除的结果，即批次数量加一
}

// Group函数将给定的切片p按照batch的大小进行分组，并返回分组后的结果。
// 参数：
//   - batch：每个分组的大小
//   - p：需要分组的切片
//
// 返回值：
//   - [][]*T：分组后的结果，每个子切片表示一个分组
func Group[T any](batch int, p ...*T) [][]*T {
	groups := [][]*T{}
	total := len(p)
	if p == nil || total == 0 {
		return groups
	}
	gCount := Count[T](batch, p...)
	if gCount == 0 {
		return groups
	}
	for i := 0; i < gCount; i++ {
		beg := i * batch
		end := beg + batch
		if end > total {
			end = total
		}
		gp := p[beg:end]
		if len(gp) > 0 {
			groups = append(groups, gp)
		}
	}
	return groups
}

// GroupFunc是一个并发执行函数的工具函数。
// 它接受一个函数f作为参数，该函数接受可变数量的指向T类型的指针作为参数，并返回一个int64和一个error。
// batch参数指定了并发执行的批次大小。
// p参数是一个可变参数列表，用于指定函数f的参数。
// 函数返回一个int64和一个map[int]error，分别表示所有并发执行的函数的总影响量和每个并发执行的函数的错误信息。
func GroupFunc[T any](f func(v ...*T) (int64, error), batch int, p ...*T) (int64, map[int]error) {
	var (
		wg          sync.WaitGroup
		errs        sync.Map
		affectTotal int64
		errM        = map[int]error{}
	)
	groups := Group[T](batch, p...) // 调用Group函数，将参数p划分为多个批次，并返回每个批次的参数列表

	// 遍历每个批次的参数列表
	for index, g := range groups {
		wg.Add(1)
		go func(i int, param ...*T) {
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
