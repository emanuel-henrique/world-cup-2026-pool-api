// internal/ranking/service.go
package ranking

import (
	"bolao-copa/internal/pagination"
	"context"
	"database/sql"
	"fmt"
)

type Service interface {
	GetRanking(ctx context.Context, page int, limit int) (RankingResponse, error)
}

type service struct {
	db *sql.DB
}

func NewService(db *sql.DB) Service {
	return &service{db: db}
}

func (s *service) GetRanking(ctx context.Context, page int, limit int) (RankingResponse, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	offset := (page - 1) * limit

	// Conta total de participantes
	var total int
	err := s.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM users`).Scan(&total)
	if err != nil {
		return RankingResponse{}, fmt.Errorf("erro ao contar usuários: %w", err)
	}

	query := `
		SELECT
			u.id,
			u.name,
			COALESCE(SUM(p.points), 0) + COALESCE(sp.points, 0) AS total_points,
			COUNT(p.id) FILTER (WHERE p.points = 10)             AS exact_scores
		FROM users u
		LEFT JOIN predictions p         ON p.user_id = u.id
		LEFT JOIN special_predictions sp ON sp.user_id = u.id
		GROUP BY u.id, u.name, sp.points
		ORDER BY total_points DESC, exact_scores DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := s.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return RankingResponse{}, fmt.Errorf("erro ao buscar ranking: %w", err)
	}
	defer rows.Close()

	entries := []RankingEntry{}
	position := offset + 1

	for rows.Next() {
		var e RankingEntry
		err := rows.Scan(&e.UserID, &e.Name, &e.TotalPoints, &e.ExactScores)
		if err != nil {
			return RankingResponse{}, fmt.Errorf("erro ao ler entrada do ranking: %w", err)
		}
		e.Position = position
		position++
		entries = append(entries, e)
	}

	if err := rows.Err(); err != nil {
		return RankingResponse{}, fmt.Errorf("erro ao iterar ranking: %w", err)
	}

	return RankingResponse{
		Entries: entries,
		Meta:    pagination.NewMeta(page, limit, total),
	}, nil
}