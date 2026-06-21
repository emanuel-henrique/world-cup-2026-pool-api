// internal/predictions/repository.go
package predictions

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type Repository interface {
	FindByUser(ctx context.Context, userID string, limit, offset int) ([]PredictionResponse, int, error)
	Upsert(ctx context.Context, p Prediction) error
	FindSpecialByUser(ctx context.Context, userID string) (SpecialPredictionResponse, error)
	UpsertSpecial(ctx context.Context, sp SpecialPrediction) error
	FindMatchKickoff(ctx context.Context, matchID string) (time.Time, string, error)
}

type postgresRepository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &postgresRepository{db: db}
}

func (r *postgresRepository) FindByUser(ctx context.Context, userID string, limit, offset int) ([]PredictionResponse, int, error) {
	var total int
	err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM predictions WHERE user_id = $1", userID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("erro ao contar palpites: %w", err)
	}

	query := `
		SELECT id, match_id, home_score, away_score, points, scored, updated_at
		FROM predictions
		WHERE user_id = $1
		ORDER BY updated_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.QueryContext(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("erro ao buscar palpites: %w", err)
	}
	defer rows.Close()

	var result []PredictionResponse
	for rows.Next() {
		var p PredictionResponse
		err := rows.Scan(&p.ID, &p.MatchID, &p.HomeScore, &p.AwayScore, &p.Points, &p.Scored, &p.UpdatedAt)
		if err != nil {
			return nil, 0, fmt.Errorf("erro ao ler palpite: %w", err)
		}
		result = append(result, p)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("erro ao iterar palpites: %w", err)
	}

	return result, total, nil
}

func (r *postgresRepository) Upsert(ctx context.Context, p Prediction) error {
	query := `
		INSERT INTO predictions (id, user_id, match_id, home_score, away_score, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (user_id, match_id) DO UPDATE SET
			home_score = EXCLUDED.home_score,
			away_score = EXCLUDED.away_score,
			updated_at = EXCLUDED.updated_at
	`

	_, err := r.db.ExecContext(ctx, query,
		p.ID, p.UserID, p.MatchID,
		p.HomeScore, p.AwayScore,
		p.CreatedAt, p.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("erro ao salvar palpite: %w", err)
	}

	return nil
}

func (r *postgresRepository) FindSpecialByUser(ctx context.Context, userID string) (SpecialPredictionResponse, error) {
	query := `
		SELECT champion_id, top_scorer_id, points, scored
		FROM special_predictions
		WHERE user_id = $1
	`

	var sp SpecialPredictionResponse
	err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&sp.ChampionID, &sp.TopScorerID, &sp.Points, &sp.Scored,
	)
	if err == sql.ErrNoRows {
		return SpecialPredictionResponse{}, ErrSpecialNotFound
	}
	if err != nil {
		return SpecialPredictionResponse{}, fmt.Errorf("erro ao buscar palpite especial: %w", err)
	}

	return sp, nil
}

func (r *postgresRepository) UpsertSpecial(ctx context.Context, sp SpecialPrediction) error {
	query := `
		INSERT INTO special_predictions (id, user_id, champion_id, top_scorer_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (user_id) DO UPDATE SET
			champion_id   = EXCLUDED.champion_id,
			top_scorer_id = EXCLUDED.top_scorer_id,
			updated_at    = EXCLUDED.updated_at
	`

	_, err := r.db.ExecContext(ctx, query,
		sp.ID, sp.UserID, sp.ChampionID, sp.TopScorerID,
		sp.CreatedAt, sp.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("erro ao salvar palpite especial: %w", err)
	}

	return nil
}

func (r *postgresRepository) FindMatchKickoff(ctx context.Context, matchID string) (time.Time, string, error) {
	query := `SELECT kickoff_at, status FROM matches WHERE id = $1`

	var kickoff time.Time
	var status string
	err := r.db.QueryRowContext(ctx, query, matchID).Scan(&kickoff, &status)
	if err == sql.ErrNoRows {
		return time.Time{}, "", ErrMatchNotFound
	}
	if err != nil {
		return time.Time{}, "", fmt.Errorf("erro ao buscar jogo: %w", err)
	}

	return kickoff, status, nil
}

var (
	ErrSpecialNotFound  = fmt.Errorf("palpite especial não encontrado")
	ErrMatchNotFound    = fmt.Errorf("jogo não encontrado")
	ErrMatchLocked      = fmt.Errorf("palpite encerrado — jogo já começou")
	ErrSpecialLocked    = fmt.Errorf("palpite especial encerrado — copa já começou")
)