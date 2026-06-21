// internal/matches/service.go
package matches

import (
	"bolao-copa/internal/pagination"
	"context"
)

type Service interface {
    ListMatches(ctx context.Context, filters MatchFilters, page, limit int) (ListMatchesResponse, error)
    GetMatch(ctx context.Context, id string) (MatchResponse, error)
}

type service struct {
    repo Repository
}

func NewService(repo Repository) Service {
    return &service{repo: repo}
}

func (s *service) ListMatches(ctx context.Context, filters MatchFilters, page, limit int) (ListMatchesResponse, error) {
    offset := pagination.GetOffset(page, limit)
    matches, total, err := s.repo.FindAll(ctx, filters, limit, offset)
    if err != nil {
        return ListMatchesResponse{}, err
    }

    return ListMatchesResponse{
        Matches: matches,
        Meta:    pagination.NewMeta(page, limit, total),
    }, nil
}

func (s *service) GetMatch(ctx context.Context, id string) (MatchResponse, error) {
    match, err := s.repo.FindByID(ctx, id)
    if err != nil {
        return MatchResponse{}, err
    }

    return match, nil
}