package main

import (
	"fmt"
	"log"
	"net/http"

	"UAS/config"
	"UAS/database"
)

func main() {
	cfg := config.Load()

	// Connect DB
	database.Connect(cfg.DB_DSN)

	// Simple test route
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "🚀 API Running. Database Connected Successfully.")
	})

	log.Printf("Server running on http://localhost:%s\n", cfg.AppPort)

	// Start server
	err := http.ListenAndServe(":"+cfg.AppPort, nil)
	if err != nil {
		log.Fatal("❌ Server failed:", err)
	}
}
