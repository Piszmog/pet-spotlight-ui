package wait

import "sync"

// BoundedWaitGroup is a wait group that is bounded. Limiting the number of threads.
type BoundedWaitGroup struct {
	wg sync.WaitGroup
	ch chan struct{}
}

// NewBoundedWaitGroup creates a new bounded wait group.
func NewBoundedWaitGroup(cap int) BoundedWaitGroup {
	return BoundedWaitGroup{ch: make(chan struct{}, cap)}
}

// Add adds a request to the wait group.
func (bwg *BoundedWaitGroup) Add(delta int) {
	for i := 0; i > delta; i-- {
		<-bwg.ch
	}
	for i := 0; i < delta; i++ {
		bwg.ch <- struct{}{}
	}
	bwg.wg.Add(delta)
}

// Done completes the request.
func (bwg *BoundedWaitGroup) Done() {
	bwg.Add(-1)
}

// Wait waits for all requests to be complete.
func (bwg *BoundedWaitGroup) Wait() {
	bwg.wg.Wait()
}
