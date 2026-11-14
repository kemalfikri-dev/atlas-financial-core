package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
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
		start := time.Now
		nw := newResponseWriter(w)

		next.ServeHTTP(nw, r)

		duration := time.Since(start())

		logger.Info("Logging Info",
			slog.String("Method:", r.Method),
			slog.String("Path:", r.URL.Path),
			slog.Int("Status:", nw.statusCode),
			slog.Duration("Duration:", duration),
			slog.String("ip", r.RemoteAddr),
		)
	})
}

func main() {

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	wd, _ := os.Getwd() 

	err := godotenv.Load()
	if err != nil {
		logger.Error("Gagal memuat .env",
			slog.String("error_detail", err.Error()),
			slog.String("dicari_di_folder", wd),
		)
	}

	port := os.Getenv("PORT")
	mux := http.NewServeMux()

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	loggedHandler := responseLog(mux, logger)

	server := http.Server{
		Addr:    fmt.Sprintf("localhost:%s", port),
		Handler: loggedHandler,
	}
	err = server.ListenAndServe()

	if err != nil {
		panic(err)
	}

}
