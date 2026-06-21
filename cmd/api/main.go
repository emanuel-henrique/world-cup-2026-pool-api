package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"

	"bolao-copa/internal/auth"
	"bolao-copa/internal/bracket"
	db "bolao-copa/internal/database"
	"bolao-copa/internal/groups"
	"bolao-copa/internal/matches"
	"bolao-copa/internal/players"
	"bolao-copa/internal/predictions"
	"bolao-copa/internal/ranking"
	"bolao-copa/internal/router"
)

func main() {
    // Carrega o .env — ignora erro em produção (variáveis já vêm do ambiente)
    if err := godotenv.Load(); err != nil {
        log.Println("aviso: arquivo .env não encontrado, usando variáveis do ambiente")
    }

    // Banco
    database, err := db.Connect()
    if err != nil {
        log.Fatalf("falha ao conectar no banco: %v", err)
    }
    defer database.Close()

    // Auth
    authRepo    := auth.NewRepository(database)
    authService := auth.NewService(authRepo)
    authHandler := auth.NewHandler(authService)

    matchRepo    := matches.NewRepository(database)
    matchService := matches.NewService(matchRepo)
    matchHandler := matches.NewHandler(matchService)

    playerRepo    := players.NewRepository(database)
    playerHandler := players.NewHandler(playerRepo)

    groupService := groups.NewService(database)
    groupHandler := groups.NewHandler(groupService)

    bracketService := bracket.NewService(database)
    bracketHandler := bracket.NewHandler(bracketService)

    predRepo    := predictions.NewRepository(database)
    predService := predictions.NewService(predRepo)
    predHandler := predictions.NewHandler(predService)

    rankingService := ranking.NewService(database)
    rankingHandler := ranking.NewHandler(rankingService)

    r := router.New(authHandler, matchHandler, playerHandler, groupHandler, bracketHandler, predHandler, rankingHandler)

    r.Run(":" + os.Getenv("PORT"))
}