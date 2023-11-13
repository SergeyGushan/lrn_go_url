package logger

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"go.uber.org/zap"
)

var Log = zap.NewNop()

func Initialize(level string) error {
	// преобразуем текстовый уровень логирования в zap.AtomicLevel
	lvl, err := zap.ParseAtomicLevel(level)
	if err != nil {
		return err
	}
	// создаём новую конфигурацию логера
	cfg := zap.NewProductionConfig()
	// устанавливаем уровень
	cfg.Level = lvl
	// создаём логер на основе конфигурации
	zl, err := cfg.Build()
	if err != nil {
		return err
	}
	// устанавливаем синглтон
	Log = zl
	return nil
}

type ResponseWriter interface {
	Write([]byte) (int, error)
	WriteHeader(statusCode int)
}

type (
	// берём структуру для хранения сведений об ответе
	responseData struct {
		status int
		size   int
	}

	// добавляем реализацию http.ResponseWriter
	loggingResponseWriter struct {
		http.ResponseWriter // встраиваем оригинальный http.ResponseWriter
		responseData        *responseData
		Body                *bytes.Buffer // добавляем поле для записи тела ответа
	}
)

func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	// записываем ответ, используя оригинальный http.ResponseWriter
	size, err := r.ResponseWriter.Write(b)
	r.responseData.size += size // захватываем размер
	// также записываем тело ответа в буфер
	r.Body.Write(b)
	return size, err
}

func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	// записываем код статуса, используя оригинальный http.ResponseWriter
	r.ResponseWriter.WriteHeader(statusCode)
	r.responseData.status = statusCode // захватываем код статуса
}

func Handler(next http.Handler) http.Handler {
	logFn := func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		responseData := &responseData{
			status: 0,
			size:   0,
		}
		// создаем буфер для записи тела ответа
		var buf bytes.Buffer
		lw := loggingResponseWriter{
			ResponseWriter: w, // встраиваем оригинальный http.ResponseWriter
			responseData:   responseData,
			Body:           &buf, // передаем буфер для записи тела ответа
		}
		next.ServeHTTP(&lw, r) // внедряем реализацию http.ResponseWriter

		duration := time.Since(start)
		var reqBuf bytes.Buffer
		teeBody := io.TeeReader(r.Body, &reqBuf)
		r.Body = io.NopCloser(teeBody)

		Log.Info("Request",
			zap.String("uri", r.RequestURI),
			zap.String("method", r.Method),
			zap.String("status", strconv.Itoa(responseData.status)),
			zap.String("duration", duration.String()),
			zap.String("size", strconv.Itoa(responseData.size)),
			zap.String("request_headers", headersToString(r.Header)),
			zap.String("request_body", reqBuf.String()),
			zap.String("response_headers", headersToString(lw.Header())),
			zap.String("response_body", buf.String()), // записываем тело ответа из буфера
		)
	}

	return http.HandlerFunc(logFn)
}

// Вспомогательная функция для преобразования заголовков в строку
func headersToString(headers http.Header) string {
	var headerString strings.Builder
	for key, values := range headers {
		headerString.WriteString(fmt.Sprintf("%s: %s\n", key, strings.Join(values, ", ")))
	}
	return headerString.String()
}
