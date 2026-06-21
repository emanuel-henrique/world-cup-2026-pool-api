// internal/players/dto.go
package players

import "bolao-copa/internal/pagination"

type PlayerResponse struct {
    ID       string      `json:"id"`
    Name     string      `json:"name"`
    Team     TeamSummary `json:"team"`
}

type TeamSummary struct {
    ID   string `json:"id"`
    Name string `json:"name"`
    Flag string `json:"flag"`
}

type ListPlayersResponse struct {
    Players []PlayerResponse `json:"players"`
    Meta    pagination.Meta  `json:"meta"`
}