// internal/ranking/dto.go
package ranking

type RankingResponse struct {
	Entries []RankingEntry `json:"entries"`
	Total   int            `json:"total"`
	Page    int            `json:"page"`
	Limit   int            `json:"limit"`
}