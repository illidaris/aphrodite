package structure

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestKVs(t *testing.T) {
	t.Run("Set", func(t *testing.T) {
		t.Run("adds new key", func(t *testing.T) {
			kvs := NewKVs[string, int]()
			kvs.Set("a", 1)

			assert.Equal(t, 1, kvs.Len())
			assert.Equal(t, 1, kvs.Get("a"))
			assert.Len(t, kvs.ToSlice(), 1)
		})

		t.Run("ignores duplicate keys", func(t *testing.T) {
			kvs := NewKVs[int, string]()
			kvs.Set(1, "first")
			kvs.Set(1, "second")

			assert.Equal(t, 1, kvs.Len())
			assert.Equal(t, "first", kvs.Get(1))
		})
	})

	t.Run("Get", func(t *testing.T) {
		kvs := NewKVs[rune, float64]()
		kvs.Set('x', 3.14)

		t.Run("existing key", func(t *testing.T) {
			assert.Equal(t, 3.14, kvs.Get('x'))
		})

		t.Run("non-existent key panics", func(t *testing.T) {
			defer func() {
				if r := recover(); r == nil {
					t.Error("Expected panic for missing key")
				}
			}()
			kvs.Get('y') // This should panic
		})
	})

	t.Run("ToSlice", func(t *testing.T) {
		t.Run("returns insertion order", func(t *testing.T) {
			kvs := NewKVs[string, interface{}]()
			values := []*KV[string, interface{}]{
				{Key: "a", Value: 1},
				{Key: "b", Value: "two"},
				{Key: "c", Value: struct{}{}},
			}

			for _, v := range values {
				kvs.Set(v.Key, v.Value)
			}

			assert.Equal(t, values, kvs.ToSlice())
		})

		t.Run("empty collection", func(t *testing.T) {
			kvs := NewKVs[bool, bool]()
			assert.Empty(t, kvs.ToSlice())
		})
	})

	t.Run("Len", func(t *testing.T) {
		t.Run("empty", func(t *testing.T) {
			kvs := NewKVs[complex128, uint]()
			assert.Zero(t, kvs.Len())
		})

		t.Run("after multiple inserts", func(t *testing.T) {
			kvs := NewKVs[int, int]()
			for i := 0; i < 100; i++ {
				kvs.Set(i, i*2)
			}
			assert.Equal(t, 100, kvs.Len())
		})
	})
}
