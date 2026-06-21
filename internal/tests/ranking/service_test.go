// internal/ranking/service_test.go
package ranking_test

import (
	"context"

	"bolao-copa/internal/ranking"
)

// Mock do Service
type mockService struct{}

func (m *mockService) GetRanking(ctx context.Context, page int, limit int) (ranking.RankingResponse, error) {
	return ranking.RankingResponse{
		Total: 3,
		Page:  page,
		Limit: limit,
		Entries: []ranking.RankingEntry{
			{Position: 1, UserID: "user-1", Name: "Emanuel", TotalPoints: 45, ExactScores: 3},
			{Position: 2, UserID: "user-2", Name: "João", TotalPoints: 38, ExactScores: 2},
			{Position: 3, UserID: "user-3", Name: "Maria", TotalPoints: 30, ExactScores: 1},
		},
	}, nil
}