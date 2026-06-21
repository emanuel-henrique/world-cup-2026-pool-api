// internal/predictions/handler_test.go
package predictions_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"bolao-copa/internal/auth"
	"bolao-copa/internal/pagination"
	"bolao-copa/internal/predictions"

	"github.com/gin-gonic/gin"
)

// Mock do Service
type mockService struct{}

func (m *mockService) ListPredictions(ctx context.Context, userID string, page, limit int) (predictions.ListPredictionsResponse, error) {
	return predictions.ListPredictionsResponse{
		Meta: pagination.Meta{
			Page:       1,
			Limit:      50,
			TotalItems: 1,
			TotalPages: 1,
		},
		Predictions: []predictions.PredictionResponse{{ID: "pred-1", MatchID: "match-1"}},
	}, nil
}

func (m *mockService) UpsertPrediction(ctx context.Context, userID string, req predictions.UpsertPredictionRequest) error {
	if req.MatchID == "locked" {
		return predictions.ErrMatchLocked
	}
	return nil
}

func (m *mockService) GetSpecial(ctx context.Context, userID string) (predictions.SpecialPredictionResponse, error) {
	return predictions.SpecialPredictionResponse{}, nil
}

func (m *mockService) UpsertSpecial(ctx context.Context, userID string, req predictions.UpsertSpecialRequest) error {
	return nil
}

func setupRouter(handler *predictions.Handler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	// Injeta userID no contexto simulando o middleware JWT
	r.Use(func(c *gin.Context) {
		c.Set(auth.UserIDKey, "user-1")
		c.Next()
	})

	r.GET("/predictions",          handler.List)
	r.POST("/predictions",         handler.Upsert)
	r.GET("/predictions/special",  handler.GetSpecial)
	r.POST("/predictions/special", handler.UpsertSpecial)
	return r
}

func TestHandlerList_Success(t *testing.T) {
	handler := predictions.NewHandler(&mockService{})
	r       := setupRouter(handler)

	req  := httptest.NewRequest(http.MethodGet, "/predictions", nil)
	resp := httptest.NewRecorder()

	r.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("esperava 200, got %d", resp.Code)
	}

	var result predictions.ListPredictionsResponse
	json.NewDecoder(resp.Body).Decode(&result)

	if result.Meta.TotalItems != 1 {
		t.Fatalf("esperava 1 palpite, got %d", result.Meta.TotalItems)
	}
}

func TestHandlerUpsert_Success(t *testing.T) {
	handler := predictions.NewHandler(&mockService{})
	r       := setupRouter(handler)

	body, _ := json.Marshal(predictions.UpsertPredictionRequest{
		MatchID:   "match-1",
		HomeScore: 2,
		AwayScore: 1,
	})

	req  := httptest.NewRequest(http.MethodPost, "/predictions", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	r.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("esperava 200, got %d", resp.Code)
	}
}

func TestHandlerUpsert_MatchLocked(t *testing.T) {
	handler := predictions.NewHandler(&mockService{})
	r       := setupRouter(handler)

	body, _ := json.Marshal(predictions.UpsertPredictionRequest{
		MatchID:   "locked",
		HomeScore: 1,
		AwayScore: 0,
	})

	req  := httptest.NewRequest(http.MethodPost, "/predictions", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	r.ServeHTTP(resp, req)

	if resp.Code != http.StatusUnprocessableEntity {
		t.Fatalf("esperava 422, got %d", resp.Code)
	}
}

func TestHandlerGetSpecial_Success(t *testing.T) {
	handler := predictions.NewHandler(&mockService{})
	r       := setupRouter(handler)

	req  := httptest.NewRequest(http.MethodGet, "/predictions/special", nil)
	resp := httptest.NewRecorder()

	r.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("esperava 200, got %d", resp.Code)
	}
}

func TestHandlerUpsertSpecial_Success(t *testing.T) {
	handler := predictions.NewHandler(&mockService{})
	r       := setupRouter(handler)

	champion := "BRA"
	body, _  := json.Marshal(predictions.UpsertSpecialRequest{ChampionID: &champion})

	req  := httptest.NewRequest(http.MethodPost, "/predictions/special", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	r.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("esperava 200, got %d", resp.Code)
	}
}