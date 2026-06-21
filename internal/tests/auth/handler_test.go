package auth_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"bolao-copa/internal/auth"

	"github.com/gin-gonic/gin"
)

// Mock do Service
type mockService struct{}

func TestMain(m *testing.M) {
    os.Setenv("JWT_SECRET", "secret_de_teste")
    os.Exit(m.Run())
}

func (m *mockService) Register(ctx context.Context, req auth.RegisterRequest) (auth.TokenResponse, error) {
    if req.Email == "existe@email.com" {
        return auth.TokenResponse{}, fmt.Errorf("email já cadastrado")
    }
    return auth.TokenResponse{Token: "token-fake"}, nil
}

func (m *mockService) Login(ctx context.Context, req auth.LoginRequest) (auth.TokenResponse, error) {
    if req.Password != "senha123" {
        return auth.TokenResponse{}, auth.ErrInvalidCredentials
    }
    return auth.TokenResponse{Token: "token-fake"}, nil
}

func setupRouter(handler *auth.Handler) *gin.Engine {
    gin.SetMode(gin.TestMode)
    r := gin.New()
    r.POST("/auth/register", handler.Register)
    r.POST("/auth/login", handler.Login)
    return r
}

func TestHandlerRegister_Success(t *testing.T) {
    handler := auth.NewHandler(&mockService{})
    r       := setupRouter(handler)

    body, _ := json.Marshal(map[string]string{
        "name":     "Emanuel",
        "email":    "novo@email.com",
        "password": "senha123",
    })

    req  := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewBuffer(body))
    req.Header.Set("Content-Type", "application/json")
    resp := httptest.NewRecorder()

    r.ServeHTTP(resp, req)

    if resp.Code != http.StatusCreated {
        t.Fatalf("esperava 201, got %d", resp.Code)
    }

    var result auth.TokenResponse
    json.NewDecoder(resp.Body).Decode(&result)
    if result.Token == "" {
        t.Fatal("esperava token na resposta")
    }
}

func TestHandlerRegister_InvalidBody(t *testing.T) {
    handler := auth.NewHandler(&mockService{})
    r       := setupRouter(handler)

    // manda sem email
    body, _ := json.Marshal(map[string]string{
        "name":     "Emanuel",
        "password": "senha123",
    })

    req  := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewBuffer(body))
    req.Header.Set("Content-Type", "application/json")
    resp := httptest.NewRecorder()

    r.ServeHTTP(resp, req)

    if resp.Code != http.StatusBadRequest {
        t.Fatalf("esperava 400, got %d", resp.Code)
    }
}

func TestHandlerLogin_WrongPassword(t *testing.T) {
    handler := auth.NewHandler(&mockService{})
    r       := setupRouter(handler)

    body, _ := json.Marshal(map[string]string{
        "email":    "emanuel@email.com",
        "password": "errada",
    })

    req  := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(body))
    req.Header.Set("Content-Type", "application/json")
    resp := httptest.NewRecorder()

    r.ServeHTTP(resp, req)

    if resp.Code != http.StatusUnauthorized {
        t.Fatalf("esperava 401, got %d", resp.Code)
    }
}