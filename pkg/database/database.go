package database

import "sync"

type Database interface {
	Get(key string) string
	Set(key, value string)
	Lock()
	Unlock()
}

type MapDatabase struct {
	mu   sync.Mutex
	data map[string]string
}

func NewMapDatabase() *MapDatabase {
	return &MapDatabase{
		data: make(map[string]string),
	}
}

func (db *MapDatabase) Get(key string) string {
	return db.data[key]
}

func (db *MapDatabase) Set(key, value string) {
	db.data[key] = value
}

func (db *MapDatabase) Lock() {
	db.mu.Lock()
}

func (db *MapDatabase) Unlock() {
	db.mu.Unlock()
}
