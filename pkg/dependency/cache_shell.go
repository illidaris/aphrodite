package dependency

import "time"

// ICacheShellKey interface defines the fundamental operations for a cache shell key.
// This interface is used to retrieve the cache key, cache duration, and a flag indicating whether to skip the cache.
type ICacheShellKey interface {
	// GetCacheKey returns the cache's key value.
	// Returns a string which represents a unique identifier for the cache.
	// Suggest example:"{AppName}:{BizId}:{BusinessName}:{EntityName/ObjectName}:{OtherKey}"
	GetCacheKey() string

	// GetCacheDuration returns the cache's duration.
	// Returns a time.Duration indicating the cache's expiration time.
	GetCacheDuration() time.Duration

	// GetSkip returns a flag that indicates if the cache should be skipped.
	// Returns a bool where true signifies that the current request should bypass the cache and fetch data directly from the source.
	GetSkip() bool
}
