package handlers

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

// HealthHandler adalah struct untuk handler health check.
type HealthHandler struct {
	db *sql.DB
}

// NewHealthHandler membuat instance baru dari HealthHandler.
func NewHealthHandler(db *sql.DB) *HealthHandler {
	return &HealthHandler{db: db}
}

// Check adalah metode yang menangani request ke /health.
func (h *HealthHandler) Check(c *gin.Context) {
	err := h.db.Ping()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":   "DOWN",
			"database": "Error: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":   "UP",
		"database": "Connected",
	})
}