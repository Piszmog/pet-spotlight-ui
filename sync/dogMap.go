package sync

import (
	"strings"
	"sync"
)

// DogMap is a thread-safe map of dogs.
type DogMap struct {
	m *sync.Map
}

// InitializeMap creates the map with the provided slice of dogs.
func InitializeMap(dogs []string) *DogMap {
	var m sync.Map
	for _, dog := range dogs {
		m.Store(strings.TrimSpace(strings.ToLower(dog)), false)
	}
	return &DogMap{m: &m}
}

// IsMatch determines if the provided name matches an entry in the map. The provided names must be contained an in a key.
func (m *DogMap) IsMatch(name string) bool {
	dogMatch := false
	m.m.Range(func(dog, alreadyDownloaded interface{}) bool {
		if !alreadyDownloaded.(bool) {
			if strings.Contains(name, dog.(string)) {
				dogMatch = true
				m.m.Store(dog, true)
				return false
			}
		}
		return true
	})
	return dogMatch
}

// IsCompete determines if all the entries in the map are true.
func (m *DogMap) IsCompete() bool {
	isComplete := true
	m.m.Range(func(_, found interface{}) bool {
		if !found.(bool) {
			isComplete = false
			return false
		}
		return true
	})
	return isComplete
}

// GetMissing returns all entries with a value of false.
func (m *DogMap) GetMissing() []string {
	var missing []string
	m.m.Range(func(name, found interface{}) bool {
		if !found.(bool) {
			missing = append(missing, name.(string))
		}
		return true
	})
	return missing
}
