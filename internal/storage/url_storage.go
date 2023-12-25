package storage

type URLStorage interface {
	GetOriginalURL(shortURL string) (string, error)
	Save(shortURL, originalURL string) error
}

var Service URLStorage
