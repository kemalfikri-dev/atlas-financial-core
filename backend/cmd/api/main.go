package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

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

	if port == "" {
		port = "8080"
	}
	mux := http.NewServeMux()

	mux.HandleFunc("/health", healthHandler	)
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
