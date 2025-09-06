package database

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type DB = sqlx.DB

func Connect(databaseURL string) (*DB, error) {
	fmt.Println("databaseURL:", databaseURL)
	db, err := sqlx.Connect("postgres", databaseURL)
	if err != nil {
		return nil, err
	}

	// Connection pool settings
	db.SetMaxOpenConns(25) // Max connections
	db.SetMaxIdleConns(5)  // Idle connections
	db.SetConnMaxLifetime(5 * time.Minute)
	db.SetConnMaxIdleTime(1 * time.Minute)

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
