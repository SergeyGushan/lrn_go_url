package main

import (
	"compress/gzip"
	"github.com/SergeyGushan/lrn_go_url/cmd/config"
	"github.com/SergeyGushan/lrn_go_url/internal/logger"
	"github.com/SergeyGushan/lrn_go_url/internal/urlhandlers"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strings"
)

func URLRouter() chi.Router {
	r := chi.NewRouter()
	r.Use(GzipMiddleware)
	r.Use(logger.Handler)
	r.Post("/", urlhandlers.Save)
	r.Post("/api/shorten", urlhandlers.Shorten)
	r.Get("/{shortCode}", urlhandlers.Get)

	return r
}

func main() {
	err := logger.Initialize("Info")
	if err != nil {
		panic(err)
	}

	config.SetOptions()
	err = http.ListenAndServe(config.Opt.ServerAddress, URLRouter())
	if err != nil {
		panic(err)
	}
}

func GzipMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Проверяем, поддерживает ли клиент сжатие
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}

		// Устанавливаем заголовки для сжатия
		w.Header().Set("Content-Encoding", "gzip")
		w.Header().Set("Vary", "Accept-Encoding")

		// Создаем gzip.Writer и устанавливаем в ResponseWriter
		gzWriter := gzip.NewWriter(w)
		defer gzWriter.Close()
		gzResponseWriter := &gzipResponseWriter{Writer: gzWriter, ResponseWriter: w}

		// Проходим дальше по цепочке сжатия данных
		next.ServeHTTP(gzResponseWriter, r)
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
