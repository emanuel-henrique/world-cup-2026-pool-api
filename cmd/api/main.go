// cmd/api/main.go
package main

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"

	"bolao-copa/internal/auth"
	"bolao-copa/internal/bracket"
	"bolao-copa/internal/database"
	"bolao-copa/internal/groups"
	"bolao-copa/internal/matches"
	"bolao-copa/internal/players"
	"bolao-copa/internal/predictions"
	"bolao-copa/internal/ranking"
	"bolao-copa/internal/router"
	"bolao-copa/internal/worker"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("aviso: .env não encontrado, usando variáveis do ambiente")
	}

	db, err := database.Connect()
	if err != nil {
		log.Fatalf("falha ao conectar no banco: %v", err)
	}
	defer db.Close()

	// Auth
	authRepo    := auth.NewRepository(db)
	authService := auth.NewService(authRepo)
	authHandler := auth.NewHandler(authService)

	// Matches
	matchRepo    := matches.NewRepository(db)
	matchService := matches.NewService(matchRepo)
	matchHandler := matches.NewHandler(matchService)

	// Players
	playerRepo    := players.NewRepository(db)
	playerHandler := players.NewHandler(playerRepo)

	// Groups
	groupService := groups.NewService(db)
	groupHandler := groups.NewHandler(groupService)

	// Bracket
	bracketService := bracket.NewService(db)
	bracketHandler := bracket.NewHandler(bracketService)

	// Predictions
	predRepo    := predictions.NewRepository(db)
	predService := predictions.NewService(predRepo)
	predHandler := predictions.NewHandler(predService)

	// Ranking
	rankingService := ranking.NewService(db)
	rankingHandler := ranking.NewHandler(rankingService)

	// Worker — roda em goroutine separada
	w := worker.New(db)
	go w.Start(context.Background())

	// Router
	r := router.New(
		authHandler,
		matchHandler,
		playerHandler,
		groupHandler,
		bracketHandler,
		predHandler,
		rankingHandler,
	)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	log.Printf("servidor rodando na porta %s", port)
	r.Run(":" + port)
}