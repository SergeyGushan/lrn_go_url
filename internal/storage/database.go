package storage

import (
	"database/sql"
	"fmt"
	"github.com/SergeyGushan/lrn_go_url/cmd/config"
	"log"
	"strings"
	"sync"
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

type DeletedError struct {
}

func (e DeletedError) Error() string { return "url deleted" }

func NewDatabaseStorage(db *sql.DB) *DatabaseStorage {
	return &DatabaseStorage{db: db}
}

func (ds *DatabaseStorage) GetOriginalURL(shortURL string) (string, error) {
	var originalURL string
	var isDeleted bool

	err := ds.db.QueryRow("SELECT original_url, is_deleted FROM urls WHERE short_url = $1", shortURL).Scan(&originalURL, &isDeleted)

	if err != nil {
		return "", fmt.Errorf("shortURL not found")
	}

	if isDeleted {
		return "", &DeletedError{}
	}

	return originalURL, nil
}

func (ds *DatabaseStorage) Save(shortURL, originalURL, userID string) error {
	result, err := ds.db.Exec("INSERT INTO urls (user_id, short_url, original_url) VALUES ($1, $2, $3) ON CONFLICT (original_url) DO NOTHING RETURNING id", userID, shortURL, originalURL)
	rowsAffected, _ := result.RowsAffected()

	if rowsAffected == 0 {
		return &DuplicateError{
			Field:    "originalURL",
			ShortURL: fmt.Sprintf("%s/%s", config.Opt.BaseURL, shortURL),
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
			ShortURL:      fmt.Sprintf("%s/%s", config.Opt.BaseURL, item.ShortURL),
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
		return results
	}

	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(rows)

	if rows.Err() != nil {
		log.Fatal(rows.Err())
		return results
	}

	for rows.Next() {
		var result URLSByUserIDResult
		err := rows.Scan(&result.ShortURL, &result.OriginalURL)
		if err != nil {
			log.Fatal(err)
		}

		result.ShortURL = fmt.Sprintf("%s/%s", config.Opt.BaseURL, result.ShortURL)

		results = append(results, result)
	}

	return results
}

func (ds *DatabaseStorage) DeleteURLS(urls []string, userID string) {
	placeholders := strings.Repeat("$2,", len(urls))
	if len(urls) > 0 {
		placeholders = placeholders[:len(placeholders)-1] // Убираем последнюю запятую
	}

	query := fmt.Sprintf("UPDATE urls SET is_deleted = true WHERE short_url IN (%s) AND user_id = $1;", placeholders)

	// Готовим запрос
	stmt, err := ds.db.Prepare(query)
	if err != nil {
		return
	}

	// Создаем канал для передачи URL в fan-in
	updatesCh := make(chan string, len(urls))

	// Функция для ожидания завершения fan-in
	var wg sync.WaitGroup
	defer wg.Wait()

	for _, url := range urls {
		wg.Add(1)
		go func(url string) {
			defer wg.Done()
			_, err := stmt.Exec(userID, url)
			if err != nil {
				println(err.Error())
			}
		}(url)
	}

	for _, url := range urls {
		updatesCh <- url
	}

	close(updatesCh) // Закрываем канал после передачи всех URL
}
