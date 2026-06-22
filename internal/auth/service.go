// internal/auth/service.go
package auth

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Service interface {
    Register(ctx context.Context, req RegisterRequest) (TokenResponse, error)
    Login(ctx context.Context, req LoginRequest) (TokenResponse, error)
}

type service struct {
    repo Repository
}

func NewService(repo Repository) Service {
    return &service{repo: repo}
}

func (s *service) Register(ctx context.Context, req RegisterRequest) (TokenResponse, error) {
    hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
    if err != nil {
        return TokenResponse{}, fmt.Errorf("erro ao gerar hash da senha: %w", err)
    }

    user := User{
        ID:        uuid.NewString(),
        Name:      req.Name,
        Email:     req.Email,
        Password:  string(hash),
        CreatedAt: time.Now(),
    }

    if err := s.repo.CreateUser(ctx, user); err != nil {
        return TokenResponse{}, fmt.Errorf("erro ao criar usuário: %w", err)
    }

    token, err := generateToken(user)
    if err != nil {
        return TokenResponse{}, fmt.Errorf("erro ao gerar token: %w", err)
    }

    return TokenResponse{Token: token}, nil
}

func (s *service) Login(ctx context.Context, req LoginRequest) (TokenResponse, error) {
    user, err := s.repo.FindUserByEmail(ctx, req.Email)
    if err == ErrUserNotFound {
        return TokenResponse{}, ErrInvalidCredentials
    }
    if err != nil {
        return TokenResponse{}, fmt.Errorf("erro ao buscar usuário: %w", err)
    }

    if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
        return TokenResponse{}, ErrInvalidCredentials
    }

    token, err := generateToken(user)
    if err != nil {
        return TokenResponse{}, fmt.Errorf("erro ao gerar token: %w", err)
    }

   return TokenResponse{Token: token}, nil
}

func generateToken(user User) (string, error) {
    secret := os.Getenv("JWT_SECRET")
    if secret == "" {
        return "", fmt.Errorf("JWT_SECRET não definido")
    }

    claims := jwt.MapClaims{
        "sub":     user.ID,
        "user_id": user.ID,
        "name":    user.Name,
        "email":   user.Email,
        "exp":     time.Now().Add(7 * 24 * time.Hour).Unix(),
        "iat":     time.Now().Unix(),
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(secret))
}

// Erros de domínio
var ErrInvalidCredentials = fmt.Errorf("email ou senha inválidos")