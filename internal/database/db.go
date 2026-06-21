package db

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/lib/pq"
)

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
    sql, err := os.ReadFile("db/migrations/001_init.sql")
    if err != nil {
        return fmt.Errorf("erro ao ler migration: %w", err)
    }

    _, err = db.Exec(string(sql))
    if err != nil {
        return fmt.Errorf("erro ao executar migration: %w", err)
    }

    return nil
}