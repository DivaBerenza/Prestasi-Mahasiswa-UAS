package database

import (
	"context"
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func Connect(dsn string) {
	conn, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("❌ Failed connecting to DB:", err)
	}

	// Test connection
	if err := conn.PingContext(context.Background()); err != nil {
		log.Fatal("❌ Database NOT reachable:", err)
	}

	log.Println("✅ PostgreSQL Connected Successfully!")
	DB = conn
}
