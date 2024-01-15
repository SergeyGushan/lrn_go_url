package storage

import "fmt"

type URL struct {
	Urls       map[string]string
	UrlsByUser map[string]map[string]string
}

func NewStorage() (*URL, error) {
	url := &URL{
		Urls:       make(map[string]string),
		UrlsByUser: make(map[string]map[string]string),
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

func (u *URL) Save(shortURL, originalURL, userID string) error {
	u.Urls[shortURL] = originalURL

	_, exists := u.UrlsByUser[userID]

	if !exists {
		u.UrlsByUser[userID] = make(map[string]string)
	}

	u.UrlsByUser[userID][shortURL] = originalURL
	return nil
}

func (u *URL) SaveBatch(batch []BatchItem) ([]BatchResult, error) {
	results := make([]BatchResult, 0, len(batch))

	for _, item := range batch {
		u.Urls[item.ShortURL] = item.OriginalURL

		_, exists := u.UrlsByUser[item.UserID]
		if !exists {
			u.UrlsByUser[item.UserID] = make(map[string]string)
		}

		u.UrlsByUser[item.UserID][item.ShortURL] = item.OriginalURL

		results = append(results, BatchResult{
			CorrelationID: item.CorrelationID,
			ShortURL:      item.ShortURL,
		})
	}

	return results, nil
}

func (u *URL) GetURLByUserID(userID string) []URLSByUserIDResult {
	items, exists := u.UrlsByUser[userID]
	results := make([]URLSByUserIDResult, 0, len(items))

	if exists {
		for shortURL, getOriginalURL := range items {
			results = append(results, URLSByUserIDResult{
				ShortURL:    shortURL,
				OriginalURL: getOriginalURL,
			})
		}
	}

	return results
}

func (u *URL) DeleteURLS(urls []string, userID string) {

}
