package structure

type KV[K comparable, V any] struct {
	Key   K
	Value V
}

func NewKVs[K comparable, V any]() *KVs[K, V] {
	return &KVs[K, V]{
		s: make([]*KV[K, V], 0),
		m: make(map[K]*KV[K, V]),
	}
}

type KVs[K comparable, V any] struct {
	s []*KV[K, V]
	m map[K]*KV[K, V]
}

func (n *KVs[K, V]) Set(k K, v V) {
	if _, ok := n.m[k]; !ok {
		pair := &KV[K, V]{Key: k, Value: v}
		n.s = append(n.s, pair)
		n.m[k] = pair
	}
}

func (n KVs[K, V]) Get(k K) V {
	return n.m[k].Value
}

func (n KVs[K, V]) ToSlice() []*KV[K, V] {
	return n.s
}

func (n KVs[K, V]) Len() int {
	return len(n.s)
}
