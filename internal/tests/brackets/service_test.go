// internal/bracket/service_test.go
package bracket_test

import (
	"testing"

	"bolao-copa/internal/bracket"
)

func TestGetBracket_WinnerCalculation(t *testing.T) {
	score2 := 2
	score1 := 1

	// Testa a lógica de winner diretamente nas structs
	match := bracket.BracketMatch{
		HomeTeam:  &bracket.TeamSummary{ID: "BRA", Name: "Brasil"},
		AwayTeam:  &bracket.TeamSummary{ID: "ARG", Name: "Argentina"},
		HomeScore: &score2,
		AwayScore: &score1,
		Status:    "finished",
	}

	winner := bracket.CalculateWinner(match)
	if winner == nil {
		t.Fatal("esperava winner, got nil")
	}
	if winner.ID != "BRA" {
		t.Fatalf("esperava BRA, got %s", winner.ID)
	}
}

func TestGetBracket_WinnerNilWhenDraw(t *testing.T) {
	score1 := 1

	match := bracket.BracketMatch{
		HomeTeam:  &bracket.TeamSummary{ID: "BRA", Name: "Brasil"},
		AwayTeam:  &bracket.TeamSummary{ID: "ARG", Name: "Argentina"},
		HomeScore: &score1,
		AwayScore: &score1,
		Status:    "finished",
	}

	winner := bracket.CalculateWinner(match)
	if winner != nil {
		t.Fatal("esperava nil em empate, got winner")
	}
}

func TestGetBracket_WinnerNilWhenScheduled(t *testing.T) {
	match := bracket.BracketMatch{
		HomeTeam: &bracket.TeamSummary{ID: "BRA", Name: "Brasil"},
		AwayTeam: &bracket.TeamSummary{ID: "ARG", Name: "Argentina"},
		Status:   "scheduled",
	}

	winner := bracket.CalculateWinner(match)
	if winner != nil {
		t.Fatal("esperava nil em jogo agendado, got winner")
	}
}

func TestStageOrder(t *testing.T) {
	if bracket.GetStageOrder("final") != 5 {
		t.Fatal("esperava order 5 pra final")
	}
	if bracket.GetStageOrder("semi") != 4 {
		t.Fatal("esperava order 4 pra semi")
	}
}