package zync

import "sync"

// Map is a type safe wrapper for `sync.Map`.
//
// The `sync.Map` is re-used with `any` replaced by concrete types.
//
// An additional `TryRange` method is provided for faillible iteration.
type Map[K comparable, V any] struct {
	inner sync.Map
}

// Load returns the value stored in the map for a key, or nil if no
// value is present.
// The ok result indicates whether value was found in the map.
func (m *Map[K, V]) Load(key K) (value V, ok bool) {
	val, ok := m.inner.Load(key)

	return val.(V), ok //nolint:forcetypeassert
}

// Store sets the value for a key.
func (m *Map[K, V]) Store(key K, value V) {
	m.inner.Store(key, value)
}

// Delete deletes the value for a key.
func (m *Map[K, V]) Delete(key K) {
	m.inner.Delete(key)
}

// LoadOrStore returns the existing value for the key if present.
// Otherwise, it stores and returns the given value.
// The loaded result is true if the value was loaded, false if stored.
func (m *Map[K, V]) LoadOrStore(key K, value V) (actual V, loaded bool) {
	val, ok := m.inner.LoadOrStore(key, value)

	return val.(V), ok //nolint:forcetypeassert
}

// LoadAndDelete deletes the value for a key, returning the previous value if any.
// The loaded result reports whether the key was present.
func (m *Map[K, V]) LoadAndDelete(key K) (value V, loaded bool) {
	val, ok := m.inner.LoadAndDelete(key)

	return val.(V), ok //nolint:forcetypeassert
}

// Range calls f sequentially for each key and value present in the map.
// If f returns false, range stops the iteration.
//
// Range does not necessarily correspond to any consistent snapshot of the Map's
// contents: no key will be visited more than once, but if the value for any key
// is stored or deleted concurrently (including by f), Range may reflect any
// mapping for that key from any point during the Range call. Range does not
// block other methods on the receiver; even f itself may call any method on m.
//
// Range may be O(N) with the number of elements in the map even if f returns
// false after a constant number of calls.
func (m *Map[K, V]) Range(f func(key K, value V) bool) {
	m.inner.Range(func(key, value any) bool {
		return f(key.(K), value.(V)) //nolint:forcetypeassert
	})
}

// TryRange is like `Range`, but with a faillible callback.
//
// If `f` returns an error, iteration stops.
func (m *Map[K, V]) TryRange(f func(key K, value V) error) error {
	var err error

	m.inner.Range(func(key, value any) bool {
		err = f(key.(K), value.(V)) //nolint:forcetypeassert

		return err == nil
	})

	return err
}
