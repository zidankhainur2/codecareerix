package handlers

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/zidankhainur2/codecareerix/backend/internal/repositories"
)

type LearningHandler struct {
	learningRepo *repositories.LearningRepository
	userRepo     *repositories.UserRepository
}

func NewLearningHandler(learningRepo *repositories.LearningRepository, userRepo *repositories.UserRepository) *LearningHandler {
	return &LearningHandler{
		learningRepo: learningRepo,
		userRepo:     userRepo,
	}
}

// GetMyRoadmap mengambil roadmap belajar berdasarkan jalur karier aktif pengguna.
func (h *LearningHandler) GetMyRoadmap(c *gin.Context) {
	// 1. Ambil userID dari context
	userIDString, _ := c.Get("userID")
	userID, err := uuid.Parse(userIDString.(string))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "ID pengguna tidak valid"})
		return
	}

	// 2. Ambil data pengguna untuk mendapatkan active_career_path_id
	user, err := h.userRepo.GetByID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data pengguna"})
		return
	}

	// 3. Cek apakah pengguna sudah memilih jalur karier
	if user.ActiveCareerPathID == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Anda belum memilih jalur karier."})
		return
	}

	// 4. Panggil repository untuk mendapatkan roadmap
	roadmap, err := h.learningRepo.GetRoadmapByCareerPathID(*user.ActiveCareerPathID)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Roadmap untuk jalur karier ini tidak ditemukan."})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil roadmap belajar."})
		return
	}

	c.JSON(http.StatusOK, roadmap)
}