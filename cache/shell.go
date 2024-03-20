package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/illidaris/aphrodite/pkg/dependency"
	"github.com/illidaris/aphrodite/pkg/exception"
	"github.com/spf13/cast"
)

// Define shell operation options and utility functions related to caching.

const (
	KEY_LOCK_SUFFIX = "_locked" // Suffix appended to cache keys for locking purposes
)

// ShellOptionFunc is a function type that configures ShellOptions.
type ShellOptionFunc func(option *ShellOptions)

// NewShellOptions initializes and returns a new ShellOptions instance.
// @param cache The dependency cache interface used for cache operations.
// @param request The dependency cache key request interface, used to retrieve cache keys and durations.
// @param opts One or more ShellOptionFuncs to further configure the Shell options.
// @return A configured ShellOptions instance.
func NewShellOptions(cache dependency.ICache, request dependency.ICacheShellKey, opts ...ShellOptionFunc) *ShellOptions {
	option := &ShellOptions{
		cache: cache,
		key:   request.GetCacheKey(),
		dur:   request.GetCacheDuration(),
		skip:  request.GetSkip(),
	}
	for _, opt := range opts {
		opt(option)
	}
	return option
}

// ShellOptions defines configuration options for shell operations.
type ShellOptions struct {
	cache dependency.ICache // Cache instance
	key   string            // Cache key
	dur   time.Duration     // Cache expiration duration
	skip  bool              // Whether to skip caching
}

// WithCache provides an option to customize the cache instance.
// @return A ShellOptionFunc to set the cache instance.
func WithCache(cache dependency.ICache) ShellOptionFunc {
	return func(option *ShellOptions) {
		option.cache = cache
	}
}

// WithKey provides an option to customize the cache key.
// @return A ShellOptionFunc to set the cache key.
func WithKey(key string) ShellOptionFunc {
	return func(option *ShellOptions) {
		option.key = key
	}
}

// WithDuration provides an option to customize the cache duration.
// @return A ShellOptionFunc to set the cache duration.
func WithDuration(dur time.Duration) ShellOptionFunc {
	return func(option *ShellOptions) {
		option.dur = dur
	}
}

// WithSkip provides an option to skip caching.
// @return A ShellOptionFunc to set whether to skip caching.
func WithSkip(skip bool) ShellOptionFunc {
	return func(option *ShellOptions) {
		option.skip = skip
	}
}

// ShellClear clears the specified cache key and its lock identifier.
// @param cache The cache instance to perform cache operations.
// @param request The cache key request interface to get the cache key.
// @param opts One or more ShellOptionFuncs to further configure the Shell options.
// @return An exception instance if there's an error clearing the cache.
func ShellClear(request dependency.ICacheShellKey, opts ...ShellOptionFunc) exception.Exception {
	option := NewShellOptions(nil, request, opts...)
	if option.cache == nil {
		return nil
	}
	keyLocked := option.key + KEY_LOCK_SUFFIX
	err := option.cache.Delete(keyLocked)
	if err != nil {
		return exception.ERR_BUSI.Wrap(err)
	}
	return nil
}

// Shell is a generic function that executes caching logic.
// @param ctx The context for logging and cancellation purposes.
// @param request The cache key request interface to get the cache key.
// @param f A function executed when cache is missed or unavailable, generating the result.
// @param opts One or more ShellOptionFuncs to further configure the Shell options.
// @return Returns the result of function f and a possible exception.
func Shell[T any](ctx context.Context, request dependency.ICacheShellKey, f func() (T, exception.Exception), opts ...ShellOptionFunc) (T, exception.Exception) {
	option := NewShellOptions(nil, request, opts...)
	if option.cache == nil || option.skip {
		return f()
	}
	key := option.key
	dur := option.dur
	cache := option.cache
	keyLocked := key + KEY_LOCK_SUFFIX
	if b, err := cache.SetNX(keyLocked, key, dur); err != nil || !b {
		response := new(T)
		cacheValue := cache.Get(key)
		resStr := cast.ToString(cacheValue)
		if len(resStr) > 0 {
			if err := json.Unmarshal([]byte(resStr), response); err != nil {
				logger().Warn(ctx, err.Error())
			}
		}
		logger().Info(ctx, "%s fallback to cache, value is %s", key, resStr)
		return *response, nil
	}
	res, ex := f()
	if ex != nil {
		defer func() {
			_ = cache.Delete(keyLocked)
		}()
		return res, ex
	}
	bs, _ := json.Marshal(res)
	if err := cache.Set(key, string(bs), dur*5); err != nil {
		logger().Warn(context.TODO(), err.Error())
	}
	return res, ex
}
