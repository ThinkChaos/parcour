package zync

import "sync"

// Mutex is a type safe wrapper for `sync.Mutex`.
//
// The API is slightly changed to make it less error prone:
//   - The protected value is stored by the mutex, thus the only way
//     to access it is via the mutex.
//   - A callback API is provided to ensure unlocking is done properly.
//     This is the recommended way of using the mutex.
//   - The lock methods return an `UnlockFunc` to make forgetting to unlock hard.
type Mutex[T any] struct {
	inner sync.Mutex
	value T
}

// NewMutex creates and returns a new `Mutex` with the given value.
func NewMutex[T any](value T) Mutex[T] {
	return Mutex[T]{value: value}
}

// WithLock denotes a critical section with access to the underlying value.
//
// The given function is called with the lock being held.
// It is unsafe to leak the given value out of the scope of the function.
func (m *Mutex[T]) WithLock(f func(*T)) {
	valPtr, unlock := m.Lock()
	defer unlock()

	f(valPtr)
}

// Lock blocks until the mutex is available.
//
// The `WithLock` function should be preferred when possible.
func (m *Mutex[T]) Lock() (*T, UnlockFunc) {
	m.inner.Lock()

	return &m.value, m.inner.Unlock
}

// TryLock tries to lock `m` and reports whether it succeeded.
// Both return values are `nil` if the attempt failed.
//
// The `WithLock` function should be preferred when possible.
// Note that while correct uses of `TryLock` do exist, they are rare,
// and use of TryLock is often a sign of a deeper problem in a particular
// use of mutexes.
func (m *Mutex[T]) TryLock() (*T, UnlockFunc) {
	if !m.inner.TryLock() {
		return nil, nil
	}

	return &m.value, m.inner.Unlock
}
