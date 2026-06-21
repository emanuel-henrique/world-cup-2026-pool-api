// internal/predictions/service_test.go
package predictions_test

import (
	"context"
	"testing"
	"time"

	"bolao-copa/internal/predictions"
)

// Mock do Repository
type mockRepository struct {
	predictions []predictions.PredictionResponse
	special     *predictions.SpecialPredictionResponse
	kickoff     time.Time
	status      string
}

func newMockRepository() *mockRepository {
	return &mockRepository{
		kickoff: time.Now().Add(24 * time.Hour), // jogo amanhã por padrão
		status:  "scheduled",
	}
}

func (m *mockRepository) FindByUser(ctx context.Context, userID string) ([]predictions.PredictionResponse, error) {
	return m.predictions, nil
}

func (m *mockRepository) Upsert(ctx context.Context, p predictions.Prediction) error {
	m.predictions = append(m.predictions, predictions.PredictionResponse{
		ID:        p.ID,
		MatchID:   p.MatchID,
		HomeScore: p.HomeScore,
		AwayScore: p.AwayScore,
	})
	return nil
}

func (m *mockRepository) FindSpecialByUser(ctx context.Context, userID string) (predictions.SpecialPredictionResponse, error) {
	if m.special == nil {
		return predictions.SpecialPredictionResponse{}, predictions.ErrSpecialNotFound
	}
	return *m.special, nil
}

func (m *mockRepository) UpsertSpecial(ctx context.Context, sp predictions.SpecialPrediction) error {
	champion := sp.ChampionID
	scorer   := sp.TopScorerID
	m.special = &predictions.SpecialPredictionResponse{
		ChampionID:  champion,
		TopScorerID: scorer,
	}
	return nil
}

func (m *mockRepository) FindMatchKickoff(ctx context.Context, matchID string) (time.Time, string, error) {
	if matchID == "nao-existe" {
		return time.Time{}, "", predictions.ErrMatchNotFound
	}
	return m.kickoff, m.status, nil
}

// Testes do UpsertPrediction

func TestUpsertPrediction_Success(t *testing.T) {
	repo    := newMockRepository()
	service := predictions.NewService(repo)

	err := service.UpsertPrediction(context.Background(), "user-1", predictions.UpsertPredictionRequest{
		MatchID:   "match-1",
		HomeScore: 2,
		AwayScore: 1,
	})

	if err != nil {
		t.Fatalf("esperava sucesso, got erro: %v", err)
	}
	if len(repo.predictions) != 1 {
		t.Fatalf("esperava 1 palpite salvo, got %d", len(repo.predictions))
	}
}

func TestUpsertPrediction_MatchLocked_Live(t *testing.T) {
	repo         := newMockRepository()
	repo.status   = "live"
	service      := predictions.NewService(repo)

	err := service.UpsertPrediction(context.Background(), "user-1", predictions.UpsertPredictionRequest{
		MatchID:   "match-1",
		HomeScore: 1,
		AwayScore: 0,
	})

	if err != predictions.ErrMatchLocked {
		t.Fatalf("esperava ErrMatchLocked, got %v", err)
	}
}

func TestUpsertPrediction_MatchLocked_KickoffPassed(t *testing.T) {
	repo          := newMockRepository()
	repo.kickoff   = time.Now().Add(-1 * time.Hour) // kickoff já passou
	service       := predictions.NewService(repo)

	err := service.UpsertPrediction(context.Background(), "user-1", predictions.UpsertPredictionRequest{
		MatchID:   "match-1",
		HomeScore: 1,
		AwayScore: 0,
	})

	if err != predictions.ErrMatchLocked {
		t.Fatalf("esperava ErrMatchLocked, got %v", err)
	}
}

func TestUpsertPrediction_MatchNotFound(t *testing.T) {
	repo    := newMockRepository()
	service := predictions.NewService(repo)

	err := service.UpsertPrediction(context.Background(), "user-1", predictions.UpsertPredictionRequest{
		MatchID:   "nao-existe",
		HomeScore: 1,
		AwayScore: 0,
	})

	if err != predictions.ErrMatchNotFound {
		t.Fatalf("esperava ErrMatchNotFound, got %v", err)
	}
}

func TestListPredictions_Success(t *testing.T) {
	repo    := newMockRepository()
	service := predictions.NewService(repo)

	// Salva um palpite primeiro
	_ = service.UpsertPrediction(context.Background(), "user-1", predictions.UpsertPredictionRequest{
		MatchID:   "match-1",
		HomeScore: 2,
		AwayScore: 1,
	})

	resp, err := service.ListPredictions(context.Background(), "user-1")
	if err != nil {
		t.Fatalf("esperava sucesso, got erro: %v", err)
	}
	if resp.Total != 1 {
		t.Fatalf("esperava 1 palpite, got %d", resp.Total)
	}
}

// Testes do UpsertSpecial

func TestUpsertSpecial_Success(t *testing.T) {
	repo    := newMockRepository()
	service := predictions.NewService(repo)

	champion := "BRA"
	scorer   := "player-1"

	err := service.UpsertSpecial(context.Background(), "user-1", predictions.UpsertSpecialRequest{
		ChampionID:  &champion,
		TopScorerID: &scorer,
	})

	if err != nil {
		t.Fatalf("esperava sucesso, got erro: %v", err)
	}
	if repo.special == nil {
		t.Fatal("esperava palpite especial salvo")
	}
	if *repo.special.ChampionID != "BRA" {
		t.Fatalf("esperava BRA, got %s", *repo.special.ChampionID)
	}
}

func TestGetSpecial_NotFound_ReturnsEmpty(t *testing.T) {
	repo    := newMockRepository()
	service := predictions.NewService(repo)

	resp, err := service.GetSpecial(context.Background(), "user-1")
	if err != nil {
		t.Fatalf("esperava sucesso mesmo sem palpite, got erro: %v", err)
	}
	if resp.ChampionID != nil {
		t.Fatal("esperava ChampionID nil")
	}
}