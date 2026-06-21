package auth

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

func (h *Handler) Register(c *gin.Context) {
    var req RegisterRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    token, err := h.service.Register(c.Request.Context(), req)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "erro ao criar conta"})
        return
    }

    c.JSON(http.StatusCreated, token)
}

func (h *Handler) Login(c *gin.Context) {
    var req LoginRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    token, err := h.service.Login(c.Request.Context(), req)
    if err == ErrInvalidCredentials {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "email ou senha inválidos"})
        return
    }
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "erro ao fazer login"})
        return
    }

    c.JSON(http.StatusOK, token)
}