// internal/players/repository.go
package players

import (
	"context"
	"database/sql"
	"fmt"
)

type Repository interface {
    FindAll(ctx context.Context) ([]PlayerResponse, error)
}

type postgresRepository struct {
    db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
    return &postgresRepository{db: db}
}

func (r *postgresRepository) FindAll(ctx context.Context) ([]PlayerResponse, error) {
    query := `
        SELECT
            p.id, p.name,
            t.id, t.name, t.flag
        FROM players p
        JOIN teams t ON t.id = p.team_id
        ORDER BY t.name ASC, p.name ASC
    `

    rows, err := r.db.QueryContext(ctx, query)
    if err != nil {
        return nil, fmt.Errorf("erro ao buscar jogadores: %w", err)
    }
    defer rows.Close()

    var result []PlayerResponse
    for rows.Next() {
        var p PlayerResponse
        err := rows.Scan(
            &p.ID, &p.Name,
            &p.Team.ID, &p.Team.Name, &p.Team.Flag,
        )
        if err != nil {
            return nil, fmt.Errorf("erro ao ler jogador: %w", err)
        }
        result = append(result, p)
    }

    if err := rows.Err(); err != nil {
        return nil, fmt.Errorf("erro ao iterar jogadores: %w", err)
    }

    return result, nil
}