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

type Database struct {
	table map[string]string
}

func initDatabase() Database {
	return Database{table: make(map[string]string)}
}

func (db *Database) Set(key, value string) error {
	// FIXME: Should it repeatedly update key value or create it once?
	//
	// _, ok := db.table[key]
	// if ok {
	//     return errors.New("key already exists")
	// }

	db.table[key] = value

	return nil
}

func (db *Database) Get(key string) (string, error) {
	value, ok := db.table[key]
	if !ok {
		return "", errors.New("key not found")
	}
	return value, nil
}
