// internal/predictions/model.go
package predictions

import "time"

type Prediction struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	MatchID   string    `json:"match_id"`
	HomeScore int       `json:"home_score"`
	AwayScore int       `json:"away_score"`
	Points    int       `json:"points"`
	Scored    bool      `json:"scored"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type SpecialPrediction struct {
	ID          string    `json:"id"`
	UserID      string    `json:"user_id"`
	ChampionID  *string   `json:"champion_id"`
	TopScorerID *string   `json:"top_scorer_id"`
	Points      int       `json:"points"`
	Scored      bool      `json:"scored"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}