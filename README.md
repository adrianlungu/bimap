# bimap
[![GoDoc](https://godoc.org/github.com/adrianlungu/bimap?status.svg)](https://godoc.org/github.com/adrianlungu/bimap)

Generic bidirectional maps written in Go

Based on vishalkuo's [bimap](https://github.com/vishalkuo/bimap).

## Installation
```
go get github.com/adrianlungu/bimap
```

## Usage

### Mutable BiMap

```go
import "github.com/adrianlungu/bimap"

// Create an empty BiMap
b := bimap.NewBiMap[string, int]()
b.Insert("apples", 1)
b.Insert("bananas", 2)

// Look up by key or by value
val, ok := b.GetByKey("apples")    // 1, true
key, ok := b.GetByValue(2)         // "bananas", true

// Single-line retrieval with a fallback default
val = b.GetByKeyWithFallback("missing", -1)    // -1
key = b.GetByValueWithFallback(99, "unknown")  // "unknown"

// Check existence
b.ExistsByKey("apples")   // true
b.ExistsByValue(99)       // false

// Delete
b.DeleteByKey("apples")
b.DeleteByValue(2)

b.Size() // 0

// Build from an existing map
b2 := bimap.NewBiMapFromMap(map[string]int{"a": 1, "b": 2})
```

### Immutable BiMap

`ImmutableBiMap` is a read-only snapshot — no mutating methods, no mutex overhead, safe for concurrent use.

```go
// Build directly from a map
ib := bimap.NewImmutableBiMapFromMap(map[string]int{"a": 1, "b": 2})

// Freeze a mutable BiMap into an immutable snapshot
b := bimap.NewBiMap[string, int]()
b.Insert("x", 10)
ib = b.Freeze() // b is unaffected and remains mutable

val, ok := ib.GetByKey("x")               // 10, true
key, ok := ib.GetByValue(10)              // "x", true
val = ib.GetByKeyWithFallback("y", -1)    // -1
ib.ExistsByKey("x")                       // true
ib.Size()                                 // 1

// GetForwardMap / GetInverseMap return copies (mutations do not affect ib)
fwd := ib.GetForwardMap() // map[string]int{"x": 10}
inv := ib.GetInverseMap() // map[int]string{10: "x"}
```

### Thread safety

`BiMap` uses a `sync.RWMutex` internally. Use `Lock`/`Unlock` if you need to hold the mutex across multiple operations.

`ImmutableBiMap` requires no locking — its data never changes after construction.

### MakeImmutable

`BiMap` also supports in-place freezing via `MakeImmutable()`. After this call the map panics on any write attempt. Use `Freeze()` instead when you want a separate immutable copy while keeping the original mutable.

```go
b := bimap.NewBiMap[string, int]()
b.Insert("x", 1)
b.MakeImmutable()
b.Insert("y", 2) // panics
```

### Deprecated methods

The following methods are kept for backward compatibility but should be replaced with their `ByKey`/`ByValue` equivalents:

| Deprecated | Replacement |
|---|---|
| `Get` | `GetByKey` |
| `GetInverse` | `GetByValue` |
| `Exists` | `ExistsByKey` |
| `ExistsInverse` | `ExistsByValue` |
| `Delete` | `DeleteByKey` |
| `DeleteInverse` | `DeleteByValue` |
