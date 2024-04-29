package structure

import (
	"fmt"
	"sync"
	"testing"
)

func TestItemMap_GetOrSet(t *testing.T) {
	im := ItemMap[string]{
		mut: &sync.RWMutex{},
		kv:  map[string]*string{},
	}

	// Test case 1: Item exists
	key1 := "key1"
	value1 := "value1"
	im.SetItem(key1, &value1)

	got := im.GetOrSet(key1, func(key string) (*string, error) {
		return nil, nil
	})

	if got == nil || *got != value1 {
		t.Errorf("Expected %s, got %v", value1, got)
	}

	// Test case 2: Item does not exist, generation successful
	key2 := "key2"
	got = im.GetOrSet(key2, func(key string) (*string, error) {
		return &key2, nil
	})

	if got == nil || *got != key2 {
		t.Errorf("Expected %s, got %v", key2, got)
	}

	// Test case 3: Item does not exist, generation failed
	key3 := "key3"
	got = im.GetOrSet(key3, func(key string) (*string, error) {
		return nil, fmt.Errorf("error generating item")
	})

	if got != nil {
		t.Errorf("Expected nil, got %v", got)
	}
}

func TestItemMap_GetItem(t *testing.T) {
	im := ItemMap[string]{
		mut: &sync.RWMutex{},
		kv:  map[string]*string{},
	}

	// Test case 1: Item exists
	key1 := "key1"
	value1 := "value1"
	im.SetItem(key1, &value1)

	got, ok := im.GetItem(key1)

	if !ok || got == nil || *got != value1 {
		t.Errorf("Expected (%s, true), got (%v, %v)", value1, got, ok)
	}

	// Test case 2: Item does not exist
	key2 := "key2"

	got, ok = im.GetItem(key2)

	if ok || got != nil {
		t.Errorf("Expected (nil, false), got (%v, %v)", got, ok)
	}
}

func TestItemMap_SetItem(t *testing.T) {
	im := ItemMap[string]{
		mut: &sync.RWMutex{},
		kv:  map[string]*string{},
	}

	// Test case: Set item
	key := "key"
	value := "value"
	im.SetItem(key, &value)

	got, ok := im.GetItem(key)

	if !ok || got == nil || *got != value {
		t.Errorf("Expected (%s, true), got (%v, %v)", value, got, ok)
	}
}
