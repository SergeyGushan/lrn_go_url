package storage

type URLStorage interface {
	GetOriginalURL(shortURL string) (string, error)
	Save(shortURL, originalURL, userID string) error
	SaveBatch(batch []BatchItem) ([]BatchResult, error)
	GetURLByUserID(userID string) []URLSByUserIDResult
}

var Service URLStorage

type URLSByUserIDResult struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

type BatchResult struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}

type BatchItemReq struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}

type BatchItem struct {
	UserID        string `json:"user_id"`
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
	ShortURL      string `json:"short_url"`
}
