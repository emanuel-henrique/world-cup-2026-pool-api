// internal/groups/service.go
package groups

import (
	"context"
	"database/sql"
	"fmt"
)

type Service interface {
    ListGroups(ctx context.Context) (ListGroupsResponse, error)
    GetGroup(ctx context.Context, name string) (GroupDetailResponse, error)
}

type service struct {
    db *sql.DB
}

func NewService(db *sql.DB) Service {
    return &service{db: db}
}

func (s *service) ListGroups(ctx context.Context) (ListGroupsResponse, error) {
    query := `
        SELECT
            group_name, team_id, team_name, COALESCE(flag, ''),
            played, wins, draws, losses,
            goals_for, goals_against, points
        FROM group_standings
        ORDER BY group_name ASC, points DESC
    `

    rows, err := s.db.QueryContext(ctx, query)
    if err != nil {
        return ListGroupsResponse{}, fmt.Errorf("erro ao buscar grupos: %w", err)
    }
    defer rows.Close()

    groupMap := make(map[string]*Group)
    groupOrder := []string{}

    for rows.Next() {
        var groupName string
        var s GroupStanding

        err := rows.Scan(
            &groupName, &s.TeamID, &s.TeamName, &s.Flag,
            &s.Played, &s.Wins, &s.Draws, &s.Losses,
            &s.GoalsFor, &s.GoalsAgainst, &s.Points,
        )
        if err != nil {
            return ListGroupsResponse{}, fmt.Errorf("erro ao ler grupo: %w", err)
        }

        if _, exists := groupMap[groupName]; !exists {
            groupMap[groupName] = &Group{Name: groupName}
            groupOrder = append(groupOrder, groupName)
        }
        groupMap[groupName].Standings = append(groupMap[groupName].Standings, s)
    }

    if err := rows.Err(); err != nil {
        return ListGroupsResponse{}, fmt.Errorf("erro ao iterar grupos: %w", err)
    }

    var groups []Group
    for _, name := range groupOrder {
        groups = append(groups, *groupMap[name])
    }

    return ListGroupsResponse{Groups: groups}, nil
}

func (s *service) GetGroup(ctx context.Context, name string) (GroupDetailResponse, error) {
    // Busca standings do grupo
    standingsQuery := `
        SELECT
            team_id, team_name, COALESCE(flag, ''),
            played, wins, draws, losses,
            goals_for, goals_against, points
        FROM group_standings
        WHERE group_name = $1
        ORDER BY points DESC
    `

    rows, err := s.db.QueryContext(ctx, standingsQuery, name)
    if err != nil {
        return GroupDetailResponse{}, fmt.Errorf("erro ao buscar standings: %w", err)
    }
    defer rows.Close()

    var standings []GroupStanding
    for rows.Next() {
        var s GroupStanding
        err := rows.Scan(
            &s.TeamID, &s.TeamName, &s.Flag,
            &s.Played, &s.Wins, &s.Draws, &s.Losses,
            &s.GoalsFor, &s.GoalsAgainst, &s.Points,
        )
        if err != nil {
            return GroupDetailResponse{}, fmt.Errorf("erro ao ler standing: %w", err)
        }
        standings = append(standings, s)
    }

    if err := rows.Err(); err != nil {
        return GroupDetailResponse{}, fmt.Errorf("erro ao iterar standings: %w", err)
    }

    if len(standings) == 0 {
        return GroupDetailResponse{}, ErrGroupNotFound
    }

    // Busca jogos do grupo
    matchesQuery := `
        SELECT
            m.id, m.status, m.kickoff_at,
            m.home_score, m.away_score,
            COALESCE(ht.name, ''), COALESCE(ht.flag, ''),
            COALESCE(at.name, ''), COALESCE(at.flag, '')
        FROM matches m
        LEFT JOIN teams ht ON ht.id = m.home_team_id
        LEFT JOIN teams at ON at.id = m.away_team_id
        WHERE m.group_name = $1
        ORDER BY m.kickoff_at ASC
    `

    mrows, err := s.db.QueryContext(ctx, matchesQuery, name)
    if err != nil {
        return GroupDetailResponse{}, fmt.Errorf("erro ao buscar jogos do grupo: %w", err)
    }
    defer mrows.Close()

    var groupMatches []GroupMatch
    for mrows.Next() {
        var m GroupMatch
        err := mrows.Scan(
            &m.ID, &m.Status, &m.KickoffAt,
            &m.HomeScore, &m.AwayScore,
            &m.HomeTeam, &m.HomeFlag,
            &m.AwayTeam, &m.AwayFlag,
        )
        if err != nil {
            return GroupDetailResponse{}, fmt.Errorf("erro ao ler jogo do grupo: %w", err)
        }
        groupMatches = append(groupMatches, m)
    }

    if err := mrows.Err(); err != nil {
        return GroupDetailResponse{}, fmt.Errorf("erro ao iterar jogos: %w", err)
    }

    return GroupDetailResponse{
        Group:   Group{Name: name, Standings: standings},
        Matches: groupMatches,
    }, nil
}

var ErrGroupNotFound = fmt.Errorf("grupo não encontrado")