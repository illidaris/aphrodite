package structure

import "sync"

func NewItemMap[T any]() *ItemMap[T] {
	return &ItemMap[T]{
		mut: sync.RWMutex{},
		kv:  map[string]*T{},
	}
}

// ItemMap is a thread-safe key-value map for storing and retrieving items of any type T.
type ItemMap[T any] struct {
	mut sync.RWMutex  // Read-write mutex for concurrent access synchronization.
	kv  map[string]*T // The underlying map storing key-value pairs.
}

// GetOrSet attempts to retrieve an item by the given key. If the item does not exist,
// it uses the provided function `f` to generate and set the item.
// - key: The key for the item to get or set.
// - f: A function that generates the item when the key is not found, taking the key as an argument and returning the item and an error if any.
// Returns the item corresponding to the key, or nil if retrieval or setting fails.
func (i *ItemMap[T]) GetOrSet(key string, f func(key string) (*T, error)) *T {
	v, ok := i.GetItem(key) // Try to get the item.
	if ok {
		return v
	}
	v, err := f(key) // Generate the item using function `f`.
	if err != nil {
		return nil
	}
	i.SetItem(key, v)                  // Set the generated item.
	if res, ok := i.GetItem(key); ok { // Re-get to ensure successful setting.
		return res
	}
	return nil
}

// GetItem retrieves an item by its key.
// - key: The key of the item to retrieve.
// Returns the item pointer and a boolean indicating whether the item was successfully retrieved.
func (i *ItemMap[T]) GetItem(key string) (*T, bool) {
	i.mut.RLock() // Acquire read lock for safety.
	defer i.mut.RUnlock()
	v, ok := i.kv[key] // Attempt to get the item.
	return v, ok
}

// SetItem sets the item for the given key.
// - key: The key for the item to set.
// - value: The value of the item to set.
func (i *ItemMap[T]) SetItem(key string, value *T) {
	i.mut.Lock() // Acquire write lock for safety.
	defer i.mut.Unlock()
	if i.kv == nil {
		i.kv = map[string]*T{} // Initialize map if it's nil.
	}
	i.kv[key] = value // Set the item.
}
