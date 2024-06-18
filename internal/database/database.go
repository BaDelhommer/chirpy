package database

import (
	"encoding/json"
	"errors"
	"os"
	"sync"
)

type DB struct {
	path string
	mux  *sync.RWMutex
}

type Chirp struct {
	ID   int    `json:"id"`
	Body string `json:"body"`
}

type DBstructure struct {
	Chirps map[int]Chirp `json:"chirps"`
}

func (db *DB) ensureDB() error {
	_, err := os.ReadFile(db.path)
	if errors.Is(err, os.ErrNotExist) {
		return db.createDB()
	}

	return err
}

func (db *DB) createDB() error {
	dbStructure := DBstructure{
		Chirps: map[int]Chirp{},
	}
	return db.writeDB(dbStructure)
}

func (db *DB) writeDB(dbs DBstructure) error {
	db.mux.Lock()
	defer db.mux.Unlock()

	dat, err := json.Marshal(dbs)
	if err != nil {
		return err
	}

	err = os.WriteFile(db.path, dat, 0600)
	if err != nil {
		return err
	}

	return nil
}

func (db *DB) loadDB() (DBstructure, error) {
	db.mux.RLock()
	defer db.mux.RUnlock()

	dbStructure := DBstructure{}
	dat, err := os.ReadFile(db.path)
	if errors.Is(err, os.ErrNotExist) {
		return dbStructure, nil
	}

	err = json.Unmarshal(dat, &dbStructure)
	if err != nil {
		return dbStructure, nil
	}

	return dbStructure, nil
}

func (db *DB) CreateChirp(body string) (Chirp, error) {
	dbs, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}

	id := len(dbs.Chirps) + 1

	chirp := Chirp{
		ID:   id,
		Body: body,
	}

	dbs.Chirps[id] = chirp

	err = db.writeDB(dbs)
	if err != nil {
		return Chirp{}, err
	}

	return chirp, nil
}

func (db *DB) GetChirps() ([]Chirp, error) {
	dbs, err := db.loadDB()
	if err != nil {
		return nil, err
	}

	chirps := make([]Chirp, 0, len(dbs.Chirps))
	chirps = append(chirps, chirps...)

	return chirps, nil
}

func NewDB(path string) (*DB, error) {
	db := &DB{
		path: path,
		mux:  &sync.RWMutex{},
	}

	err := db.ensureDB()

	return db, err
}
