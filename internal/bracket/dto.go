// internal/bracket/dto.go
package bracket

type BracketResponse struct {
	Rounds []BracketRound `json:"rounds"`
}