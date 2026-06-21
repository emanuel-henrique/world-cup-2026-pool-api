// internal/matches/dto.go
package matches

import (
	"bolao-copa/internal/pagination"
	"time"
)

type ListMatchesResponse struct {
    Matches []MatchResponse `json:"matches"`
    Meta    pagination.Meta `json:"meta"`
}

type MatchResponse struct {
    ID         string      `json:"id"`
    HomeTeam   *TeamSummary `json:"home_team"`
    AwayTeam   *TeamSummary `json:"away_team"`
    HomeScore  *int         `json:"home_score"`
    AwayScore  *int         `json:"away_score"`
    Stage      string       `json:"stage"`
    GroupName  *string      `json:"group_name"`
    KickoffAt  time.Time    `json:"kickoff_at"`
    Status     string       `json:"status"`
}

type TeamSummary struct {
    ID   string `json:"id"`
    Name string `json:"name"`
    Flag string `json:"flag"`
}

type MatchFilters struct {
    Stage  string
    Status string
    Group  string
}