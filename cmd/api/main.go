package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"

	"bolao-copa/internal/auth"
	db "bolao-copa/internal/database"
	"bolao-copa/internal/matches"
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

    r := router.New(authHandler, matchHandler)

    r.Run(":" + os.Getenv("PORT"))
}