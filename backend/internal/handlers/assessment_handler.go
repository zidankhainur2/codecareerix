package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

func (h *AssessmentHandler) SubmitAssessment(c *gin.Context) {
	var input models.SubmitAssessmentInput

	// 1. Bind & Validasi Input JSON
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 2. Ambil userID dari context yang di-set oleh middleware
	userIDString, _ := c.Get("userID")
	userID, err := uuid.Parse(userIDString.(string))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "ID pengguna tidak valid"})
		return
	}

	// 3. Panggil repository untuk menyimpan jawaban
	err = h.repo.SubmitAnswers(userID, input.Answers)
	if err != nil {
		// Log error di server untuk debugging
		log.Printf("Gagal menyimpan jawaban asesmen untuk user %s: %v", userID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan jawaban Anda"})
		return
	}

	// 4. Kirim respons sukses
	c.JSON(http.StatusOK, gin.H{"message": "Asesmen berhasil diselesaikan!"})
}