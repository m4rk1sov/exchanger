package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"log/slog"
	"m4rk1sov/exchanger/internal/app"
	"net/http"
	"os"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("Error loading .env file")
	}
	
	log := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	
	http.HandleFunc("/", app.Handler)
	log.Info("Starting server on port 4000")
	err = http.ListenAndServe("localhost:4000", nil)
	if err != nil {
		log.Error("Error starting server: ", err)
	}
}
