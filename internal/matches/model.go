// internal/matches/model.go
package matches

import "time"

type Match struct {
    ID         string    `json:"id"`
    ExternalID string    `json:"external_id"`
    HomeTeamID *string   `json:"home_team_id"`
    AwayTeamID *string   `json:"away_team_id"`
    HomeScore  *int      `json:"home_score"`
    AwayScore  *int      `json:"away_score"`
    HomeHalfTime *int      `json:"home_half_time"`
    AwayHalfTime *int      `json:"away_half_time"`
    Stage      string    `json:"stage"`
    GroupName  *string   `json:"group_name"`
    KickoffAt  time.Time `json:"kickoff_at"`
    Status     string    `json:"status"`
    CreatedAt  time.Time `json:"created_at"`
    Goals      []Goal    `json:"goals"`
    Minute     *int      `json:"minute"` // current minute for live matches
}

type Status string

const (
    StatusScheduled Status = "scheduled"
    StatusLive      Status = "live"
    StatusFinished  Status = "finished"
)

type Stage string

const (
    StageGroup    Stage = "group"
    StageRound32  Stage = "round_of_32"
    StageRound16  Stage = "round_of_16"
    StageQuarter  Stage = "quarter"
    StageSemi     Stage = "semi"
    StageFinal    Stage = "final"
)