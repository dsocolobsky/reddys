package database

import "sync"

// Database is an interface that defines the methods that a database must implement
type Database interface {
	Get(key string) string
	Set(key, value string)
	HGet(key, field string) string
	HSet(key, field, value string)
	Lock()
	Unlock()
}

// MapDatabase is an in-memory database that uses a map as the underlying data store
type MapDatabase struct {
	mu   sync.Mutex
	data map[string]string
	hset map[string]map[string]string
}

// NewMapDatabase creates a new MapDatabase using a map as the underlying data store
func NewMapDatabase() *MapDatabase {
	return &MapDatabase{
		data: make(map[string]string),
		hset: make(map[string]map[string]string),
	}
}

// Get retrieves the value of a key from the database
func (db *MapDatabase) Get(key string) string {
	return db.data[key]
}

// Set sets the value of a key in the database
func (db *MapDatabase) Set(key, value string) {
	db.data[key] = value
}

// HGet retrieves the value of a field from a hash, returning an empty string if the field doesn't exist
func (db *MapDatabase) HGet(key, field string) string {
	if _, ok := db.hset[key]; !ok {
		return ""
	}
	return db.hset[key][field]
}

// HSet sets the value of a field in a hash, creating the hash if it doesn't exist
func (db *MapDatabase) HSet(key, field, value string) {
	if _, ok := db.hset[key]; !ok {
		db.hset[key] = make(map[string]string)
	}
	db.hset[key][field] = value
}

// Lock locks the database
func (db *MapDatabase) Lock() {
	db.mu.Lock()
}

// Unlock unlocks the database
func (db *MapDatabase) Unlock() {
	db.mu.Unlock()
}
