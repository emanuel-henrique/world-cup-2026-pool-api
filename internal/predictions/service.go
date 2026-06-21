// internal/predictions/service.go
package predictions

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Service interface {
	ListPredictions(ctx context.Context, userID string) (ListPredictionsResponse, error)
	UpsertPrediction(ctx context.Context, userID string, req UpsertPredictionRequest) error
	GetSpecial(ctx context.Context, userID string) (SpecialPredictionResponse, error)
	UpsertSpecial(ctx context.Context, userID string, req UpsertSpecialRequest) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) ListPredictions(ctx context.Context, userID string) (ListPredictionsResponse, error) {
	predictions, err := s.repo.FindByUser(ctx, userID)
	if err != nil {
		return ListPredictionsResponse{}, err
	}

	return ListPredictionsResponse{
		Predictions: predictions,
		Total:       len(predictions),
	}, nil
}

func (s *service) UpsertPrediction(ctx context.Context, userID string, req UpsertPredictionRequest) error {
	// RN-01 e RN-02 — valida kickoff no servidor
	kickoff, status, err := s.repo.FindMatchKickoff(ctx, req.MatchID)
	if err != nil {
		return err
	}

	if status != "scheduled" || time.Now().After(kickoff) {
		return ErrMatchLocked
	}

	p := Prediction{
		ID:        uuid.NewString(),
		UserID:    userID,
		MatchID:   req.MatchID,
		HomeScore: req.HomeScore,
		AwayScore: req.AwayScore,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	return s.repo.Upsert(ctx, p)
}

func (s *service) GetSpecial(ctx context.Context, userID string) (SpecialPredictionResponse, error) {
	sp, err := s.repo.FindSpecialByUser(ctx, userID)
	if err == ErrSpecialNotFound {
		// Retorna vazio em vez de erro — usuário ainda não fez palpite especial
		return SpecialPredictionResponse{}, nil
	}
	return sp, err
}

func (s *service) UpsertSpecial(ctx context.Context, userID string, req UpsertSpecialRequest) error {
	// RN-03 — bloqueia após início do primeiro jogo da Copa
	firstMatch, _, err := s.repo.FindMatchKickoff(ctx, "first")
	if err == nil && time.Now().After(firstMatch) {
		return ErrSpecialLocked
	}

	sp := SpecialPrediction{
		ID:          uuid.NewString(),
		UserID:      userID,
		ChampionID:  req.ChampionID,
		TopScorerID: req.TopScorerID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	return s.repo.UpsertSpecial(ctx, sp)
}