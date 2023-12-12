package storage

type URL struct {
	Urls   map[string]string
	Record *Records
}

var URLStore *URL

func NewURL(fileName string) (*URL, error) {
	record, err := (&Records{}).New(fileName)
	if err != nil {
		return nil, err
	}

	url := &URL{
		Urls:   make(map[string]string),
		Record: record,
	}

	// Загрузка записей из файла в URL.Urls
	if err := url.LoadFromFile(fileName); err != nil {
		return nil, err
	}

	return url, nil
}

func (s *URL) Push(key, value string) {
	// Добавление записи в URL.Urls
	s.Urls[key] = value

	// Создание новой записи и сохранение в файл
	record := &Record{
		UUID:        0, // Используйте соответствующий UUID
		ShortURL:    key,
		OriginalURL: value,
	}

	// Сохранение в файл
	if err := s.Record.WriteEvent(record); err != nil {
		panic(err) // Обработка ошибки по вашему усмотрению
	}
}

func (s *URL) GetByKey(key string) (string, bool) {
	value, exists := s.Urls[key]
	return value, exists
}

func (s *URL) LoadFromFile(fileName string) error {
	// Загрузка записей из файла в URL.Urls
	if err := s.Record.LoadFromFile(fileName); err != nil {
		return err
	}

	// Заполнение URL.Urls
	for _, record := range s.Record.Items {
		s.Urls[record.ShortURL] = record.OriginalURL
	}

	return nil
}
