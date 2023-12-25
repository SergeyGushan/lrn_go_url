package storage

import "fmt"

type URL struct {
	Urls map[string]string
}

func NewStorage() (*URL, error) {
	url := &URL{
		Urls: make(map[string]string),
	}

	return url, nil
}

func (u *URL) GetOriginalURL(shortURL string) (string, error) {
	value, exists := u.Urls[shortURL]

	if !exists {
		return "", fmt.Errorf("shortURL not found")
	}
	return value, nil
}

func (u *URL) Save(shortURL, originalURL string) error {
	u.Urls[shortURL] = originalURL
	return nil
}
