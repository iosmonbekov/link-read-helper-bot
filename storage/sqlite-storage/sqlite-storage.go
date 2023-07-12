package sqlitestorage

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

type Storage struct {
	db *sql.DB
}

func New(path string) (*Storage, error) {
	db, err := sql.Open("sqlite3", "__"+path)
	if err != nil {
		return nil, fmt.Errorf("sqlitestorage.New: error on sql.Open -> %v\n", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("sqlitestorage.New: error on db.Ping -> %v\n", err)
	}
	return &Storage{db: db}, nil
}

func (s *Storage) Save(username string, content []byte) error {
	_, err := s.db.Exec("REPLACE INTO storage(username, content) VALUES(?, ?)", username, content)
	if err != nil {
		return fmt.Errorf("sqlitestorage.Save: error on db.Exec -> %v\n", err)
	}
	return nil
}

func (s *Storage) Load(username string) ([]byte, error) {
	var content []byte
	err := s.db.QueryRow("SELECT content FROM storage WHERE username=?", username).Scan(&content)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("sqlitestorage.Load: error on db.QueryRow -> %v\n", err)
	}
	return content, nil
}

func (s *Storage) Remove(username string) error {
	_, err := s.db.Exec("DELETE FROM storage WHERE username=?", username)
	if err != nil {
		return fmt.Errorf("sqlitestorage.Remove: error on db.Exec -> %v\n", err)
	}
	return nil
}

func (s *Storage) Init() error {
	q := `CREATE TABLE IF NOT EXISTS storage (username TEXT PRIMARY KEY, content BLOB)`

	_, err := s.db.Exec(q)
	if err != nil {
		return fmt.Errorf("sqlitestorage.Init: error on db.Exec -> %v\n", err)
	}

	return nil
}
