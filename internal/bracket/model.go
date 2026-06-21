// internal/bracket/model.go
package bracket

type BracketMatch struct {
	ID        string       `json:"id"`
	HomeTeam  *TeamSummary `json:"home_team"`
	AwayTeam  *TeamSummary `json:"away_team"`
	HomeScore *int         `json:"home_score"`
	AwayScore *int         `json:"away_score"`
	Status    string       `json:"status"`
	KickoffAt string       `json:"kickoff_at"`
	Winner    *TeamSummary `json:"winner"`
}

type TeamSummary struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Flag string `json:"flag"`
}

type BracketRound struct {
	Name    string         `json:"name"`
	Order   int            `json:"order"`
	Matches []BracketMatch `json:"matches"`
}