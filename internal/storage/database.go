package storage

import (
	"database/sql"
	"fmt"
)

type DatabaseStorage struct {
	db *sql.DB
}

func NewDatabaseStorage(db *sql.DB) *DatabaseStorage {
	return &DatabaseStorage{db: db}
}

func (ds *DatabaseStorage) GetOriginalURL(shortURL string) (string, error) {
	var originalURL string
	err := ds.db.QueryRow("SELECT original_url FROM urls WHERE short_url = $1", shortURL).Scan(&originalURL)
	if err != nil {
		return "", fmt.Errorf("shortURL not found")
	}
	return originalURL, nil
}

func (ds *DatabaseStorage) Save(shortURL, originalURL string) error {
	_, err := ds.db.Exec("INSERT INTO urls (short_url, original_url) VALUES ($1, $2)", shortURL, originalURL)
	return err
}
