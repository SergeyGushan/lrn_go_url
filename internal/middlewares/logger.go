package middlewares

import (
	"bytes"
	"fmt"
	"github.com/SergeyGushan/lrn_go_url/internal/logger"
	"net/http"
	"strconv"
	"strings"
	"time"

	"go.uber.org/zap"
)

type (
	ResponseWriter interface {
		Write([]byte) (int, error)
		WriteHeader(statusCode int)
	}

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

func LoggerMiddleware(next http.Handler) http.Handler {
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

		logger.Log.Info("Request",
			zap.String("uri", r.RequestURI),
			zap.String("method", r.Method),
			zap.String("status", strconv.Itoa(responseData.status)),
			zap.String("duration", duration.String()),
			zap.String("size", strconv.Itoa(responseData.size)),
			zap.String("request_headers", headersToString(r.Header)),
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
