package internal

type Database interface {
	Get(key string) string
	Set(key, value string)
}

type MapDatabase struct {
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
