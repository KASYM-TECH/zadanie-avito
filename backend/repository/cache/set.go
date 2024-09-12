package cache

import "sync"

type Set struct {
	store map[string]struct{}
	mutex *sync.RWMutex
}

func NewSet() *Set {
	return &Set{
		store: make(map[string]struct{}),
		mutex: &sync.RWMutex{},
	}
}

func (s *Set) Add(value string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.store[value] = struct{}{}
}

func (s *Set) Exists(value string) bool {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	_, ok := s.store[value]

	return ok
}

func (s *Set) Remove(value string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	delete(s.store, value)
}

func (s *Set) WarmUp(values []string) {
	if len(values) == 0 {
		return
	}
	s.mutex.Lock()
	defer s.mutex.Unlock()

	for _, value := range values {
		s.store[value] = struct{}{}
	}
}
