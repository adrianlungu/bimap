package bimap

// ImmutableBiMap is a read-only bidirectional map. Safe for concurrent use
// without any locking — the data never changes after construction.
type ImmutableBiMap[K comparable, V comparable] struct {
	forward map[K]V
	inverse map[V]K
}

// NewImmutableBiMapFromMap builds an ImmutableBiMap from a map[K]V.
func NewImmutableBiMapFromMap[K, V comparable](forwardMap map[K]V) *ImmutableBiMap[K, V] {
	forward := make(map[K]V, len(forwardMap))
	inverse := make(map[V]K, len(forwardMap))
	for k, v := range forwardMap {
		forward[k] = v
		inverse[v] = k
	}
	return &ImmutableBiMap[K, V]{forward: forward, inverse: inverse}
}

// GetByKey returns the value for a given key and whether or not the element was present.
func (b *ImmutableBiMap[K, V]) GetByKey(k K) (V, bool) {
	v, ok := b.forward[k]
	return v, ok
}

// GetByValue returns the key for a given value and whether or not the element was present.
func (b *ImmutableBiMap[K, V]) GetByValue(v V) (K, bool) {
	k, ok := b.inverse[v]
	return k, ok
}

// GetByKeyWithFallback returns the value for k, or fallback if k is not present.
func (b *ImmutableBiMap[K, V]) GetByKeyWithFallback(k K, fallback V) V {
	if v, ok := b.forward[k]; ok {
		return v
	}
	return fallback
}

// GetByValueWithFallback returns the key for v, or fallback if v is not present.
func (b *ImmutableBiMap[K, V]) GetByValueWithFallback(v V, fallback K) K {
	if k, ok := b.inverse[v]; ok {
		return k
	}
	return fallback
}

// ExistsByKey checks whether or not a key exists in the ImmutableBiMap.
func (b *ImmutableBiMap[K, V]) ExistsByKey(k K) bool {
	_, ok := b.forward[k]
	return ok
}

// ExistsByValue checks whether or not a value exists in the ImmutableBiMap.
func (b *ImmutableBiMap[K, V]) ExistsByValue(v V) bool {
	_, ok := b.inverse[v]
	return ok
}

// Size returns the number of elements in the ImmutableBiMap.
func (b *ImmutableBiMap[K, V]) Size() int {
	return len(b.forward)
}

// GetForwardMap returns a copy of the forward map (key → value).
func (b *ImmutableBiMap[K, V]) GetForwardMap() map[K]V {
	m := make(map[K]V, len(b.forward))
	for k, v := range b.forward {
		m[k] = v
	}
	return m
}

// GetInverseMap returns a copy of the inverse map (value → key).
func (b *ImmutableBiMap[K, V]) GetInverseMap() map[V]K {
	m := make(map[V]K, len(b.inverse))
	for v, k := range b.inverse {
		m[v] = k
	}
	return m
}
