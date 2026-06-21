// internal/database/db.go
package database

import (
	"database/sql"
	_ "embed"
	"fmt"
	"os"
	"time"

	_ "github.com/lib/pq"
)

//go:embed migrations/001_init.sql
var initSQL string

func Connect() (*sql.DB, error) {
    dsn := os.Getenv("DATABASE_URL")
    if dsn == "" {
        return nil, fmt.Errorf("DATABASE_URL não definida")
    }

    db, err := sql.Open("postgres", dsn)
    if err != nil {
        return nil, fmt.Errorf("erro ao abrir conexão: %w", err)
    }

    if err := db.Ping(); err != nil {
        return nil, fmt.Errorf("erro ao conectar no banco: %w", err)
    }

    db.SetMaxOpenConns(25)
    db.SetMaxIdleConns(10)
    db.SetConnMaxLifetime(5 * time.Minute)

    if err := migrate(db); err != nil {
        return nil, fmt.Errorf("erro ao executar migrations: %w", err)
    }

    return db, nil
}

func migrate(db *sql.DB) error {
    _, err := db.Exec(initSQL)
    if err != nil {
        return fmt.Errorf("erro ao executar migration: %w", err)
    }
    return nil
}