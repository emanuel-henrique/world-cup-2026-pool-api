// internal/ranking/handler_test.go
package ranking_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"bolao-copa/internal/ranking"

	"github.com/gin-gonic/gin"
)

func setupRouter(handler *ranking.Handler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/ranking", handler.List)
	return r
}

func TestHandlerList_Success(t *testing.T) {
	handler := ranking.NewHandler(&mockService{})
	r       := setupRouter(handler)

	req  := httptest.NewRequest(http.MethodGet, "/ranking", nil)
	resp := httptest.NewRecorder()

	r.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("esperava 200, got %d", resp.Code)
	}

	var result ranking.RankingResponse
	json.NewDecoder(resp.Body).Decode(&result)

	if result.Meta.TotalItems != 3 {
		t.Fatalf("esperava total 3, got %d", result.Meta.TotalItems)
	}
}

func TestHandlerList_DefaultPagination(t *testing.T) {
	handler := ranking.NewHandler(&mockService{})
	r       := setupRouter(handler)

	req  := httptest.NewRequest(http.MethodGet, "/ranking", nil)
	resp := httptest.NewRecorder()

	r.ServeHTTP(resp, req)

	var result ranking.RankingResponse
	json.NewDecoder(resp.Body).Decode(&result)

	if result.Meta.Page != 1 {
		t.Fatalf("esperava page 1, got %d", result.Meta.Page)
	}
	if result.Meta.Limit != 20 {
		t.Fatalf("esperava limit 20, got %d", result.Meta.Limit)
	}
}

func TestHandlerList_CustomPagination(t *testing.T) {
	handler := ranking.NewHandler(&mockService{})
	r       := setupRouter(handler)

	req  := httptest.NewRequest(http.MethodGet, "/ranking?page=2&limit=10", nil)
	resp := httptest.NewRecorder()

	r.ServeHTTP(resp, req)

	var result ranking.RankingResponse
	json.NewDecoder(resp.Body).Decode(&result)

	if result.Meta.Page != 2 {
		t.Fatalf("esperava page 2, got %d", result.Meta.Page)
	}
	if result.Meta.Limit != 10 {
		t.Fatalf("esperava limit 10, got %d", result.Meta.Limit)
	}
}

func TestHandlerList_PositionOrder(t *testing.T) {
	handler := ranking.NewHandler(&mockService{})
	r       := setupRouter(handler)

	req  := httptest.NewRequest(http.MethodGet, "/ranking", nil)
	resp := httptest.NewRecorder()

	r.ServeHTTP(resp, req)

	var result ranking.RankingResponse
	json.NewDecoder(resp.Body).Decode(&result)

	if result.Entries[0].Position != 1 {
		t.Fatalf("esperava posição 1 no primeiro, got %d", result.Entries[0].Position)
	}
	if result.Entries[0].TotalPoints < result.Entries[1].TotalPoints {
		t.Fatal("esperava primeiro lugar com mais pontos que segundo")
	}
}

func TestHandlerList_EntriesCount(t *testing.T) {
	handler := ranking.NewHandler(&mockService{})
	r       := setupRouter(handler)

	req  := httptest.NewRequest(http.MethodGet, "/ranking", nil)
	resp := httptest.NewRecorder()

	r.ServeHTTP(resp, req)

	var result ranking.RankingResponse
	json.NewDecoder(resp.Body).Decode(&result)

	if len(result.Entries) != 3 {
		t.Fatalf("esperava 3 entradas, got %d", len(result.Entries))
	}
}