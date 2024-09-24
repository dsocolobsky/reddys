package internal

type Database struct {
	data map[string]string
}

func NewDatabase() *Database {
	return &Database{
		data: make(map[string]string),
	}
}

func (db *Database) Get(key string) string {
	return db.data[key]
}

func (db *Database) Set(key, value string) {
	db.data[key] = value
}
