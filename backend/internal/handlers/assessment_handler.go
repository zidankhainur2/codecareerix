package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zidankhainur2/codecareerix/backend/internal/models"
	"github.com/zidankhainur2/codecareerix/backend/internal/repositories"
)

type AssessmentHandler struct {
	repo *repositories.AssessmentRepository
}

func NewAssessmentHandler(repo *repositories.AssessmentRepository) *AssessmentHandler {
	return &AssessmentHandler{repo: repo}
}

// GetAssessmentQuestions menangani permintaan untuk mendapatkan semua pertanyaan asesmen.
func (h *AssessmentHandler) GetAssessmentQuestions(c *gin.Context) {
	questions, err := h.repo.GetAllQuestionsWithOptions()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil pertanyaan asesmen"})
		return
	}

	// Jika tidak ada pertanyaan sama sekali
	if len(questions) == 0 {
		c.JSON(http.StatusOK, []models.AssessmentQuestion{})
		return
	}

	c.JSON(http.StatusOK, questions)
}