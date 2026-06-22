// internal/matches/dto.go
package matches

import (
	"bolao-copa/internal/pagination"
	"time"
)

// --- Mapeamento da API Externa (API-Football) ---
// Adicionamos esta estrutura no fim ou topo do arquivo para receber os dados da API-Football
type APIFootballResponse struct {
	Response []struct {
		Teams struct {
			Home struct { Name string `json:"name"` } `json:"home"`
			Away struct { Name string `json:"name"` } `json:"away"`
		} `json:"teams"`
		Events []struct {
			Type   string `json:"type"` 
			Team struct {
				Name string `json:"name"`
			} `json:"team"`
			Time   struct {
				Elapsed int `json:"elapsed"` 
			} `json:"time"`
			Player struct {
				Name string `json:"name"` 
			} `json:"player"`
		} `json:"events"`
	} `json:"response"`
}

// --- Suas Estruturas Originais Mantidas ---

type ListMatchesResponse struct {
	Matches []MatchResponse `json:"matches"`
	Meta    pagination.Meta `json:"meta"`
}

type MatchResponse struct {
	ID           string       `json:"id"`
	HomeTeam     *TeamSummary `json:"home_team"`
	AwayTeam     *TeamSummary `json:"away_team"`
	HomeScore    *int         `json:"home_score"`
	AwayScore    *int         `json:"away_score"`
	HomeHalfTime *int         `json:"home_half_time"`
	AwayHalfTime *int         `json:"away_half_time"`
	Stage        string       `json:"stage"`
	GroupName    *string      `json:"group_name"`
	KickoffAt    time.Time    `json:"kickoff_at"`
	Status       string       `json:"status"`
	Goals        []Goal       `json:"goals"` // Usará a sua estrutura abaixo
	Minute       *int         `json:"minute"` 
}

type Goal struct {
	PlayerName string `json:"player_name"`
	Minute     int    `json:"minute"`
	Team       string `json:"team"` // "home" ou "away"
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