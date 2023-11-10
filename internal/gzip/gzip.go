package gzip

import (
	"compress/gzip"
	"net/http"
	"strings"
)

func Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Проверяем, поддерживает ли клиент сжатие
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}

		w.Header().Set("Content-Encoding", "gzip")

		gz := gzip.NewWriter(w)
		defer func(gz *gzip.Writer) {
			err := gz.Close()
			if err != nil {
				panic(err)
			}
		}(gz)

		gzWriter := gzipResponseWriter{ResponseWriter: w, Writer: gz}
		next.ServeHTTP(gzWriter, r)
	})
}

type gzipResponseWriter struct {
	http.ResponseWriter
	*gzip.Writer
}

func (w gzipResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}
