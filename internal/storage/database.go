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

func (ds *DatabaseStorage) SaveBatch(batch []BatchItem) ([]BatchResult, error) {
	tx, err := ds.db.Begin()
	if err != nil {
		return nil, err
	}
	defer func() {
		if p := recover(); p != nil {
			err := tx.Rollback()
			if err != nil {
				return
			}
		}
	}()

	results := make([]BatchResult, 0, len(batch))
	for _, item := range batch {

		_, err := tx.Exec("INSERT INTO urls (short_url, original_url) VALUES ($1, $2)", item.ShortURL, item.OriginalURL)
		if err != nil {
			err := tx.Rollback()
			if err != nil {
				return nil, err
			}
			return nil, err
		}

		results = append(results, BatchResult{
			CorrelationID: item.CorrelationID,
			ShortURL:      item.ShortURL,
		})
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return results, nil
}
