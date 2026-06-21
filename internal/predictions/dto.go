// internal/predictions/dto.go
package predictions

import "time"

// Requests
type UpsertPredictionRequest struct {
	MatchID   string `json:"match_id"   binding:"required"`
	HomeScore int    `json:"home_score"  binding:"min=0"`
	AwayScore int    `json:"away_score"  binding:"min=0"`
}

type UpsertSpecialRequest struct {
	ChampionID  *string `json:"champion_id"`
	TopScorerID *string `json:"top_scorer_id"`
}

// Responses
type PredictionResponse struct {
	ID        string    `json:"id"`
	MatchID   string    `json:"match_id"`
	HomeScore int       `json:"home_score"`
	AwayScore int       `json:"away_score"`
	Points    int       `json:"points"`
	Scored    bool      `json:"scored"`
	UpdatedAt time.Time `json:"updated_at"`
}

type SpecialPredictionResponse struct {
	ChampionID  *string `json:"champion_id"`
	TopScorerID *string `json:"top_scorer_id"`
	Points      int     `json:"points"`
	Scored      bool    `json:"scored"`
}

type ListPredictionsResponse struct {
	Predictions []PredictionResponse `json:"predictions"`
	Total       int                  `json:"total"`
}