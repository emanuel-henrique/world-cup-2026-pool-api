// internal/matches/repository.go
package matches

import (
	"context"
	"database/sql"
	"fmt"
)

func derefStr(s *string) string {
	if s != nil {
		return *s
	}
	return ""
}

type Repository interface {
    FindAll(ctx context.Context, filters MatchFilters, limit, offset int) ([]MatchResponse, int, error)
    FindByID(ctx context.Context, id string) (MatchResponse, error)
    Upsert(ctx context.Context, match Match) error
    UpdateStatus(ctx context.Context, externalID string, status string, homeScore *int, awayScore *int) error
}

type postgresRepository struct {
    db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
    return &postgresRepository{db: db}
}

func (r *postgresRepository) FindAll(ctx context.Context, filters MatchFilters, limit, offset int) ([]MatchResponse, int, error) {
    var total int
    countQuery := `
        SELECT COUNT(*)
        FROM matches m
        WHERE ($1 = '' OR m.stage = $1)
          AND ($2 = '' OR m.status = $2)
          AND ($3 = '' OR m.group_name = $3)
    `
    err := r.db.QueryRowContext(ctx, countQuery, filters.Stage, filters.Status, filters.Group).Scan(&total)
    if err != nil {
        return nil, 0, fmt.Errorf("erro ao contar jogos: %w", err)
    }

    query := `
        SELECT
            m.id,
            m.home_score,
            m.away_score,
            m.stage,
            m.group_name,
            m.kickoff_at,
            m.status,
            m.minute,
            ht.id,   ht.name, COALESCE(ht.flag, ''),
            at.id,   at.name, COALESCE(at.flag, '')
        FROM matches m
        INNER JOIN teams ht ON ht.id = m.home_team_id
        INNER JOIN teams at ON at.id = m.away_team_id
        WHERE ($1 = '' OR m.stage = $1)
          AND ($2 = '' OR m.status = $2)
          AND ($3 = '' OR m.group_name = $3)
        ORDER BY m.kickoff_at ASC
        LIMIT $4 OFFSET $5
    `

    rows, err := r.db.QueryContext(ctx, query, filters.Stage, filters.Status, filters.Group, limit, offset)
    if err != nil {
        return nil, 0, fmt.Errorf("erro ao buscar jogos: %w", err)
    }
    defer rows.Close()

    matches := []MatchResponse{}
    matchIDs := []string{}
    for rows.Next() {
        var m MatchResponse
        var home, away TeamSummary
        var homeID, awayID *string
        var homeName, awayName *string
        var homeFlag, awayFlag *string

        err := rows.Scan(
            &m.ID,
            &m.HomeScore,
            &m.AwayScore,
            &m.Stage,
            &m.GroupName,
            &m.KickoffAt,
            &m.Status,
            &m.Minute,
            &homeID, &homeName, &homeFlag,
            &awayID, &awayName, &awayFlag,
        )
        if err != nil {
            return nil, 0, fmt.Errorf("erro ao ler jogo: %w", err)
        }

        if homeID != nil {
            home = TeamSummary{ID: *homeID, Name: derefStr(homeName), Flag: derefStr(homeFlag)}
            m.HomeTeam = &home
        }
        if awayID != nil {
            away = TeamSummary{ID: *awayID, Name: derefStr(awayName), Flag: derefStr(awayFlag)}
            m.AwayTeam = &away
        }

        matches = append(matches, m)
        matchIDs = append(matchIDs, m.ID)
    }

    if err := rows.Err(); err != nil {
        return nil, 0, fmt.Errorf("erro ao iterar jogos: %w", err)
    }

    // Buscar gols para todos os matches
    if len(matchIDs) > 0 {
        goalsQuery := `
            SELECT match_id, player_name, minute, team
            FROM goals
            WHERE match_id = ANY($1)
            ORDER BY minute ASC
        `
        goalRows, err := r.db.QueryContext(ctx, goalsQuery, matchIDs)
        if err != nil {
            return nil, 0, fmt.Errorf("erro ao buscar gols: %w", err)
        }
        defer goalRows.Close()

        goalsByMatch := make(map[string][]Goal)
        for goalRows.Next() {
            var g Goal
            var matchID string
            err := goalRows.Scan(&matchID, &g.PlayerName, &g.Minute, &g.Team)
            if err != nil {
                return nil, 0, fmt.Errorf("erro ao ler gol: %w", err)
            }
            goalsByMatch[matchID] = append(goalsByMatch[matchID], g)
        }

        // Adicionar gols aos matches
        for i := range matches {
            if goals, ok := goalsByMatch[matches[i].ID]; ok {
                matches[i].Goals = goals
            } else {
                matches[i].Goals = []Goal{}
            }
        }
    }

    return matches, total, nil
}

func (r *postgresRepository) FindByID(ctx context.Context, id string) (MatchResponse, error) {
    query := `
        SELECT
            m.id,
            m.home_score,
            m.away_score,
            m.stage,
            m.group_name,
            m.kickoff_at,
            m.status,
            m.minute,
            ht.id,   ht.name, COALESCE(ht.flag, ''),
            at.id,   at.name, COALESCE(at.flag, '')
        FROM matches m
        LEFT JOIN teams ht ON ht.id = m.home_team_id
        LEFT JOIN teams at ON at.id = m.away_team_id
        WHERE m.id = $1
    `

    var m MatchResponse
    var home, away TeamSummary
    var homeID, awayID *string
    var homeName, awayName *string
    var homeFlag, awayFlag *string

    err := r.db.QueryRowContext(ctx, query, id).Scan(
        &m.ID,
        &m.HomeScore,
        &m.AwayScore,
        &m.Stage,
        &m.GroupName,
        &m.KickoffAt,
        &m.Status,
        &m.Minute,
        &homeID, &homeName, &homeFlag,
        &awayID, &awayName, &awayFlag,
    )
    if err == sql.ErrNoRows {
        return MatchResponse{}, ErrMatchNotFound
    }
    if err != nil {
        return MatchResponse{}, fmt.Errorf("erro ao buscar jogo: %w", err)
    }

    if homeID != nil {
        home = TeamSummary{ID: *homeID, Name: derefStr(homeName), Flag: derefStr(homeFlag)}
        m.HomeTeam = &home
    }
    if awayID != nil {
        away = TeamSummary{ID: *awayID, Name: derefStr(awayName), Flag: derefStr(awayFlag)}
        m.AwayTeam = &away
    }

    // Buscar gols para este match
    goalsQuery := `
        SELECT player_name, minute, team
        FROM goals
        WHERE match_id = $1
        ORDER BY minute ASC
    `
    goalRows, err := r.db.QueryContext(ctx, goalsQuery, id)
    if err == nil {
        defer goalRows.Close()
        var goals []Goal
        for goalRows.Next() {
            var g Goal
            err := goalRows.Scan(&g.PlayerName, &g.Minute, &g.Team)
            if err == nil {
                goals = append(goals, g)
            }
        }
        m.Goals = goals
    } else {
        m.Goals = []Goal{}
    }

    return m, nil
}

func (r *postgresRepository) Upsert(ctx context.Context, match Match) error {
    query := `
        INSERT INTO matches (
            id, external_id, home_team_id, away_team_id,
            home_score, away_score, stage, group_name, kickoff_at, status, minute
        ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
        ON CONFLICT (external_id) DO UPDATE SET
            home_team_id = EXCLUDED.home_team_id,
            away_team_id = EXCLUDED.away_team_id,
            home_score   = EXCLUDED.home_score,
            away_score   = EXCLUDED.away_score,
            status       = EXCLUDED.status,
            minute       = EXCLUDED.minute
    `

    _, err := r.db.ExecContext(ctx, query,
        match.ID,
        match.ExternalID,
        match.HomeTeamID,
        match.AwayTeamID,
        match.HomeScore,
        match.AwayScore,
        match.Stage,
        match.GroupName,
        match.KickoffAt,
        match.Status,
        match.Minute,
    )
    if err != nil {
        return fmt.Errorf("erro ao upsert jogo: %w", err)
    }

    return nil
}

func (r *postgresRepository) UpdateStatus(ctx context.Context, externalID string, status string, homeScore *int, awayScore *int) error {
    query := `
        UPDATE matches
        SET status     = $1,
            home_score = $2,
            away_score = $3
        WHERE external_id = $4
    `

    _, err := r.db.ExecContext(ctx, query, status, homeScore, awayScore, externalID)
    if err != nil {
        return fmt.Errorf("erro ao atualizar status do jogo: %w", err)
    }

    return nil
}

var ErrMatchNotFound = fmt.Errorf("jogo não encontrado")