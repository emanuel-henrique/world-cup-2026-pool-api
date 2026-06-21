package pagination

import (
	"math"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Meta struct {
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	TotalItems int `json:"totalItems"`
	TotalPages int `json:"totalPages"`
}

func GetParams(c *gin.Context, defaultLimit int) (int, int) {
	pageStr := c.Query("page")
	limitStr := c.Query("limit")

	page := 1
	if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
		page = p
	}

	limit := defaultLimit
	if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
		limit = l
	}

	return page, limit
}

func NewMeta(page, limit, totalItems int) Meta {
	totalPages := int(math.Ceil(float64(totalItems) / float64(limit)))
	if totalPages == 0 {
		totalPages = 1 // At least 1 page even if empty
	}
	return Meta{
		Page:       page,
		Limit:      limit,
		TotalItems: totalItems,
		TotalPages: totalPages,
	}
}

func GetOffset(page, limit int) int {
	return (page - 1) * limit
}
