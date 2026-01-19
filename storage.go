package main

import "errors"

type Storage struct {
	databases map[string]Database
}

func initStorage() Storage {
	return Storage{databases: make(map[string]Database)}
}

func (s *Storage) CreateDatabase(db string) error {
	_, ok := s.databases[db]
	if ok {
		return errors.New("database already exists")
	}

	s.databases[db] = initDatabase()

	return nil
}

func (s *Storage) GetDatabase(db string) (*Database, error) {
	database, ok := s.databases[db]
	if !ok {
		return nil, errors.New("database not found")
	}

	return &database, nil
}

func (s *Storage) DeleteDatabase(db string) error {
	_, ok := s.databases[db]
	if !ok {
		return errors.New("database not found")
	}

	delete(s.databases, db)
	return nil
}

func (s *Storage) GetDBNames() (names []string) {
	for key := range s.databases {
		names = append(names, key)
	}
	return
}

type Database struct {
	table map[string]string
}

func initDatabase() Database {
	return Database{table: make(map[string]string)}
}

func (db *Database) Set(key, value string) error {
	db.table[key] = value

	return nil
}

func (db *Database) Get(key string) *string {
	value, ok := db.table[key]
	if !ok {
		return nil
	}
	return &value
}

func (db *Database) Delete(key string) error {
	_, ok := db.table[key]
	if !ok {
		return errors.New("key not found")
	}

	delete(db.table, key)
	return nil
}

func (db *Database) GetAllKeys() *map[string]string {
	return &db.table
}
