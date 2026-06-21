// internal/groups/handler.go
package groups

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
    resp, err := h.service.ListGroups(c.Request.Context())
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "erro ao buscar grupos"})
        return
    }

    c.JSON(http.StatusOK, resp)
}

func (h *Handler) GetByName(c *gin.Context) {
    name := c.Param("name")

    resp, err := h.service.GetGroup(c.Request.Context(), name)
    if err == ErrGroupNotFound {
        c.JSON(http.StatusNotFound, gin.H{"error": "grupo não encontrado"})
        return
    }
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "erro ao buscar grupo"})
        return
    }

    c.JSON(http.StatusOK, resp)
}