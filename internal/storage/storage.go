package storage

import (
	"sync"
)

func NewMemoryStorage(initial map[string]string) *MemoryStorage {
	data := make(map[string]string)

	for key, value := range initial {
		data[key] = value
	}

	return &MemoryStorage{
		data: data,
	}
}

type MemoryStorage struct {
	data map[string]string
	mu   sync.RWMutex
}

func (m *MemoryStorage) Add(key string, value string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.data[key] = value
}

func (m *MemoryStorage) Get(key string) string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.data[key]
}
