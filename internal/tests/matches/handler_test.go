// internal/matches/handler_test.go
package matches_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"bolao-copa/internal/matches"

	"github.com/gin-gonic/gin"
)

// Mock do Service
type mockService struct{}

func (m *mockService) ListMatches(ctx context.Context, filters matches.MatchFilters) (matches.ListMatchesResponse, error) {
	score := 1
	group := "C"
	return matches.ListMatchesResponse{
		Total: 1,
		Matches: []matches.MatchResponse{
			{
				ID:        "match-1",
				HomeTeam:  &matches.TeamSummary{ID: "BRA", Name: "Brasil", Flag: "https://flagcdn.com/w40/br.png"},
				AwayTeam:  &matches.TeamSummary{ID: "MAR", Name: "Marrocos", Flag: "https://flagcdn.com/w40/ma.png"},
				HomeScore: &score,
				AwayScore: &score,
				Stage:     string(matches.StageGroup),
				GroupName: &group,
				Status:    string(matches.StatusLive),
			},
		},
	}, nil
}

func (m *mockService) GetMatch(ctx context.Context, id string) (matches.MatchResponse, error) {
	if id != "match-1" {
		return matches.MatchResponse{}, matches.ErrMatchNotFound
	}
	return matches.MatchResponse{ID: "match-1", Status: string(matches.StatusLive)}, nil
}

func setupRouter(handler *matches.Handler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/matches",      handler.List)
	r.GET("/matches/:id",  handler.GetByID)
	return r
}

// Testes

func TestHandlerList_Success(t *testing.T) {
	handler := matches.NewHandler(&mockService{})
	r       := setupRouter(handler)

	req  := httptest.NewRequest(http.MethodGet, "/matches", nil)
	resp := httptest.NewRecorder()

	r.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("esperava 200, got %d", resp.Code)
	}

	var result matches.ListMatchesResponse
	json.NewDecoder(resp.Body).Decode(&result)

	if result.Total != 1 {
		t.Fatalf("esperava 1 jogo, got %d", result.Total)
	}
}

func TestHandlerList_WithStatusFilter(t *testing.T) {
	handler := matches.NewHandler(&mockService{})
	r       := setupRouter(handler)

	req  := httptest.NewRequest(http.MethodGet, "/matches?status=live", nil)
	resp := httptest.NewRecorder()

	r.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("esperava 200, got %d", resp.Code)
	}
}

func TestHandlerGetByID_Found(t *testing.T) {
	handler := matches.NewHandler(&mockService{})
	r       := setupRouter(handler)

	req  := httptest.NewRequest(http.MethodGet, "/matches/match-1", nil)
	resp := httptest.NewRecorder()

	r.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("esperava 200, got %d", resp.Code)
	}

	var result matches.MatchResponse
	json.NewDecoder(resp.Body).Decode(&result)

	if result.ID != "match-1" {
		t.Fatalf("esperava match-1, got %s", result.ID)
	}
}

func TestHandlerGetByID_NotFound(t *testing.T) {
	handler := matches.NewHandler(&mockService{})
	r       := setupRouter(handler)

	req  := httptest.NewRequest(http.MethodGet, "/matches/nao-existe", nil)
	resp := httptest.NewRecorder()

	r.ServeHTTP(resp, req)

	if resp.Code != http.StatusNotFound {
		t.Fatalf("esperava 404, got %d", resp.Code)
	}
}