// cmd/api/main.go
package main

import (
	"log"
	"os"

	"bolao-copa/internal/auth"
	db "bolao-copa/internal/database"
	"bolao-copa/internal/router"
)

func main() {
    // 1. Banco
    database, err := db.Connect()
    if err != nil {
        log.Fatalf("falha ao conectar no banco: %v", err)
    }
    defer database.Close()

    // 2. Auth — monta as dependências em cadeia
    authRepo    := auth.NewRepository(database)
    authService := auth.NewService(authRepo)
    authHandler := auth.NewHandler(authService)

    // 3. Router
    r := router.New(authHandler)
    r.Run(":" + os.Getenv("PORT"))
}