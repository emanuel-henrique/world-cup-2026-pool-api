// internal/scoring/scoring.go
package scoring

import (
	"context"
	"database/sql"
	"fmt"
)

type MatchResult struct {
	HomeScore int
	AwayScore int
}

func CalculatePoints(prediction MatchResult, result MatchResult) int {
	// Placar exato
	if prediction.HomeScore == result.HomeScore && prediction.AwayScore == result.AwayScore {
		return 10
	}

	predWinner := winner(prediction)
	realWinner := winner(result)

	// Vencedores diferentes
	if predWinner != realWinner {
		return 0
	}

	// Empate correto (placar errado — já coberto acima)
	if predWinner == "draw" {
		return 5
	}

	// Vencedor correto — verifica margem de vitória
	predMargin := abs(prediction.HomeScore - prediction.AwayScore)
	realMargin := abs(result.HomeScore - result.AwayScore)

	if predMargin == realMargin {
		return 7
	}

	return 5
}

func winner(r MatchResult) string {
	if r.HomeScore > r.AwayScore {
		return "home"
	}
	if r.AwayScore > r.HomeScore {
		return "away"
	}
	return "draw"
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func CalculateMatch(ctx context.Context, db *sql.DB, matchID string) error {
	var homeScore, awayScore int
	err := db.QueryRowContext(ctx, `
		SELECT home_score, away_score FROM matches WHERE id = $1
	`, matchID).Scan(&homeScore, &awayScore)
	if err != nil {
		return fmt.Errorf("erro ao buscar resultado do jogo: %w", err)
	}

	result := MatchResult{HomeScore: homeScore, AwayScore: awayScore}

	rows, err := db.QueryContext(ctx, `
		SELECT id, home_score, away_score
		FROM predictions
		WHERE match_id = $1 AND scored = FALSE
	`, matchID)
	if err != nil {
		return fmt.Errorf("erro ao buscar palpites: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var predID string
		var predHome, predAway int

		if err := rows.Scan(&predID, &predHome, &predAway); err != nil {
			return fmt.Errorf("erro ao ler palpite: %w", err)
		}

		points := CalculatePoints(
			MatchResult{HomeScore: predHome, AwayScore: predAway},
			result,
		)

		_, err := db.ExecContext(ctx, `
			UPDATE predictions
			SET points = $1, scored = TRUE
			WHERE id = $2
		`, points, predID)
		if err != nil {
			return fmt.Errorf("erro ao atualizar pontos: %w", err)
		}
	}

	return rows.Err()
}

func CalculateSpecial(ctx context.Context, db *sql.DB, championID string, topScorerID string) error {
	_, err := db.ExecContext(ctx, `
		UPDATE special_predictions
		SET points = points + 30, scored = TRUE
		WHERE champion_id = $1 AND scored = FALSE
	`, championID)
	if err != nil {
		return fmt.Errorf("erro ao pontuar campeão: %w", err)
	}

	_, err = db.ExecContext(ctx, `
		UPDATE special_predictions
		SET points = points + 20
		WHERE top_scorer_id = $1 AND scored = FALSE
	`, topScorerID)
	if err != nil {
		return fmt.Errorf("erro ao pontuar artilheiro: %w", err)
	}

	return nil
}