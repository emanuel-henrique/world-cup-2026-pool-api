// internal/bracket/handler.go
package bracket

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

func (h *Handler) Get(c *gin.Context) {
	resp, err := h.service.GetBracket(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "erro ao buscar chaveamento"})
		return
	}

	c.JSON(http.StatusOK, resp)
}