// internal/groups/service_test.go
package groups_test

import (
	"context"

	"bolao-copa/internal/groups"
)

// Mock do banco — simula o Service diretamente
type mockService struct{}

func (m *mockService) ListGroups(ctx context.Context) (groups.ListGroupsResponse, error) {
	return groups.ListGroupsResponse{
		Groups: []groups.Group{
			{
				Name: "C",
				Standings: []groups.GroupStanding{
					{TeamID: "BRA", TeamName: "Brasil", Flag: "https://flagcdn.com/w40/br.png", Points: 3},
					{TeamID: "MAR", TeamName: "Marrocos", Flag: "https://flagcdn.com/w40/ma.png", Points: 1},
					{TeamID: "HTI", TeamName: "Haiti", Flag: "https://flagcdn.com/w40/ht.png", Points: 0},
					{TeamID: "SCO", TeamName: "Escócia", Flag: "https://flagcdn.com/w40/gb-sct.png", Points: 0},
				},
			},
		},
	}, nil
}

func (m *mockService) GetGroup(ctx context.Context, name string) (groups.GroupDetailResponse, error) {
	if name != "C" {
		return groups.GroupDetailResponse{}, groups.ErrGroupNotFound
	}

	score1 := 1
	score0 := 0

	return groups.GroupDetailResponse{
		Group: groups.Group{
			Name: "C",
			Standings: []groups.GroupStanding{
				{TeamID: "BRA", TeamName: "Brasil", Points: 3},
				{TeamID: "MAR", TeamName: "Marrocos", Points: 1},
			},
		},
		Matches: []groups.GroupMatch{
			{
				ID:        "match-1",
				HomeTeam:  "Brasil",
				AwayTeam:  "Marrocos",
				HomeScore: &score1,
				AwayScore: &score0,
				Status:    "finished",
				KickoffAt: "2026-06-14T18:00:00Z",
			},
		},
	}, nil
}