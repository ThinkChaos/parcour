package zync

import "sync"

// RWMutex is a type safe wrapper for `sync.RWMutex`.
//
// The API is slightly changed to make it less error prone:
//   - The protected value is stored by the mutex, thus the only way
//     to access it is via the mutex.
//   - A callback API is provided to ensure unlocking is done properly.
//     This is the recommended way of using the mutex.
//   - The lock methods return an `UnlockFunc` to make forgetting to unlock hard,
//     and guarantee the lock/unlock calls match.
type RWMutex[T any] struct {
	inner sync.RWMutex
	value T
}

// NewRWMutex creates and returns a new `RWMutex` with the given value.
func NewRWMutex[T any](value T) RWMutex[T] {
	return RWMutex[T]{value: value}
}

// WithRLock denotes a critical section with read access to the underlying value.
//
// The given function is called with `m` locked for reading.
// It is unsafe to leak the given value out of the scope of the function.
func (m *RWMutex[T]) WithRLock(f func(*T)) {
	valPtr, unlock := m.RLock()
	defer unlock()

	f(valPtr)
}

// RLock locks `m` for reading.
//
// The `With` functions should be preferred when possible.
//
// It should not be used for recursive read locking; a blocked Lock
// call excludes new readers from acquiring the lock. See the
// documentation on the `sync.RWMutex` type.
func (m *RWMutex[T]) RLock() (*T, UnlockFunc) {
	m.inner.RLock()

	return &m.value, m.inner.RUnlock
}

// TryRLock tries to lock `m` for reading and reports whether it succeeded.
// Both return values are `nil` if the attempt failed.
//
// The `WithRLock` function should be preferred when possible.
// Note that while correct uses of `TryRLock` do exist, they are rare,
// and use of `TryRLock` is often a sign of a deeper problem in a particular
// use of mutexes.
func (m *RWMutex[T]) TryRLock() (*T, UnlockFunc) {
	if !m.inner.TryRLock() {
		return nil, nil
	}

	return &m.value, m.inner.RUnlock
}

// WithWLock denotes a critical section with write access to the underlying value.
//
// The given function is called with the write lock being held.
// It is unsafe to leak the given value out of the scope of the function.
func (m *RWMutex[T]) WithWLock(f func(*T)) {
	valPtr, unlock := m.WLock()
	defer unlock()

	f(valPtr)
}

// WLock locks `m` for writing.
//
// The `With` functions should be preferred when possible.
//
// If the lock is already locked for reading or writing,
// Lock blocks until the lock is available.
func (m *RWMutex[T]) WLock() (*T, UnlockFunc) {
	m.inner.Lock()

	return &m.value, m.inner.Unlock
}

// TryWLock tries to lock `m` for writing and reports whether it succeeded.
// Both return values are `nil` if the attempt failed.
//
// The `WithWLock` function should be preferred when possible.
// Note that while correct uses of `TryWLock` do exist, they are rare,
// and use of `TryWLock` is often a sign of a deeper problem in a particular
// use of mutexes.
func (m *RWMutex[T]) TryWLock() (*T, UnlockFunc) {
	if !m.inner.TryLock() {
		return nil, nil
	}

	return &m.value, m.inner.RUnlock
}
