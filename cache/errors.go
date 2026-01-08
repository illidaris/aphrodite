package cache

import "errors"

var (
	ErrLimit    = errors.New("达到上限限制")
	ErrCacheNil = errors.New("缓存配置错误")
)
