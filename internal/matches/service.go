// internal/matches/service.go
package matches

import "context"

type Service interface {
    ListMatches(ctx context.Context, filters MatchFilters) (ListMatchesResponse, error)
    GetMatch(ctx context.Context, id string) (MatchResponse, error)
}

type service struct {
    repo Repository
}

func NewService(repo Repository) Service {
    return &service{repo: repo}
}

func (s *service) ListMatches(ctx context.Context, filters MatchFilters) (ListMatchesResponse, error) {
    matches, err := s.repo.FindAll(ctx, filters)
    if err != nil {
        return ListMatchesResponse{}, err
    }

    return ListMatchesResponse{
        Matches: matches,
        Total:   len(matches),
    }, nil
}

func (s *service) GetMatch(ctx context.Context, id string) (MatchResponse, error) {
    match, err := s.repo.FindByID(ctx, id)
    if err != nil {
        return MatchResponse{}, err
    }

    return match, nil
}