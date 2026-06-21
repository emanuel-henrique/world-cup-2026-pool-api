// internal/brackets/handler_test.go
package bracket_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"bolao-copa/internal/bracket"

	"github.com/gin-gonic/gin"
)

// Mock do Service
type mockService struct{}

func (m *mockService) GetBracket(ctx context.Context) (bracket.BracketResponse, error) {
	score1 := 2
	score2 := 1

	return bracket.BracketResponse{
		Rounds: []bracket.BracketRound{
			{
				Name:  "Quartas de Final",
				Order: 3,
				Matches: []bracket.BracketMatch{
					{
						ID:        "match-1",
						HomeTeam:  &bracket.TeamSummary{ID: "BRA", Name: "Brasil", Flag: "https://flagcdn.com/w40/br.png"},
						AwayTeam:  &bracket.TeamSummary{ID: "ARG", Name: "Argentina", Flag: "https://flagcdn.com/w40/ar.png"},
						HomeScore: &score1,
						AwayScore: &score2,
						Status:    "finished",
						KickoffAt: "2026-07-04T18:00:00Z",
						Winner:    &bracket.TeamSummary{ID: "BRA", Name: "Brasil", Flag: "https://flagcdn.com/w40/br.png"},
					},
				},
			},
			{
				Name:  "Semifinais",
				Order: 4,
				Matches: []bracket.BracketMatch{
					{
						ID:        "match-2",
						HomeTeam:  &bracket.TeamSummary{ID: "BRA", Name: "Brasil", Flag: "https://flagcdn.com/w40/br.png"},
						AwayTeam:  &bracket.TeamSummary{ID: "FRA", Name: "França", Flag: "https://flagcdn.com/w40/fr.png"},
						Status:    "scheduled",
						KickoffAt: "2026-07-08T18:00:00Z",
						Winner:    nil,
					},
				},
			},
		},
	}, nil
}

func setupRouter(handler *bracket.Handler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/bracket", handler.Get)
	return r
}

// Testes

func TestHandlerGet_Success(t *testing.T) {
	handler := bracket.NewHandler(&mockService{})
	r       := setupRouter(handler)

	req  := httptest.NewRequest(http.MethodGet, "/bracket", nil)
	resp := httptest.NewRecorder()

	r.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("esperava 200, got %d", resp.Code)
	}
}

func TestHandlerGet_RoundsCount(t *testing.T) {
	handler := bracket.NewHandler(&mockService{})
	r       := setupRouter(handler)

	req  := httptest.NewRequest(http.MethodGet, "/bracket", nil)
	resp := httptest.NewRecorder()

	r.ServeHTTP(resp, req)

	var result bracket.BracketResponse
	json.NewDecoder(resp.Body).Decode(&result)

	if len(result.Rounds) != 2 {
		t.Fatalf("esperava 2 rounds, got %d", len(result.Rounds))
	}
}

func TestHandlerGet_RoundOrder(t *testing.T) {
	handler := bracket.NewHandler(&mockService{})
	r       := setupRouter(handler)

	req  := httptest.NewRequest(http.MethodGet, "/bracket", nil)
	resp := httptest.NewRecorder()

	r.ServeHTTP(resp, req)

	var result bracket.BracketResponse
	json.NewDecoder(resp.Body).Decode(&result)

	if result.Rounds[0].Order != 3 {
		t.Fatalf("esperava order 3 no primeiro round, got %d", result.Rounds[0].Order)
	}
	if result.Rounds[1].Order != 4 {
		t.Fatalf("esperava order 4 no segundo round, got %d", result.Rounds[1].Order)
	}
}

func TestHandlerGet_WinnerPopulated(t *testing.T) {
	handler := bracket.NewHandler(&mockService{})
	r       := setupRouter(handler)

	req  := httptest.NewRequest(http.MethodGet, "/bracket", nil)
	resp := httptest.NewRecorder()

	r.ServeHTTP(resp, req)

	var result bracket.BracketResponse
	json.NewDecoder(resp.Body).Decode(&result)

	firstMatch := result.Rounds[0].Matches[0]
	if firstMatch.Winner == nil {
		t.Fatal("esperava winner preenchido no jogo finalizado")
	}
	if firstMatch.Winner.ID != "BRA" {
		t.Fatalf("esperava BRA como winner, got %s", firstMatch.Winner.ID)
	}
}

func TestHandlerGet_WinnerNilWhenScheduled(t *testing.T) {
	handler := bracket.NewHandler(&mockService{})
	r       := setupRouter(handler)

	req  := httptest.NewRequest(http.MethodGet, "/bracket", nil)
	resp := httptest.NewRecorder()

	r.ServeHTTP(resp, req)

	var result bracket.BracketResponse
	json.NewDecoder(resp.Body).Decode(&result)

	secondMatch := result.Rounds[1].Matches[0]
	if secondMatch.Winner != nil {
		t.Fatal("esperava winner nil em jogo agendado")
	}
}