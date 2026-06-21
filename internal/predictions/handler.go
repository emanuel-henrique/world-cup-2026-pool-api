// internal/predictions/handler.go
package predictions

import (
	"net/http"

	"bolao-copa/internal/auth"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) List(c *gin.Context) {
	userID := auth.GetUserID(c)

	resp, err := h.service.ListPredictions(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "erro ao buscar palpites"})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *Handler) Upsert(c *gin.Context) {
	userID := auth.GetUserID(c)

	var req UpsertPredictionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.service.UpsertPrediction(c.Request.Context(), userID, req)
	if err == ErrMatchLocked {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}
	if err == ErrMatchNotFound {
		c.JSON(http.StatusNotFound, gin.H{"error": "jogo não encontrado"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "erro ao salvar palpite"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "palpite salvo com sucesso"})
}

func (h *Handler) GetSpecial(c *gin.Context) {
	userID := auth.GetUserID(c)

	resp, err := h.service.GetSpecial(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "erro ao buscar palpite especial"})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *Handler) UpsertSpecial(c *gin.Context) {
	userID := auth.GetUserID(c)

	var req UpsertSpecialRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.service.UpsertSpecial(c.Request.Context(), userID, req)
	if err == ErrSpecialLocked {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "erro ao salvar palpite especial"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "palpite especial salvo com sucesso"})
}