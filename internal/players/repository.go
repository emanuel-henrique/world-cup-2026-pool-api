// internal/players/repository.go
package players

import (
	"context"
	"database/sql"
	"fmt"
)

type Repository interface {
    FindAll(ctx context.Context, limit, offset int) ([]PlayerResponse, int, error)
}

type postgresRepository struct {
    db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
    return &postgresRepository{db: db}
}

func (r *postgresRepository) FindAll(ctx context.Context, limit, offset int) ([]PlayerResponse, int, error) {
    var total int
    err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM players").Scan(&total)
    if err != nil {
        return nil, 0, fmt.Errorf("erro ao contar jogadores: %w", err)
    }

    query := `
        SELECT
            p.id, p.name,
            t.id, t.name, t.flag
        FROM players p
        JOIN teams t ON t.id = p.team_id
        ORDER BY t.name ASC, p.name ASC
        LIMIT $1 OFFSET $2
    `

    rows, err := r.db.QueryContext(ctx, query, limit, offset)
    if err != nil {
        return nil, 0, fmt.Errorf("erro ao buscar jogadores: %w", err)
    }
    defer rows.Close()

    result := []PlayerResponse{}
    for rows.Next() {
        var p PlayerResponse
        err := rows.Scan(
            &p.ID, &p.Name,
            &p.Team.ID, &p.Team.Name, &p.Team.Flag,
        )
        if err != nil {
            return nil, 0, fmt.Errorf("erro ao ler jogador: %w", err)
        }
        result = append(result, p)
    }

    if err := rows.Err(); err != nil {
        return nil, 0, fmt.Errorf("erro ao iterar jogadores: %w", err)
    }

    return result, total, nil
}