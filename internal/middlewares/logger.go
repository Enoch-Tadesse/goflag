package middlewares

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		lrw := &loggingResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		next.ServeHTTP(lrw, r)

		duration := time.Since(start)
		log.Printf("%s %s %s %s", r.Method, r.URL.Path, colorStatusCode(lrw.statusCode), duration)
	})
}

func colorStatusCode(code int) string {
	switch {
	case code >= 200 && code < 300:
		return fmt.Sprintf("\033[32m%d\033[0m", code) // Green for success
	case code >= 300 && code < 400:
		return fmt.Sprintf("\033[36m%d\033[0m", code) // Cyan for redirect
	case code >= 400 && code < 500:
		return fmt.Sprintf("\033[33m%d\033[0m", code) // Yellow for client error
	case code >= 500:
		return fmt.Sprintf("\033[31m%d\033[0m", code) // Red for server error
	default:
		return fmt.Sprintf("%d", code) // Default color (no color)
	}
}

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}
