package storage

type URL struct {
	Urls map[string]string
}

var URLStore = URL{}

func (s *URL) Push(key, value string) {
	if s.Urls == nil {
		s.Urls = make(map[string]string)
	}

	s.Urls[key] = value
}

func (s *URL) GetByKey(key string) (string, bool) {
	value, exists := s.Urls[key]
	return value, exists
}
