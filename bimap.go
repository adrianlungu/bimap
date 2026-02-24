// Package bimap provides a threadsafe bidirectional map
package bimap

import "sync"

// BiMap is a bi-directional hashmap that is thread safe and supports immutability
type BiMap[K comparable, V comparable] struct {
	s         sync.RWMutex
	immutable bool
	forward   map[K]V
	inverse   map[V]K
}

// NewBiMap returns a an empty, mutable, biMap
func NewBiMap[K comparable, V comparable]() *BiMap[K, V] {
	return &BiMap[K, V]{forward: make(map[K]V), inverse: make(map[V]K), immutable: false}
}

// NewBiMapFromMap returns a new BiMap from a map[K, V]
func NewBiMapFromMap[K comparable, V comparable](forwardMap map[K]V) *BiMap[K, V] {
	biMap := NewBiMap[K, V]()
	for k, v := range forwardMap {
		biMap.Insert(k, v)
	}
	return biMap
}

// Insert puts a key and value into the BiMap, provided its mutable. Also creates the reverse mapping from value to key.
func (b *BiMap[K, V]) Insert(k K, v V) {
	b.s.Lock()
	defer b.s.Unlock()
	if b.immutable {
		panic("Cannot modify immutable map")
	}
	if _, ok := b.forward[k]; ok {
		delete(b.inverse, b.forward[k])
	}
	b.forward[k] = v
	b.inverse[v] = k
}

// ExistsByKey checks whether or not a key exists in the BiMap.
func (b *BiMap[K, V]) ExistsByKey(k K) bool {
	b.s.RLock()
	defer b.s.RUnlock()
	_, ok := b.forward[k]
	return ok
}

// Exists checks whether or not a key exists in the BiMap.
//
// Deprecated: Use ExistsByKey instead.
func (b *BiMap[K, V]) Exists(k K) bool { return b.ExistsByKey(k) }

// ExistsByValue checks whether or not a value exists in the BiMap.
func (b *BiMap[K, V]) ExistsByValue(k V) bool {
	b.s.RLock()
	defer b.s.RUnlock()
	_, ok := b.inverse[k]
	return ok
}

// ExistsInverse checks whether or not a value exists in the BiMap.
//
// Deprecated: Use ExistsByValue instead.
func (b *BiMap[K, V]) ExistsInverse(k V) bool { return b.ExistsByValue(k) }

// GetByKey returns the value for a given key in the BiMap and whether or not the element was present.
func (b *BiMap[K, V]) GetByKey(k K) (V, bool) {
	b.s.RLock()
	defer b.s.RUnlock()
	v, ok := b.forward[k]
	return v, ok
}

// Get returns the value for a given key in the BiMap and whether or not the element was present.
//
// Deprecated: Use GetByKey instead.
func (b *BiMap[K, V]) Get(k K) (V, bool) { return b.GetByKey(k) }

// GetByValue returns the key for a given value in the BiMap and whether or not the element was present.
func (b *BiMap[K, V]) GetByValue(v V) (K, bool) {
	b.s.RLock()
	defer b.s.RUnlock()
	k, ok := b.inverse[v]
	return k, ok
}

// GetInverse returns the key for a given value in the BiMap and whether or not the element was present.
//
// Deprecated: Use GetByValue instead.
func (b *BiMap[K, V]) GetInverse(v V) (K, bool) { return b.GetByValue(v) }

// GetByKeyWithFallback returns the value for k, or fallback if k is not present.
func (b *BiMap[K, V]) GetByKeyWithFallback(k K, fallback V) V {
	if v, ok := b.GetByKey(k); ok {
		return v
	}
	return fallback
}

// GetByValueWithFallback returns the key for v, or fallback if v is not present.
func (b *BiMap[K, V]) GetByValueWithFallback(v V, fallback K) K {
	if k, ok := b.GetByValue(v); ok {
		return k
	}
	return fallback
}

// DeleteByKey removes a key-value pair from the BiMap for a given key. Returns if the key doesn't exist.
func (b *BiMap[K, V]) DeleteByKey(k K) {
	b.s.Lock()
	defer b.s.Unlock()
	if b.immutable {
		panic("Cannot modify immutable map")
	}
	val, ok := b.forward[k]
	if !ok {
		return
	}
	delete(b.forward, k)
	delete(b.inverse, val)
}

// Delete removes a key-value pair from the BiMap for a given key. Returns if the key doesn't exist.
//
// Deprecated: Use DeleteByKey instead.
func (b *BiMap[K, V]) Delete(k K) { b.DeleteByKey(k) }

// DeleteByValue removes a key-value pair from the BiMap for a given value. Returns if the value doesn't exist.
func (b *BiMap[K, V]) DeleteByValue(v V) {
	b.s.Lock()
	defer b.s.Unlock()
	if b.immutable {
		panic("Cannot modify immutable map")
	}
	key, ok := b.inverse[v]
	if !ok {
		return
	}
	delete(b.inverse, v)
	delete(b.forward, key)
}

// DeleteInverse removes a key-value pair from the BiMap for a given value. Returns if the value doesn't exist.
//
// Deprecated: Use DeleteByValue instead.
func (b *BiMap[K, V]) DeleteInverse(v V) { b.DeleteByValue(v) }

// Size returns the number of elements in the bimap
func (b *BiMap[K, V]) Size() int {
	b.s.RLock()
	defer b.s.RUnlock()
	return len(b.forward)
}

// MakeImmutable freezes the BiMap preventing any further write actions from taking place
func (b *BiMap[K, V]) MakeImmutable() {
	b.s.Lock()
	defer b.s.Unlock()
	b.immutable = true
}

// Freeze returns a new ImmutableBiMap with a snapshot of the current state.
// The original BiMap is unaffected and remains mutable.
func (b *BiMap[K, V]) Freeze() *ImmutableBiMap[K, V] {
	b.s.RLock()
	defer b.s.RUnlock()
	forward := make(map[K]V, len(b.forward))
	inverse := make(map[V]K, len(b.inverse))
	for k, v := range b.forward {
		forward[k] = v
	}
	for v, k := range b.inverse {
		inverse[v] = k
	}
	return &ImmutableBiMap[K, V]{forward: forward, inverse: inverse}
}

// GetInverseMap returns a regular go map mapping from the BiMap's values to its keys
func (b *BiMap[K, V]) GetInverseMap() map[V]K {
	return b.inverse
}

// GetForwardMap returns a regular go map mapping from the BiMap's keys to its values
func (b *BiMap[K, V]) GetForwardMap() map[K]V {
	return b.forward
}

// Lock manually locks the BiMap's mutex
func (b *BiMap[K, V]) Lock() {
	b.s.Lock()
}

// Unlock manually unlocks the BiMap's mutex
func (b *BiMap[K, V]) Unlock() {
	b.s.Unlock()
}
