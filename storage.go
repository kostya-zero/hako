package main

import (
	"errors"
	"maps"
	"sync"
)

type Storage struct {
	mu        sync.RWMutex
	databases map[string]*Database
}

func NewStorage() Storage {
	return Storage{databases: make(map[string]*Database)}
}

func (s *Storage) CreateDatabase(db string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, ok := s.databases[db]
	if ok {
		return errors.New("database already exists")
	}

	newDB := NewDatabase()
	s.databases[db] = &newDB

	return nil
}

func (s *Storage) GetDatabase(db string) (*Database, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	database, ok := s.databases[db]
	if !ok {
		return nil, errors.New("database not found")
	}

	return database, nil
}

func (s *Storage) DeleteDatabase(db string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, ok := s.databases[db]
	if !ok {
		return errors.New("database not found")
	}

	delete(s.databases, db)
	return nil
}

func (s *Storage) GetDBNames() (names []string) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for key := range s.databases {
		names = append(names, key)
	}
	return
}

func (s *Storage) MakeSnapshot() map[string]map[string]string {
	s.mu.Lock()
	defer s.mu.Unlock()

	snapshot := make(map[string]map[string]string)

	for name, db := range s.databases {
		keys := db.GetAllKeys()
		snapshot[name] = keys
	}

	return snapshot
}

func (s *Storage) Load(data map[string]map[string]string) {
	for name, keys := range data {
		db := NewDatabaseWithData(keys)
		s.databases[name] = &db
	}
}

type Database struct {
	mu    sync.RWMutex
	table map[string]string
}

func NewDatabase() Database {
	return Database{table: make(map[string]string)}
}

func NewDatabaseWithData(data map[string]string) Database {
	return Database{table: data}
}

func (db *Database) Set(key, value string) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	db.table[key] = value

	return nil
}

func (db *Database) Get(key string) *string {
	db.mu.RLock()
	defer db.mu.RUnlock()

	value, ok := db.table[key]
	if !ok {
		return nil
	}
	return &value
}

func (db *Database) Delete(key string) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	_, ok := db.table[key]
	if !ok {
		return errors.New("key not found")
	}

	delete(db.table, key)
	return nil
}

func (db *Database) GetAllKeys() map[string]string {
	db.mu.RLock()
	defer db.mu.RUnlock()

	keys := make(map[string]string, len(db.table))
	maps.Copy(keys, db.table)

	return keys
}
