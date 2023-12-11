package gzip

import (
	"compress/gzip"
	"net/http"
	"strings"
)

func Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}

		reader, err := gzip.NewReader(r.Body)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		defer func(reader *gzip.Reader) {
			err := reader.Close()
			if err != nil {
				panic(err)
			}
		}(reader)

		r.Body = reader
		next.ServeHTTP(w, r)
	})
}

type gzipResponseWriter struct {
	http.ResponseWriter
	*gzip.Writer
}

func (w gzipResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func (w gzipResponseWriter) WriteHeader(statusCode int) {
	w.ResponseWriter.WriteHeader(statusCode)
}

func (w gzipResponseWriter) Header() http.Header {
	return w.ResponseWriter.Header()
}
