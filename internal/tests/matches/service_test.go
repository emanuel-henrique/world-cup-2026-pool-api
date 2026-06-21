// internal/matches/service_test.go
package matches_test

import (
	"context"
	"testing"

	"bolao-copa/internal/matches"
)

// Mock do Repository
type mockRepository struct {
	data []matches.MatchResponse
}

func newMockRepository() *mockRepository {
	score0 := 0
	score1 := 1
	group := "C"

	return &mockRepository{
		data: []matches.MatchResponse{
			{
				ID:        "match-1",
				HomeTeam:  &matches.TeamSummary{ID: "BRA", Name: "Brasil", Flag: "https://flagcdn.com/w40/br.png"},
				AwayTeam:  &matches.TeamSummary{ID: "MAR", Name: "Marrocos", Flag: "https://flagcdn.com/w40/ma.png"},
				HomeScore: &score0,
				AwayScore: &score1,
				Stage:     string(matches.StageGroup),
				GroupName: &group,
				Status:    string(matches.StatusFinished),
			},
			{
					ID:       "match-2",
					HomeTeam: &matches.TeamSummary{ID: "ARG", Name: "Argentina", Flag: "https://flagcdn.com/w40/ar.png"},
					AwayTeam: &matches.TeamSummary{ID: "USA", Name: "Estados Unidos", Flag: "https://flagcdn.com/w40/us.png"},
					Stage:    string(matches.StageGroup),
					Status:   string(matches.StatusScheduled),
			},
		},
	}
}

func (m *mockRepository) FindAll(ctx context.Context, filters matches.MatchFilters) ([]matches.MatchResponse, error) {
	if filters.Status == "" && filters.Stage == "" && filters.Group == "" {
		return m.data, nil
	}

	var result []matches.MatchResponse
	for _, match := range m.data {
		if filters.Status != "" && match.Status != filters.Status {
			continue
		}
		if filters.Stage != "" && match.Stage != filters.Stage {
			continue
		}
		result = append(result, match)
	}
	return result, nil
}

func (m *mockRepository) FindByID(ctx context.Context, id string) (matches.MatchResponse, error) {
	for _, match := range m.data {
		if match.ID == id {
			return match, nil
		}
	}
	return matches.MatchResponse{}, matches.ErrMatchNotFound
}

func (m *mockRepository) Upsert(ctx context.Context, match matches.Match) error {
	return nil
}

func (m *mockRepository) UpdateStatus(ctx context.Context, externalID string, status string, homeScore *int, awayScore *int) error {
	return nil
}

// Testes

func TestListMatches_NoFilter(t *testing.T) {
	repo    := newMockRepository()
	service := matches.NewService(repo)

	resp, err := service.ListMatches(context.Background(), matches.MatchFilters{})
	if err != nil {
		t.Fatalf("esperava sucesso, got erro: %v", err)
	}
	if resp.Total != 2 {
		t.Fatalf("esperava 2 jogos, got %d", resp.Total)
	}
}

func TestListMatches_FilterByStatus(t *testing.T) {
	repo    := newMockRepository()
	service := matches.NewService(repo)

	resp, err := service.ListMatches(context.Background(), matches.MatchFilters{
		Status: string(matches.StatusFinished),
	})
	if err != nil {
		t.Fatalf("esperava sucesso, got erro: %v", err)
	}
	if resp.Total != 1 {
		t.Fatalf("esperava 1 jogo finalizado, got %d", resp.Total)
	}
}

func TestGetMatch_Found(t *testing.T) {
	repo    := newMockRepository()
	service := matches.NewService(repo)

	match, err := service.GetMatch(context.Background(), "match-1")
	if err != nil {
		t.Fatalf("esperava sucesso, got erro: %v", err)
	}
	if match.ID != "match-1" {
		t.Fatalf("esperava match-1, got %s", match.ID)
	}
}

func TestGetMatch_NotFound(t *testing.T) {
	repo    := newMockRepository()
	service := matches.NewService(repo)

	_, err := service.GetMatch(context.Background(), "nao-existe")
	if err != matches.ErrMatchNotFound {
		t.Fatalf("esperava ErrMatchNotFound, got %v", err)
	}
}