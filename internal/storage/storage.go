package storage

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
}

func (m *MemoryStorage) Add(key string, value string) {
	m.data[key] = value
}

func (m *MemoryStorage) Get(key string) string {
	return m.data[key]
}
