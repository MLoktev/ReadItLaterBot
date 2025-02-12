package sqlite

import (
	"database/sql"
	"fmt"
	"read-it-later-bot/storage"
)

type Storage struct {
	db *sql.DB
}

func New(path string) (*Storage, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, fmt.Errorf("can't open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("can't connect to database: %w", err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) Save(p *storage.Page) error {

	q := `INSERT INTO pages (url, user_name) VALUES (?, ?)`

	s.db.Exec()

}
func (s *Storage) PickRandom(UserName string) (*storage.Page, error) {}
func (s *Storage) Remove(p *storage.Page) error                      {}
func (s *Storage) IsExists(p *storage.Page) (bool, error)            {}
