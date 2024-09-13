package cache

import "sync"

type KeyValue struct {
	Key   string
	Value string
}

type Storage struct {
	store map[string]string
	mutex *sync.RWMutex
}

func NewStorage() *Storage {
	return &Storage{
		store: make(map[string]string),
		mutex: &sync.RWMutex{},
	}
}

func (s *Storage) Add(key, value string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.store[key] = value
}

func (s *Storage) Get(key string) (string, bool) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	val, ok := s.store[key]
	
	return val, ok
}

func (s *Storage) Remove(key string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	delete(s.store, key)
}

func (s *Storage) WarmUp(items []KeyValue) {
	if len(items) == 0 {
		return
	}
	s.mutex.Lock()
	defer s.mutex.Unlock()

	for _, item := range items {
		s.store[item.Key] = item.Value
	}
}
