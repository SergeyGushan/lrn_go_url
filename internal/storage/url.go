package storage

type Url struct {
	Urls map[string]string
}

var UrlStore = Url{}

func (s *Url) Push(key, value string) {
	if s.Urls == nil {
		s.Urls = make(map[string]string)
	}

	s.Urls[key] = value
}

func (s *Url) GetByKey(key string) (string, bool) {
	value, exists := s.Urls[key]
	return value, exists
}
