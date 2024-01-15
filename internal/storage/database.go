package storage

import (
	"database/sql"
	"fmt"
	"log"
)

type DatabaseStorage struct {
	db *sql.DB
}

type DuplicateError struct {
	Field    string
	ShortURL string
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

func (ds *DatabaseStorage) Save(shortURL, originalURL, userID string) error {
	result, err := ds.db.Exec("INSERT INTO urls (user_id, short_url, original_url) VALUES ($1, $2, $3) ON CONFLICT (original_url) DO NOTHING RETURNING id", userID, shortURL, originalURL)
	rowsAffected, _ := result.RowsAffected()

	if rowsAffected == 0 {
		return &DuplicateError{
			Field:    "originalURL",
			ShortURL: shortURL,
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

		_, err := tx.Exec("INSERT INTO urls (user_id, correlation_id, short_url, original_url) VALUES ($1, $2, $3, $4)", item.UserID, item.CorrelationID, item.ShortURL, item.OriginalURL)

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

func (ds *DatabaseStorage) GetURLByUserID(userID string) []URLSByUserIDResult {
	results := make([]URLSByUserIDResult, 0)
	rows, err := ds.db.Query("SELECT short_url, original_url FROM urls WHERE user_id = $1", userID)
	if err != nil {
		log.Fatal(err)
	}

	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(rows)

	for rows.Next() {
		var result URLSByUserIDResult
		err := rows.Scan(&result.ShortURL, &result.OriginalURL)
		if err != nil {
			log.Fatal(err)
		}
		results = append(results, result)
	}

	return results
}
