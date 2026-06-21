// internal/scoring/scoring_test.go
package scoring_test

import (
	"testing"

	"bolao-copa/internal/scoring"
)

func TestCalculatePoints_ExactScore(t *testing.T) {
	pts := scoring.CalculatePoints(
		scoring.MatchResult{HomeScore: 2, AwayScore: 1},
		scoring.MatchResult{HomeScore: 2, AwayScore: 1},
	)
	if pts != 10 {
		t.Fatalf("esperava 10, got %d", pts)
	}
}

func TestCalculatePoints_CorrectWinnerAndDiff(t *testing.T) {
	pts := scoring.CalculatePoints(
		scoring.MatchResult{HomeScore: 3, AwayScore: 1},
		scoring.MatchResult{HomeScore: 2, AwayScore: 0},
	)
	if pts != 7 {
		t.Fatalf("esperava 7, got %d", pts)
	}
}

func TestCalculatePoints_CorrectWinner(t *testing.T) {
	pts := scoring.CalculatePoints(
		scoring.MatchResult{HomeScore: 1, AwayScore: 0},
		scoring.MatchResult{HomeScore: 3, AwayScore: 1},
	)
	if pts != 5 {
		t.Fatalf("esperava 5, got %d", pts)
	}
}

func TestCalculatePoints_CorrectDraw(t *testing.T) {
	pts := scoring.CalculatePoints(
		scoring.MatchResult{HomeScore: 1, AwayScore: 1},
		scoring.MatchResult{HomeScore: 0, AwayScore: 0},
	)
	if pts != 5 {
		t.Fatalf("esperava 5, got %d", pts)
	}
}

func TestCalculatePoints_WrongEverything(t *testing.T) {
	pts := scoring.CalculatePoints(
		scoring.MatchResult{HomeScore: 2, AwayScore: 0},
		scoring.MatchResult{HomeScore: 0, AwayScore: 1},
	)
	if pts != 0 {
		t.Fatalf("esperava 0, got %d", pts)
	}
}

func TestCalculatePoints_PredictedDrawButWinner(t *testing.T) {
	pts := scoring.CalculatePoints(
		scoring.MatchResult{HomeScore: 1, AwayScore: 1},
		scoring.MatchResult{HomeScore: 2, AwayScore: 0},
	)
	if pts != 0 {
		t.Fatalf("esperava 0, got %d", pts)
	}
}

func TestCalculatePoints_CorrectAwayWinnerSameMargin(t *testing.T) {
	// 0x2 vs 1x3 — mesmo vencedor (away) e mesma margem (2) → 7 pts
	pts := scoring.CalculatePoints(
		scoring.MatchResult{HomeScore: 0, AwayScore: 2},
		scoring.MatchResult{HomeScore: 1, AwayScore: 3},
	)
	if pts != 7 {
		t.Fatalf("esperava 7, got %d", pts)
	}
}

func TestCalculatePoints_CorrectAwayWinnerDiffMargin(t *testing.T) {
	// 0x1 vs 1x3 — mesmo vencedor (away) mas margem diferente → 5 pts
	pts := scoring.CalculatePoints(
		scoring.MatchResult{HomeScore: 0, AwayScore: 1},
		scoring.MatchResult{HomeScore: 1, AwayScore: 3},
	)
	if pts != 5 {
		t.Fatalf("esperava 5, got %d", pts)
	}
}