package gzip

import (
	"compress/gzip"
	"net/http"
	"strings"
)

type Middleware struct {
	Next http.Handler
}

func (gm *Middleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
		gm.Next.ServeHTTP(w, r)
		return
	}

	writer := gzipResponseWriter{
		ResponseWriter: w,
		Writer:         gzip.NewWriter(w),
	}

	defer func(writer *gzipResponseWriter) {
		err := writer.Close()
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	}(&writer)

	r.Header.Add("Accept-Encoding", "gzip")

	gm.Next.ServeHTTP(writer, r)
}

func Handler(next http.Handler) http.Handler {
	return &Middleware{Next: next}
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
