package sync

import "sync"

// DogList is a thread-safe slice of dogs.
type DogList struct {
	m    sync.RWMutex
	dogs []string
}

// Add adds the dog to the list.
func (l *DogList) Add(name string) {
	l.m.Lock()
	l.dogs = append(l.dogs, name)
	l.m.Unlock()
}

// Get retrieves all the dogs.
func (l *DogList) Get() []string {
	l.m.RLock()
	l.m.RUnlock()
	return l.dogs
}
