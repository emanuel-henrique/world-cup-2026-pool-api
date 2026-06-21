// internal/groups/handler_test.go
package groups_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"bolao-copa/internal/groups"

	"github.com/gin-gonic/gin"
)

func setupRouter(handler *groups.Handler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/groups",       handler.List)
	r.GET("/groups/:name", handler.GetByName)
	return r
}

func TestHandlerList_Success(t *testing.T) {
	handler := groups.NewHandler(&mockService{})
	r       := setupRouter(handler)

	req  := httptest.NewRequest(http.MethodGet, "/groups", nil)
	resp := httptest.NewRecorder()

	r.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("esperava 200, got %d", resp.Code)
	}

	var result groups.ListGroupsResponse
	json.NewDecoder(resp.Body).Decode(&result)

	if len(result.Groups) != 1 {
		t.Fatalf("esperava 1 grupo, got %d", len(result.Groups))
	}
	if result.Groups[0].Name != "C" {
		t.Fatalf("esperava grupo C, got %s", result.Groups[0].Name)
	}
}

func TestHandlerList_StandingsCount(t *testing.T) {
	handler := groups.NewHandler(&mockService{})
	r       := setupRouter(handler)

	req  := httptest.NewRequest(http.MethodGet, "/groups", nil)
	resp := httptest.NewRecorder()

	r.ServeHTTP(resp, req)

	var result groups.ListGroupsResponse
	json.NewDecoder(resp.Body).Decode(&result)

	if len(result.Groups[0].Standings) != 4 {
		t.Fatalf("esperava 4 times no grupo C, got %d", len(result.Groups[0].Standings))
	}
}

func TestHandlerGetByName_Found(t *testing.T) {
	handler := groups.NewHandler(&mockService{})
	r       := setupRouter(handler)

	req  := httptest.NewRequest(http.MethodGet, "/groups/C", nil)
	resp := httptest.NewRecorder()

	r.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("esperava 200, got %d", resp.Code)
	}

	var result groups.GroupDetailResponse
	json.NewDecoder(resp.Body).Decode(&result)

	if result.Group.Name != "C" {
		t.Fatalf("esperava grupo C, got %s", result.Group.Name)
	}
	if len(result.Matches) != 1 {
		t.Fatalf("esperava 1 jogo, got %d", len(result.Matches))
	}
}

func TestHandlerGetByName_NotFound(t *testing.T) {
	handler := groups.NewHandler(&mockService{})
	r       := setupRouter(handler)

	req  := httptest.NewRequest(http.MethodGet, "/groups/Z", nil)
	resp := httptest.NewRecorder()

	r.ServeHTTP(resp, req)

	if resp.Code != http.StatusNotFound {
		t.Fatalf("esperava 404, got %d", resp.Code)
	}
}