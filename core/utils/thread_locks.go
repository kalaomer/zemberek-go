package utils

import "sync"

// ReadWriteLock is a lock object that allows many simultaneous "read locks",
// but only one "write lock."
type ReadWriteLock struct {
	mu      sync.Mutex
	cond    *sync.Cond
	readers int
}

// NewReadWriteLock creates a new ReadWriteLock
func NewReadWriteLock() *ReadWriteLock {
	rwl := &ReadWriteLock{}
	rwl.cond = sync.NewCond(&rwl.mu)
	return rwl
}

// AcquireRead acquires a read lock. Blocks only if a thread has acquired the write lock.
func (rwl *ReadWriteLock) AcquireRead() {
	rwl.mu.Lock()
	rwl.readers++
	rwl.mu.Unlock()
}

// ReleaseRead releases a read lock
func (rwl *ReadWriteLock) ReleaseRead() {
	rwl.mu.Lock()
	rwl.readers--
	if rwl.readers == 0 {
		rwl.cond.Broadcast()
	}
	rwl.mu.Unlock()
}

// AcquireWrite acquires a write lock. Blocks until there are no acquired read or write locks.
func (rwl *ReadWriteLock) AcquireWrite() {
	rwl.mu.Lock()
	for rwl.readers > 0 {
		rwl.cond.Wait()
	}
}

// ReleaseWrite releases a write lock
func (rwl *ReadWriteLock) ReleaseWrite() {
	rwl.mu.Unlock()
}
