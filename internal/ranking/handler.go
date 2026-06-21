// internal/ranking/handler.go
package ranking

import (
	"bolao-copa/internal/pagination"
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
	page, limit := pagination.GetParams(c, 20)

	resp, err := h.service.GetRanking(c.Request.Context(), page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "erro ao buscar ranking"})
		return
	}

	c.JSON(http.StatusOK, resp)
}