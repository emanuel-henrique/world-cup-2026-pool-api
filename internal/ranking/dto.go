// internal/ranking/dto.go
package ranking

import "bolao-copa/internal/pagination"

type RankingResponse struct {
	Entries []RankingEntry  `json:"entries"`
	Meta    pagination.Meta `json:"meta"`
}