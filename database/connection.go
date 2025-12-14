package database

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func Connect(dsn string) *sql.DB {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("❌ Failed to connect DB:", err)
	}

	// Test koneksi
	if err := db.Ping(); err != nil {
		log.Fatal("❌ DB ping failed:", err)
	}

	DB = db
	log.Println("✅ Database connected")
	return db
}
