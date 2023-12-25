package storage

import (
	"database/sql"
	"fmt"
)

type DatabaseStorage struct {
	db *sql.DB
}

type DuplicateError struct {
	Field string
}

func (e DuplicateError) Error() string {
	return fmt.Sprintf("%s already exists", e.Field)
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
	result, err := ds.db.Exec("INSERT INTO urls (short_url, original_url) VALUES ($1, $2) ON CONFLICT (original_url) DO NOTHING RETURNING id", shortURL, originalURL)
	rowsAffected, _ := result.RowsAffected()

	if rowsAffected == 0 {
		return &DuplicateError{
			Field: "originalURL",
		}
	}

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

		result, err := tx.Exec("INSERT INTO urls (short_url, original_url) VALUES ($1, $2) ON CONFLICT (original_url) DO NOTHING RETURNING id", item.ShortURL, item.OriginalURL)
		rowsAffected, _ := result.RowsAffected()

		if rowsAffected != 0 {
			return nil, &DuplicateError{
				Field: "originalURL",
			}
		}

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
