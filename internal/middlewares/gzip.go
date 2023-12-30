package middlewares

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

// GzipMiddleware обрабатывает gzip-кодирование запроса.
func GzipMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		var bodyReader io.Reader = req.Body

		if strings.Contains(req.Header.Get("Content-Encoding"), "gzip") {
			reader, err := gzip.NewReader(req.Body)
			if err != nil {
				http.Error(res, "Internal Server Error", http.StatusInternalServerError)
				return
			}
			defer func(reader *gzip.Reader) {
				err := reader.Close()
				if err != nil {
					panic(err)
				}
			}(reader)
			bodyReader = reader
		}

		// Создаем кастомный ReadCloser для req.Body
		req.Body = &customReadCloser{
			Reader: bodyReader,
			Closer: req.Body,
		}

		next.ServeHTTP(res, req)
	})
}

// customReadCloser является кастомным ReadCloser, включающим io.Reader и io.Closer.
type customReadCloser struct {
	io.Reader
	io.Closer
}
