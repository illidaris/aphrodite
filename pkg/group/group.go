package group

import (
	"fmt"
	"sync"
	"sync/atomic"
)

// Count
func Count[T any](batch int, p ...T) int {
	total := len(p)
	if total == 0 {
		return 0
	}
	if int(total)%batch == 0 {
		return int(total) / batch
	}
	return int(total)/batch + 1
}

// Group
func Group[T any](batch int, p ...T) [][]T {
	groups := [][]T{}
	gCount := Count[T](batch, p...)
	if gCount == 0 {
		return groups
	}
	total := len(p)
	for i := 0; i < gCount; i++ {
		end := (i + 1) * batch
		if end > total {
			end = total
		}
		groups = append(groups, p[i*batch:end])
	}
	return groups
}

// GroupFunc
func GroupFunc[T any](f func(v ...T) (int64, error), batch int, p ...T) (int64, map[int]error) {
	var (
		wg          sync.WaitGroup
		errs        sync.Map
		affectTotal int64
		errM        = map[int]error{}
	)
	groups := Group[T](batch, p...)
	for index, g := range groups {
		wg.Add(1)
		go func(i int, param ...T) {
			var err error
			defer wg.Done()
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
			affect, err := f(param...)
			atomic.AddInt64(&affectTotal, affect)
		}(index, g...)
	}
	wg.Wait()
	errs.Range(func(key, value any) bool {
		k := key.(int)
		v := value.(error)
		errM[k] = v
		return true
	})
	return affectTotal, errM
}
