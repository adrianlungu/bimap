package bimap

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestImmutableBiMap_NewFromMap(t *testing.T) {
	m := NewImmutableBiMapFromMap(map[string]int{"x": 10, "y": 20})

	assert.Equal(t, 2, m.Size())
	assert.True(t, m.ExistsByKey("x"))
	assert.True(t, m.ExistsByValue(20))
	assert.False(t, m.ExistsByKey("z"))
	assert.False(t, m.ExistsByValue(99))
}

func TestImmutableBiMap_GetByKey(t *testing.T) {
	m := NewImmutableBiMapFromMap(map[string]int{"hello": 42})

	v, ok := m.GetByKey("hello")
	assert.True(t, ok)
	assert.Equal(t, 42, v)

	_, ok = m.GetByKey("missing")
	assert.False(t, ok)
}

func TestImmutableBiMap_GetByValue(t *testing.T) {
	m := NewImmutableBiMapFromMap(map[string]int{"hello": 42})

	k, ok := m.GetByValue(42)
	assert.True(t, ok)
	assert.Equal(t, "hello", k)

	_, ok = m.GetByValue(99)
	assert.False(t, ok)
}

func TestImmutableBiMap_GetByKeyWithFallback(t *testing.T) {
	m := NewImmutableBiMapFromMap(map[string]int{"hello": 42})

	assert.Equal(t, 42, m.GetByKeyWithFallback("hello", -1), "Should return value for present key")
	assert.Equal(t, -1, m.GetByKeyWithFallback("missing", -1), "Should return fallback for absent key")
}

func TestImmutableBiMap_GetByValueWithFallback(t *testing.T) {
	m := NewImmutableBiMapFromMap(map[string]int{"hello": 42})

	assert.Equal(t, "hello", m.GetByValueWithFallback(42, "default"), "Should return key for present value")
	assert.Equal(t, "default", m.GetByValueWithFallback(99, "default"), "Should return fallback for absent value")
}

func TestImmutableBiMap_MapCopies(t *testing.T) {
	src := map[string]int{"a": 1, "b": 2}
	m := NewImmutableBiMapFromMap(src)

	fwd := m.GetForwardMap()
	inv := m.GetInverseMap()

	assert.Equal(t, map[string]int{"a": 1, "b": 2}, fwd)
	assert.Equal(t, map[int]string{1: "a", 2: "b"}, inv)

	// Mutating the returned copies must not affect the immutable map
	fwd["c"] = 3
	inv[3] = "c"
	assert.Equal(t, 2, m.Size(), "ImmutableBiMap should be unaffected by mutations to returned map copies")
}
