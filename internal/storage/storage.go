package storage

import (
	"bufio"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"os"
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

func (m *MemoryStorage) Add(url string, short string) (err error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.data[short] = url
	return nil
}

func (m *MemoryStorage) Get(key string) (string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.data[key], nil
}

type FileStorage struct {
	path string
	file *os.File
	mu   sync.RWMutex
}

type storageItem struct {
	UUID        uuid.UUID `json:"uuid"`
	OriginalURL string    `json:"original_url"`
	ShortURL    string    `json:"short_url"`
}

func NewFileStorage(path string) (*FileStorage, error) {
	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

	if err != nil {
		return nil, err
	}

	return &FileStorage{
		path: path,
		file: file,
	}, nil
}

func (m *FileStorage) Add(url string, short string) (err error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	id := uuid.New()
	storageItem := &storageItem{
		UUID:        id,
		ShortURL:    short,
		OriginalURL: url,
	}

	data, err := json.Marshal(&storageItem)
	if err != nil {
		return err
	}
	// добавляем перенос строки
	data = append(data, '\n')

	_, err = m.file.Write(data)
	return err
}

func (m *FileStorage) Get(short string) (string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	m.file.Seek(0, 0)
	result := make(map[string]storageItem)
	scanner := bufio.NewScanner(m.file)
	fmt.Println("1")

	info, _ := m.file.Stat()
	fmt.Println("Размер файла:", info.Size())

	for scanner.Scan() {
		line := scanner.Text()
		item := storageItem{}
		fmt.Println("2", line)
		err := json.Unmarshal([]byte(line), &item)
		if err != nil {
			fmt.Println("Ошибка парсинга JSON:", err)
			continue
		}
		result[item.ShortURL] = item
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Ошибка сканирования файла:", err)
	}

	return result[short].OriginalURL, nil
}

type DatabaseStorage struct {
	db *sql.DB
	mu sync.RWMutex
}

func NewDatabaseStorage(db *sql.DB) *DatabaseStorage {
	return &DatabaseStorage{
		db: db,
	}
}

func (m *DatabaseStorage) Add(url string, short string) (err error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	id := uuid.New()
	_, err = m.db.ExecContext(context.Background(), "INSERT INTO link (uuid, original_url, short_url) VALUES ($1, $2, $3)", id, url, short)

	if err != nil {
		return err
	}

	return nil
}

func (m *DatabaseStorage) Get(short string) (string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	row := m.db.QueryRowContext(context.Background(), "SELECT original_url FROM link WHERE short_url = $1", short)
	var originalURL string
	err := row.Scan(&originalURL)

	if err != nil {
		return "", err
	}

	return originalURL, nil
}
