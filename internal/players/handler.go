// internal/players/handler.go
package players

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
    repo Repository
}

func NewHandler(repo Repository) *Handler {
    return &Handler{repo: repo}
}

func (h *Handler) List(c *gin.Context) {
    players, err := h.repo.FindAll(c.Request.Context())
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "erro ao buscar jogadores"})
        return
    }

    c.JSON(http.StatusOK, ListPlayersResponse{
        Players: players,
        Total:   len(players),
    })
}