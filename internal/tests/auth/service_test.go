package auth_test

import (
	"context"
	"fmt"
	"testing"

	"bolao-copa/internal/auth"

	"github.com/golang-jwt/jwt/v5"
)

// Mock do Repository
type mockRepository struct {
    users map[string]auth.User
}

func newMockRepository() *mockRepository {
    return &mockRepository{users: make(map[string]auth.User)}
}

func (m *mockRepository) CreateUser(ctx context.Context, user auth.User) error {
    if _, exists := m.users[user.Email]; exists {
        return fmt.Errorf("email já cadastrado")
    }
    m.users[user.Email] = user
    return nil
}

func (m *mockRepository) FindUserByEmail(ctx context.Context, email string) (auth.User, error) {
    user, exists := m.users[email]
    if !exists {
        return auth.User{}, auth.ErrUserNotFound
    }
    return user, nil
}

// Testes do Register
func TestRegister_Success(t *testing.T) {
    repo    := newMockRepository()
    service := auth.NewService(repo)

    req := auth.RegisterRequest{
        Name:     "Emanuel",
        Email:    "emanuel@email.com",
        Password: "senha123",
    }

    resp, err := service.Register(context.Background(), req)

    if err != nil {
        t.Fatalf("esperava sucesso, got erro: %v", err)
    }
    if resp.Token == "" {
        t.Fatal("esperava token, got string vazia")
    }
}

func TestRegister_DuplicateEmail(t *testing.T) {
    repo    := newMockRepository()
    service := auth.NewService(repo)

    req := auth.RegisterRequest{
        Name:     "Emanuel",
        Email:    "emanuel@email.com",
        Password: "senha123",
    }

    // primeiro cadastro — deve funcionar
    _, err := service.Register(context.Background(), req)
    if err != nil {
        t.Fatalf("primeiro cadastro falhou: %v", err)
    }

    // segundo cadastro com mesmo email — deve falhar
    _, err = service.Register(context.Background(), req)
    if err == nil {
        t.Fatal("esperava erro de email duplicado, got nil")
    }
}

// Testes do Login
func TestLogin_Success(t *testing.T) {
    repo    := newMockRepository()
    service := auth.NewService(repo)

    // cadastra primeiro
    _, err := service.Register(context.Background(), auth.RegisterRequest{
        Name:     "Emanuel",
        Email:    "emanuel@email.com",
        Password: "senha123",
    })
    if err != nil {
        t.Fatalf("cadastro falhou: %v", err)
    }

    // tenta login
    resp, err := service.Login(context.Background(), auth.LoginRequest{
        Email:    "emanuel@email.com",
        Password: "senha123",
    })

    if err != nil {
        t.Fatalf("esperava sucesso, got erro: %v", err)
    }
    if resp.Token == "" {
        t.Fatal("esperava token, got string vazia")
    }
}

func TestLogin_WrongPassword(t *testing.T) {
    repo    := newMockRepository()
    service := auth.NewService(repo)

    _, _ = service.Register(context.Background(), auth.RegisterRequest{
        Name:     "Emanuel",
        Email:    "emanuel@email.com",
        Password: "senha123",
    })

    _, err := service.Login(context.Background(), auth.LoginRequest{
        Email:    "emanuel@email.com",
        Password: "senhaerrada",
    })

    if err != auth.ErrInvalidCredentials {
        t.Fatalf("esperava ErrInvalidCredentials, got: %v", err)
    }
}

func TestLogin_UserNotFound(t *testing.T) {
    repo    := newMockRepository()
    service := auth.NewService(repo)

    _, err := service.Login(context.Background(), auth.LoginRequest{
        Email:    "naoexiste@email.com",
        Password: "senha123",
    })

    if err != auth.ErrInvalidCredentials {
        t.Fatalf("esperava ErrInvalidCredentials, got: %v", err)
    }
}

// Testa que o JWT contém os claims esperados pelo frontend (user_id, name, email)
func TestRegister_TokenContainsUserClaims(t *testing.T) {
    repo    := newMockRepository()
    service := auth.NewService(repo)

    req := auth.RegisterRequest{
        Name:     "Emanuel",
        Email:    "emanuel@email.com",
        Password: "senha123",
    }

    resp, err := service.Register(context.Background(), req)
    if err != nil {
        t.Fatalf("register falhou: %v", err)
    }

    claims := parseTokenClaims(t, resp.Token)

    if claims["user_id"] != claims["sub"] {
        t.Fatalf("user_id (%v) deveria ser igual a sub (%v)", claims["user_id"], claims["sub"])
    }
    if claims["name"] != "Emanuel" {
        t.Fatalf("esperava name=Emanuel, got %v", claims["name"])
    }
    if claims["email"] != "emanuel@email.com" {
        t.Fatalf("esperava email=emanuel@email.com, got %v", claims["email"])
    }
}

func TestLogin_TokenContainsUserClaims(t *testing.T) {
    repo    := newMockRepository()
    service := auth.NewService(repo)

    _, err := service.Register(context.Background(), auth.RegisterRequest{
        Name:     "Maria",
        Email:    "maria@email.com",
        Password: "senha123",
    })
    if err != nil {
        t.Fatalf("register falhou: %v", err)
    }

    resp, err := service.Login(context.Background(), auth.LoginRequest{
        Email:    "maria@email.com",
        Password: "senha123",
    })
    if err != nil {
        t.Fatalf("login falhou: %v", err)
    }

    claims := parseTokenClaims(t, resp.Token)

    if claims["user_id"] == nil || claims["user_id"] == "" {
        t.Fatal("esperava user_id no token")
    }
    if claims["name"] != "Maria" {
        t.Fatalf("esperava name=Maria, got %v", claims["name"])
    }
    if claims["email"] != "maria@email.com" {
        t.Fatalf("esperava email=maria@email.com, got %v", claims["email"])
    }
}

func parseTokenClaims(t *testing.T, tokenStr string) jwt.MapClaims {
    t.Helper()
    token, _, err := jwt.NewParser().ParseUnverified(tokenStr, jwt.MapClaims{})
    if err != nil {
        t.Fatalf("erro ao parsear token: %v", err)
    }
    claims, ok := token.Claims.(jwt.MapClaims)
    if !ok {
        t.Fatal("claims não é MapClaims")
    }
    return claims
}