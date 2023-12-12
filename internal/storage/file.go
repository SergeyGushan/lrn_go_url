package storage

import (
	"encoding/json"
	"os"
	"strings"
)

type Record struct {
	UUID        uint   `json:"uuid"`
	ShortUrl    string `json:"short_url"`
	OriginalUrl string `json:"original_url"`
}

type Records struct {
	file    *os.File
	encoder *json.Encoder
	Items   []Record
}

func (r *Records) New(fileName string) (*Records, error) {
	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(fileName)

	if err != nil {
		panic(err)
	}

	records := strings.Split(string(data), "\n")
	var items []Record

	for _, record := range records {
		event := &Record{}
		err := json.Unmarshal([]byte(record), event)
		if err != nil {
			continue
		}
		items = append(items, *event)
	}

	return &Records{
		file:    file,
		encoder: json.NewEncoder(file),
		Items:   items,
	}, nil
}

func (r *Records) WriteEvent(record *Record) error {
	return r.encoder.Encode(&record)
}

func (r *Records) Close() error {
	return r.file.Close()
}

func (r *Records) LoadFromFile(fileName string) error {
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
		r.Items = append(r.Items, *event)
	}

	return nil
}
