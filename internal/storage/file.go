package storage

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"os"
	"strings"
)

type Record struct {
	UUID        string `json:"uuid"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

type JSONStorage struct {
	file    *os.File
	encoder *json.Encoder
	Items   []Record
}

func NewJSONStorage(fileName string) (*JSONStorage, error) {
	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	js := &JSONStorage{
		file:    file,
		encoder: json.NewEncoder(file),
	}

	errLoadFromFile := js.loadFromFile(fileName)

	if errLoadFromFile != nil {
		panic(errLoadFromFile)
	}

	return js, nil
}

func (js *JSONStorage) Save(shortURL, originalURL string) error {
	record := Record{
		UUID:        uuid.New().String(),
		ShortURL:    shortURL,
		OriginalURL: originalURL,
	}

	js.Items = append(js.Items, record)

	return js.encoder.Encode(&record)
}

func (js *JSONStorage) GetOriginalURL(shortURL string) (string, error) {
	for _, record := range js.Items {
		if record.ShortURL == shortURL {
			return record.OriginalURL, nil
		}
	}

	return "", fmt.Errorf("URL not found")
}

func (js *JSONStorage) SaveBatch(batch []BatchItem) ([]BatchResult, error) {
	results := make([]BatchResult, 0, len(batch))

	for _, item := range batch {
		record := Record{
			UUID:        item.CorrelationID,
			ShortURL:    item.ShortURL,
			OriginalURL: item.OriginalURL,
		}

		js.Items = append(js.Items, record)
		err := js.encoder.Encode(&record)
		if err != nil {
			return nil, err
		}

		results = append(results, BatchResult{
			CorrelationID: item.CorrelationID,
			ShortURL:      record.ShortURL,
		})
	}

	return results, nil
}

func (js *JSONStorage) Close() error {
	return js.file.Close()
}

func (js *JSONStorage) loadFromFile(fileName string) error {
	data, err := os.ReadFile(fileName)
	if err != nil {
		return err
	}

	records := strings.Split(string(data), "\n")

	for _, record := range records {
		event := &Record{}
		err := json.Unmarshal([]byte(record), event)
		if err != nil {
			continue
		}
		js.Items = append(js.Items, *event)
	}

	return nil
}
