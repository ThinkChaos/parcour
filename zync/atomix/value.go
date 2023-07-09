package atomix

import "sync/atomic"

// Value is the generic counterpart of `atomic.Value`.
//
// There are two key differences to the `atomic` type:
//   - a `Value` can safely be copied
//   - the zero value is invalid and a constructor must be used.
//     This avoids having to deal with the zero value of `atomic.Value`
//     which contains a "nil interface": that is not of type `T`.
//
// Internally, we use a `*atomic.Value`, thus if you get a `nil` deref panic,
// it is because you tried to use an uninitialized `Value`.
type Value[T any] struct {
	// Use a pointer so we can ensure the zero value for `Value` is invalid.
	// This also has the benefit of making it safe to copy a `Value`.
	inner *atomic.Value
}

func NewValue[T any](val T) Value[T] {
	v := Value[T]{new(atomic.Value)}

	v.Store(val)

	return v
}

func NewZeroValue[T any]() Value[T] {
	var zero T

	return NewValue(zero)
}

func (v Value[T]) CompareAndSwap(old, new T) (swapped bool) {
	return v.inner.CompareAndSwap(old, new)
}

func (v Value[T]) Load() T {
	return v.inner.Load().(T)
}

func (v Value[T]) Store(val T) {
	v.inner.Store(val)
}

func (v Value[T]) Swap(new T) (old T) {
	return v.inner.Swap(new).(T)
}
