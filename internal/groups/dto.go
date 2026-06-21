// internal/groups/dto.go
package groups

type ListGroupsResponse struct {
    Groups []Group `json:"groups"`
}

type GroupDetailResponse struct {
    Group     Group         `json:"group"`
    Matches   []GroupMatch  `json:"matches"`
}

type GroupMatch struct {
    ID        string  `json:"id"`
    HomeTeam  string  `json:"home_team"`
    AwayTeam  string  `json:"away_team"`
    HomeFlag  string  `json:"home_flag"`
    AwayFlag  string  `json:"away_flag"`
    HomeScore *int    `json:"home_score"`
    AwayScore *int    `json:"away_score"`
    Status    string  `json:"status"`
    KickoffAt string  `json:"kickoff_at"`
}