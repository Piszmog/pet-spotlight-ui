package sync

import "sync"

// AtomicBoolean is a thread-safe boolean.
type AtomicBoolean struct {
	t bool
	m sync.RWMutex
}

// Get retrieves the value of the boolean.
func (p *AtomicBoolean) Get() bool {
	p.m.RLock()
	p.m.RUnlock()
	return p.t
}

// Set sets the value of the boolean.
func (p *AtomicBoolean) Set(value bool) {
	p.m.Lock()
	p.t = value
	p.m.Unlock()
}
