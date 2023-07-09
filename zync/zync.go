// Package zync provides generic and safer wrappers for types from the standard sync module.
package zync

// UnlockFunc unlocks the mutex which returned it.
//
// See the documentation for Mutex and RWMutex for the reasoning behind this design.
type UnlockFunc func()
