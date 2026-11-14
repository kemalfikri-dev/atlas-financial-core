package main

import (
	"log/slog"
	"net/http"
	"time"
)

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func newResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{w, http.StatusOK}
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func responseLog(next http.Handler, logger *slog.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		nw := newResponseWriter(w)

		next.ServeHTTP(nw, r)

		duration := time.Since(start)

		logger.Info("Logging Info",
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
			slog.Int("status", nw.statusCode),
			slog.Duration("duration", duration),
		)
	})
}
