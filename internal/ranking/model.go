// internal/ranking/model.go
package ranking

type RankingEntry struct {
	Position    int    `json:"position"`
	UserID      string `json:"user_id"`
	Name        string `json:"name"`
	TotalPoints int    `json:"total_points"`
	ExactScores int    `json:"exact_scores"`
}
