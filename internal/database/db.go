// internal/database/db.go
package database

import (
	"database/sql"
	_ "embed"
	"fmt"
	"os"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

//go:embed migrations/001_init.sql
var initSQL string

//go:embed migrations/002_add_minute_to_matches.sql
var addMinuteSQL string

//go:embed migrations/003_add_goals_table.sql
var addGoalsTableSQL string

func Connect() (*sql.DB, error) {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		return nil, fmt.Errorf("DATABASE_URL não definida")
	}

	db, err := sql.Open("pgx", dsn)
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
    // Criar tabela de controle de migrations se não existir
    _, err := db.Exec(`
        CREATE TABLE IF NOT EXISTS schema_migrations (
            version TEXT PRIMARY KEY,
            applied_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
        )
    `)
    if err != nil {
        return fmt.Errorf("erro ao criar tabela de migrations: %w", err)
    }

    // Função auxiliar para verificar se migration já foi aplicada
    isApplied := func(version string) (bool, error) {
        var count int
        err := db.QueryRow("SELECT COUNT(*) FROM schema_migrations WHERE version = $1", version).Scan(&count)
        return count > 0, err
    }

    // Função auxiliar para marcar migration como aplicada
    markApplied := func(version string) error {
        _, err := db.Exec("INSERT INTO schema_migrations (version) VALUES ($1)", version)
        return err
    }

    // Executar migration 001
    applied, err := isApplied("001")
    if err != nil {
        return fmt.Errorf("erro ao verificar migration 001: %w", err)
    }
    if !applied {
        if _, err := db.Exec(initSQL); err != nil {
            return fmt.Errorf("erro ao executar migration 001: %w", err)
        }
        if err := markApplied("001"); err != nil {
            return fmt.Errorf("erro ao marcar migration 001: %w", err)
        }
    }

    // Executar migration 002
    applied, err = isApplied("002")
    if err != nil {
        return fmt.Errorf("erro ao verificar migration 002: %w", err)
    }
    if !applied {
        if _, err := db.Exec(addMinuteSQL); err != nil {
            return fmt.Errorf("erro ao executar migration 002: %w", err)
        }
        if err := markApplied("002"); err != nil {
            return fmt.Errorf("erro ao marcar migration 002: %w", err)
        }
    }

    // Executar migration 003
    applied, err = isApplied("003")
    if err != nil {
        return fmt.Errorf("erro ao verificar migration 003: %w", err)
    }
    if !applied {
        if _, err := db.Exec(addGoalsTableSQL); err != nil {
            return fmt.Errorf("erro ao executar migration 003: %w", err)
        }
        if err := markApplied("003"); err != nil {
            return fmt.Errorf("erro ao marcar migration 003: %w", err)
        }
    }

    return nil
}