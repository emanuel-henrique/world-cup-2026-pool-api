// internal/players/model.go
package players

type Player struct {
    ID     string `json:"id"`
    Name   string `json:"name"`
    TeamID string `json:"team_id"`
}