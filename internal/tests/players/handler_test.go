// internal/players/handler_test.go
package players_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"bolao-copa/internal/players"

	"github.com/gin-gonic/gin"
)

// Mock do Repository
type mockRepository struct{}

func (m *mockRepository) FindAll(ctx context.Context) ([]players.PlayerResponse, error) {
	return []players.PlayerResponse{
		{
			ID:   "player-1",
			Name: "Vinicius Jr.",
			Team: players.TeamSummary{ID: "BRA", Name: "Brasil", Flag: "https://flagcdn.com/w40/br.png"},
		},
		{
			ID:   "player-2",
			Name: "Kylian Mbappé",
			Team: players.TeamSummary{ID: "FRA", Name: "França", Flag: "https://flagcdn.com/w40/fr.png"},
		},
	}, nil
}

func setupRouter(handler *players.Handler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/players", handler.List)
	return r
}

// Testes

func TestHandlerList_Success(t *testing.T) {
	handler := players.NewHandler(&mockRepository{})
	r       := setupRouter(handler)

	req  := httptest.NewRequest(http.MethodGet, "/players", nil)
	resp := httptest.NewRecorder()

	r.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("esperava 200, got %d", resp.Code)
	}

	var result players.ListPlayersResponse
	json.NewDecoder(resp.Body).Decode(&result)

	if result.Total != 2 {
		t.Fatalf("esperava 2 jogadores, got %d", result.Total)
	}
}

func TestHandlerList_CheckFields(t *testing.T) {
	handler := players.NewHandler(&mockRepository{})
	r       := setupRouter(handler)

	req  := httptest.NewRequest(http.MethodGet, "/players", nil)
	resp := httptest.NewRecorder()

	r.ServeHTTP(resp, req)

	var result players.ListPlayersResponse
	json.NewDecoder(resp.Body).Decode(&result)

	first := result.Players[0]
	if first.ID != "player-1" {
		t.Fatalf("esperava player-1, got %s", first.ID)
	}
	if first.Team.ID != "BRA" {
		t.Fatalf("esperava BRA, got %s", first.Team.ID)
	}
	if first.Team.Flag == "" {
		t.Fatal("esperava flag preenchida")
	}
}