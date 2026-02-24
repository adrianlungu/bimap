package bimap

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

const key = "key"
const value = "value"

func TestNewBiMap(t *testing.T) {
	actual := NewBiMap[string, string]()
	expected := &BiMap[string, string]{forward: make(map[string]string), inverse: make(map[string]string)}
	assert.Equal(t, expected, actual, "They should be equal")
}

func TestNewBiMapFrom(t *testing.T) {
	actual := NewBiMapFromMap(map[string]string{
		key: value,
	})
	actual.Insert(key, value)

	fwdExpected := make(map[string]string)
	invExpected := make(map[string]string)
	fwdExpected[key] = value
	invExpected[value] = key
	expected := &BiMap[string, string]{forward: fwdExpected, inverse: invExpected}

	assert.Equal(t, expected, actual, "They should be equal")
}

func TestBiMap_Insert(t *testing.T) {
	actual := NewBiMap[string, string]()
	actual.Insert(key, value)

	fwdExpected := make(map[string]string)
	invExpected := make(map[string]string)
	fwdExpected[key] = value
	invExpected[value] = key
	expected := &BiMap[string, string]{forward: fwdExpected, inverse: invExpected}

	assert.Equal(t, expected, actual, "They should be equal")
}

func TestBiMap_InsertTwice(t *testing.T) {
	additionalValue := value + value

	actual := NewBiMap[string, string]()
	actual.Insert(key, value)
	actual.Insert(key, additionalValue)

	fwdExpected := make(map[string]string)
	invExpected := make(map[string]string)
	fwdExpected[key] = additionalValue

	invExpected[additionalValue] = key
	expected := &BiMap[string, string]{forward: fwdExpected, inverse: invExpected}

	assert.Equal(t, expected, actual, "They should be equal")
}

func TestBiMap_Exists(t *testing.T) {
	actual := NewBiMap[string, string]()

	actual.Insert(key, value)
	assert.False(t, actual.ExistsByKey("ARBITARY_KEY"), "Key should not exist")
	assert.True(t, actual.ExistsByKey(key), "Inserted key should exist")
}

func TestBiMap_InverseExists(t *testing.T) {
	actual := NewBiMap[string, string]()

	actual.Insert(key, value)
	assert.False(t, actual.ExistsByValue("ARBITARY_VALUE"), "Value should not exist")
	assert.True(t, actual.ExistsByValue(value), "Inserted value should exist")
}

func TestBiMap_Get(t *testing.T) {
	actual := NewBiMap[string, string]()

	actual.Insert(key, value)

	actualVal, ok := actual.GetByKey(key)

	assert.True(t, ok, "It should return true")
	assert.Equal(t, value, actualVal, "Value and returned val should be equal")

	actualVal, ok = actual.GetByKey(value)

	assert.False(t, ok, "It should return false")
	assert.Empty(t, actualVal, "Actual val should be empty")
}

func TestBiMap_GetInverse(t *testing.T) {
	actual := NewBiMap[string, string]()

	actual.Insert(key, value)

	actualKey, ok := actual.GetByValue(value)

	assert.True(t, ok, "It should return true")
	assert.Equal(t, key, actualKey, "Key and returned key should be equal")

	actualKey, ok = actual.GetByKey(value)

	assert.False(t, ok, "It should return false")
	assert.Empty(t, actualKey, "Actual key should be empty")
}

func TestBiMap_GetByKeyWithFallback(t *testing.T) {
	actual := NewBiMap[string, string]()
	actual.Insert(key, value)

	assert.Equal(t, value, actual.GetByKeyWithFallback(key, "fallback"), "Should return value for present key")
	assert.Equal(t, "fallback", actual.GetByKeyWithFallback("missing", "fallback"), "Should return fallback for absent key")
}

func TestBiMap_GetByValueWithFallback(t *testing.T) {
	actual := NewBiMap[string, string]()
	actual.Insert(key, value)

	assert.Equal(t, key, actual.GetByValueWithFallback(value, "fallback"), "Should return key for present value")
	assert.Equal(t, "fallback", actual.GetByValueWithFallback("missing", "fallback"), "Should return fallback for absent value")
}

func TestBiMap_Size(t *testing.T) {
	actual := NewBiMap[string, string]()

	assert.Equal(t, 0, actual.Size(), "Length of empty bimap should be zero")

	actual.Insert(key, value)

	assert.Equal(t, 1, actual.Size(), "Length of bimap should be one")
}

func TestBiMap_Delete(t *testing.T) {
	actual := NewBiMap[string, string]()
	dummyKey := "DummyKey"
	dummyVal := "DummyVal"
	actual.Insert(key, value)
	actual.Insert(dummyKey, dummyVal)

	assert.Equal(t, 2, actual.Size(), "Size of bimap should be two")

	actual.DeleteByKey(dummyKey)

	fwdExpected := make(map[string]string)
	invExpected := make(map[string]string)
	fwdExpected[key] = value
	invExpected[value] = key

	expected := &BiMap[string, string]{forward: fwdExpected, inverse: invExpected}

	assert.Equal(t, 1, actual.Size(), "Size of bimap should be two")
	assert.Equal(t, expected, actual, "They should be the same")

	actual.DeleteByKey(dummyKey)

	assert.Equal(t, 1, actual.Size(), "Size of bimap should be two")
	assert.Equal(t, expected, actual, "They should be the same")
}

func TestBiMap_InverseDelete(t *testing.T) {
	actual := NewBiMap[string, string]()
	dummyKey := "DummyKey"
	dummyVal := "DummyVal"
	actual.Insert(key, value)
	actual.Insert(dummyKey, dummyVal)

	assert.Equal(t, 2, actual.Size(), "Size of bimap should be two")

	actual.DeleteByValue(dummyVal)

	fwdExpected := make(map[string]string)
	invExpected := make(map[string]string)
	fwdExpected[key] = value
	invExpected[value] = key

	expected := &BiMap[string, string]{forward: fwdExpected, inverse: invExpected}

	assert.Equal(t, 1, actual.Size(), "Size of bimap should be two")
	assert.Equal(t, expected, actual, "They should be the same")

	actual.DeleteByValue(dummyVal)

	assert.Equal(t, 1, actual.Size(), "Size of bimap should be two")
	assert.Equal(t, expected, actual, "They should be the same")
}

func TestBiMap_WithVaryingType(t *testing.T) {
	actual := NewBiMap[string, int]()
	dummyKey := "Dummy key"
	dummyVal := 3

	actual.Insert(dummyKey, dummyVal)

	res, _ := actual.GetByKey(dummyKey)
	resVal, _ := actual.GetByValue(dummyVal)
	assert.Equal(t, dummyVal, res, "Get by string key should return integer val")
	assert.Equal(t, dummyKey, resVal, "Get by integer val should return string key")

}

func TestBiMap_MakeImmutable(t *testing.T) {
	actual := NewBiMap[string, int]()
	dummyKey := "Dummy key"
	dummyVal := 3

	actual.Insert(dummyKey, dummyVal)

	actual.MakeImmutable()

	assert.Panics(t, func() {
		actual.DeleteByKey(dummyKey)
	}, "It should panic on a mutation operation")

	val, _ := actual.GetByKey(dummyKey)

	assert.Equal(t, dummyVal, val, "It should still have the value")

	assert.Panics(t, func() {
		actual.DeleteByValue(dummyVal)
	}, "It should panic on a mutation operation")

	k, _ := actual.GetByValue(dummyVal)

	assert.Equal(t, dummyKey, k, "It should still have the key")

	size := actual.Size()

	assert.Equal(t, 1, size, "Size should be one")

	assert.Panics(t, func() {
		actual.Insert("New", 1)
	}, "It should panic on a mutation operation")

	size = actual.Size()

	assert.Equal(t, 1, size, "Size should be one")

}

func TestBiMap_GetForwardMap(t *testing.T) {
	actual := NewBiMap[string, int]()
	dummyKey := "Dummy key"
	dummyVal := 42

	forwardMap := make(map[string]int)
	forwardMap[dummyKey] = dummyVal

	actual.Insert(dummyKey, dummyVal)

	actualForwardMap := actual.GetForwardMap()
	eq := reflect.DeepEqual(actualForwardMap, forwardMap)
	assert.True(t, eq, "Forward maps should be equal")
}

func TestBiMap_GetInverseMap(t *testing.T) {
	actual := NewBiMap[string, int]()
	dummyKey := "Dummy key"
	dummyVal := 42

	inverseMap := make(map[int]string)
	inverseMap[dummyVal] = dummyKey

	actual.Insert(dummyKey, dummyVal)

	actualInverseMap := actual.GetInverseMap()
	eq := reflect.DeepEqual(actualInverseMap, inverseMap)
	assert.True(t, eq, "Inverse maps should be equal")
}

func TestBiMap_Freeze(t *testing.T) {
	mutable := NewBiMap[string, int]()
	mutable.Insert("a", 1)
	mutable.Insert("b", 2)

	frozen := mutable.Freeze()

	assert.Equal(t, 2, frozen.Size(), "Frozen map should have same size")

	v, ok := frozen.GetByKey("a")
	assert.True(t, ok)
	assert.Equal(t, 1, v)

	k, ok := frozen.GetByValue(2)
	assert.True(t, ok)
	assert.Equal(t, "b", k)

	// Mutations to original do not affect the frozen copy
	mutable.Insert("c", 3)
	assert.Equal(t, 2, frozen.Size(), "Frozen map should be unaffected by mutations to the original")

	// Frozen map does not affect the original
	assert.Equal(t, 3, mutable.Size(), "Original should still have all insertions")
}
