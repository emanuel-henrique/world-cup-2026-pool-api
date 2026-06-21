// internal/groups/model.go
package groups

type GroupStanding struct {
    TeamID       string `json:"team_id"`
    TeamName     string `json:"team_name"`
    Flag         string `json:"flag"`
    Played       int    `json:"played"`
    Wins         int    `json:"wins"`
    Draws        int    `json:"draws"`
    Losses       int    `json:"losses"`
    GoalsFor     int    `json:"goals_for"`
    GoalsAgainst int    `json:"goals_against"`
    Points       int    `json:"points"`
}

type Group struct {
    Name      string          `json:"name"`
    Standings []GroupStanding `json:"standings"`
}