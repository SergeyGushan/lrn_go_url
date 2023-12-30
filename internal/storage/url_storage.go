package storage

type URLStorage interface {
	GetOriginalURL(shortURL string) (string, error)
	Save(shortURL, originalURL string) error
	SaveBatch(batch []BatchItem) ([]BatchResult, error)
}

var Service URLStorage

type BatchResult struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}

type BatchItemReq struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}

type BatchItem struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
	ShortURL      string `json:"short_url"`
}
