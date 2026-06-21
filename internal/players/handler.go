// internal/players/handler.go
package players

import (
	"bolao-copa/internal/pagination"
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
    page, limit := pagination.GetParams(c, 50)
    offset := pagination.GetOffset(page, limit)

    players, total, err := h.repo.FindAll(c.Request.Context(), limit, offset)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "erro ao buscar jogadores"})
        return
    }

    c.JSON(http.StatusOK, ListPlayersResponse{
        Players: players,
        Meta:    pagination.NewMeta(page, limit, total),
    })
}