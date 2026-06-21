package router

import (
	"bolao-copa/internal/auth"
	"bolao-copa/internal/matches"
	"bolao-copa/internal/players"

	// "bolao-copa/internal/bracket"
	// "bolao-copa/internal/groups"
	// "bolao-copa/internal/players"
	// "bolao-copa/internal/predictions"
	// "bolao-copa/internal/ranking"

	"github.com/gin-gonic/gin"
)

func New(
    authHandler     *auth.Handler,
    matchHandler    *matches.Handler,
    // groupHandler    *groups.Handler,
    // bracketHandler  *bracket.Handler,
    playerHandler   *players.Handler,
    // predHandler     *predictions.Handler,
    // rankingHandler  *ranking.Handler,
) *gin.Engine {
    r := gin.Default()

    // Auth
    a := r.Group("/auth")
    {
        a.POST("/register", authHandler.Register)
        a.POST("/login",    authHandler.Login)
    }

    // // Públicas
    r.GET("/matches",        matchHandler.List)
    r.GET("/matches/:id",    matchHandler.GetByID)
    // r.GET("/groups",         groupHandler.List)
    // r.GET("/groups/:name",   groupHandler.GetByName)
    // r.GET("/bracket",        bracketHandler.Get)
    r.GET("/players",        playerHandler.List)
    // r.GET("/ranking",        rankingHandler.List)

    // // Protegidas — JWT obrigatório
    // protected := r.Group("/")
    // protected.Use(auth.Middleware())
    // {
    //     protected.GET("/predictions",           predHandler.List)
    //     protected.POST("/predictions",          predHandler.Upsert)
    //     protected.GET("/predictions/special",   predHandler.GetSpecial)
    //     protected.POST("/predictions/special",  predHandler.UpsertSpecial)
    // }

    return r
}