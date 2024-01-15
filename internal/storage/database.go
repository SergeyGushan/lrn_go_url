package storage

import (
	"database/sql"
	"fmt"
	"github.com/SergeyGushan/lrn_go_url/cmd/config"
	"github.com/SergeyGushan/lrn_go_url/internal/logger"
	"log"
	"strings"
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
		logger.Log.Error(err.Error())
		return results
	}

	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(rows)

	if rows.Err() != nil {
		logger.Log.Error(rows.Err().Error())
		return results
	}

	for rows.Next() {
		var result URLSByUserIDResult
		err := rows.Scan(&result.ShortURL, &result.OriginalURL)
		if err != nil {
			logger.Log.Error(err.Error())
			continue
		}

		result.ShortURL = fmt.Sprintf("%s/%s", config.Opt.BaseURL, result.ShortURL)

		results = append(results, result)
	}

	return results
}

func (ds *DatabaseStorage) DeleteURLS(urls []string, userID string) {
	workersCount := 10
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

	// Создаем канал для передачи URL в пул воркеров
	updatesCh := make(chan string, len(urls))

	// Создаем канал для сигнала об окончании работы
	done := make(chan struct{})
	defer close(done)

	// Запускаем пул воркеров
	for i := 0; i < workersCount; i++ {
		go worker(updatesCh, stmt, userID, done)
	}

	// Передаем URL-ы в канал
	for _, url := range urls {
		updatesCh <- url
	}

	// Закрываем канал, чтобы завершить работу воркеров
	close(updatesCh)

	// Ожидаем завершения всех воркеров
	for i := 0; i < workersCount; i++ {
		<-done
	}
}

func worker(updatesCh <-chan string, stmt *sql.Stmt, userID string, done chan<- struct{}) {
	defer func() {
		done <- struct{}{} // Сигнализируем о завершении работы воркера
	}()

	for url := range updatesCh {
		_, err := stmt.Exec(userID, url)
		if err != nil {
			logger.Log.Error(err.Error())
		}
	}
}
