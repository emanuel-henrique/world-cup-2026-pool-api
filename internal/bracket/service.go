// internal/bracket/service.go
package bracket

import (
	"context"
	"database/sql"
	"fmt"
)

type Service interface {
	GetBracket(ctx context.Context) (BracketResponse, error)
}

type service struct {
	db *sql.DB
}

func NewService(db *sql.DB) Service {
	return &service{db: db}
}

var stageOrder = map[string]int{
	"round_of_32": 1,
	"round_of_16": 2,
	"quarter":     3,
	"semi":        4,
	"final":       5,
}

var stageNames = map[string]string{
	"round_of_32": "32avos de Final",
	"round_of_16": "Oitavas de Final",
	"quarter":     "Quartas de Final",
	"semi":        "Semifinais",
	"final":       "Final",
}

func (s *service) GetBracket(ctx context.Context) (BracketResponse, error) {
	query := `
		SELECT
			m.id, m.status, m.kickoff_at,
			m.home_score, m.away_score,
			m.stage,
			ht.id, ht.name, ht.flag,
			at.id, at.name, at.flag
		FROM matches m
		LEFT JOIN teams ht ON ht.id = m.home_team_id
		LEFT JOIN teams at ON at.id = m.away_team_id
		WHERE m.stage != 'group'
		ORDER BY m.stage ASC, m.kickoff_at ASC
	`

	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return BracketResponse{}, fmt.Errorf("erro ao buscar chaveamento: %w", err)
	}
	defer rows.Close()

	roundMap   := make(map[string]*BracketRound)
	roundOrder := []string{}

	for rows.Next() {
		var m BracketMatch
		var stage string
		var homeID, homeName, homeFlag *string
		var awayID, awayName, awayFlag *string

		err := rows.Scan(
			&m.ID, &m.Status, &m.KickoffAt,
			&m.HomeScore, &m.AwayScore,
			&stage,
			&homeID, &homeName, &homeFlag,
			&awayID, &awayName, &awayFlag,
		)
		if err != nil {
			return BracketResponse{}, fmt.Errorf("erro ao ler jogo: %w", err)
		}

		if homeID != nil {
			m.HomeTeam = &TeamSummary{ID: *homeID, Name: *homeName, Flag: *homeFlag}
		}
		if awayID != nil {
			m.AwayTeam = &TeamSummary{ID: *awayID, Name: *awayName, Flag: *awayFlag}
		}

		m.Winner = CalculateWinner(m)

		if _, exists := roundMap[stage]; !exists {
			roundMap[stage] = &BracketRound{
				Name:  stageNames[stage],
				Order: stageOrder[stage],
			}
			roundOrder = append(roundOrder, stage)
		}
		roundMap[stage].Matches = append(roundMap[stage].Matches, m)
	}

	if err := rows.Err(); err != nil {
		return BracketResponse{}, fmt.Errorf("erro ao iterar chaveamento: %w", err)
	}

	rounds := []BracketRound{}
	for _, stage := range roundOrder {
		rounds = append(rounds, *roundMap[stage])
	}

	return BracketResponse{Rounds: rounds}, nil
}

func CalculateWinner(m BracketMatch) *TeamSummary {
	if m.Status != "finished" || m.HomeScore == nil || m.AwayScore == nil {
		return nil
	}
	if *m.HomeScore > *m.AwayScore {
		return m.HomeTeam
	}
	if *m.AwayScore > *m.HomeScore {
		return m.AwayTeam
	}
	return nil
}

func GetStageOrder(stage string) int {
	return stageOrder[stage]
}