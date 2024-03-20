package dependency

import "time"

// ICache interface defines the basic operations for a cache.
type ICache interface {
	// Get retrieves a value by its key.
	// key: The key of the cache item.
	// Returns: The value of the cache item, or nil if it doesn't exist.
	Get(key string) any

	// TTL gets the remaining time-to-live of a key.
	// key: The key of the cache item.
	// Returns: The remaining time-to-live of the key in seconds, or 0 if the key doesn't exist.
	TTL(key string) time.Duration

	// Set sets a cache item with an expiration timeout.
	// key: The key of the cache item.
	// val: The value of the cache item.
	// timeout: The expiration duration of the cache item.
	// Returns: An error if the set operation fails.
	Set(key string, val any, timeout time.Duration) error

	// SetNX sets a cache item only if the key does not already exist.
	// key: The key of the cache item.
	// val: The value of the cache item.
	// timeout: The expiration duration of the cache item.
	// Returns: True if the set operation succeeds, and false if the key already exists.
	SetNX(key string, val any, timeout time.Duration) (bool, error)

	// IsExist checks if a key exists in the cache.
	// key: The key of the cache item.
	// Returns: True if the key exists, otherwise false.
	IsExist(key string) bool

	// Delete removes a cache item by its key.
	// key: The key of the cache item.
	// Returns: An error if the deletion operation fails.
	Delete(key string) error
}
