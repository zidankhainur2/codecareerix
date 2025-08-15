package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/zidankhainur2/codecareerix/backend/internal/repositories"
)

type ProgressHandler struct {
	repo *repositories.ProgressRepository
}

func NewProgressHandler(repo *repositories.ProgressRepository) *ProgressHandler {
	return &ProgressHandler{repo: repo}
}

// MarkAsComplete menangani permintaan untuk menandai sebuah resource sebagai selesai.
func (h *ProgressHandler) MarkAsComplete(c *gin.Context) {
	// 1. Ambil resource_id dari URL parameter
	resourceIDStr := c.Param("resource_id")
	resourceID, err := strconv.Atoi(resourceIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Resource ID tidak valid"})
		return
	}

	// 2. Ambil userID dari context
	userIDString, _ := c.Get("userID")
	userID, err := uuid.Parse(userIDString.(string))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "ID pengguna tidak valid"})
		return
	}

	// 3. Panggil repository untuk menyimpan progres
	err = h.repo.MarkResourceAsComplete(userID, resourceID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan progres belajar"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Materi berhasil ditandai selesai"})
}