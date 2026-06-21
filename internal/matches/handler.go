// internal/matches/handler.go
package matches

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
    service Service
}

func NewHandler(service Service) *Handler {
    return &Handler{service: service}
}

func (h *Handler) List(c *gin.Context) {
    filters := MatchFilters{
        Stage:  c.Query("stage"),
        Status: c.Query("status"),
        Group:  c.Query("group"),
    }

    resp, err := h.service.ListMatches(c.Request.Context(), filters)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "erro ao buscar jogos"})
        return
    }

    c.JSON(http.StatusOK, resp)
}

func (h *Handler) GetByID(c *gin.Context) {
    id := c.Param("id")

    match, err := h.service.GetMatch(c.Request.Context(), id)
    if err == ErrMatchNotFound {
        c.JSON(http.StatusNotFound, gin.H{"error": "jogo não encontrado"})
        return
    }
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "erro ao buscar jogo"})
        return
    }

    c.JSON(http.StatusOK, match)
}