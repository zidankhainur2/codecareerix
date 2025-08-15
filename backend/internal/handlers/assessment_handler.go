package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/zidankhainur2/codecareerix/backend/internal/models"
	"github.com/zidankhainur2/codecareerix/backend/internal/repositories"
	"github.com/zidankhainur2/codecareerix/backend/internal/services"
)

type AssessmentHandler struct {
	repo    *repositories.AssessmentRepository
	service *services.AssessmentService
}

func NewAssessmentHandler(repo *repositories.AssessmentRepository, service *services.AssessmentService) *AssessmentHandler {
	return &AssessmentHandler{
		repo:    repo,
		service: service,
	}
}

// GetAssessmentQuestions menangani permintaan untuk mendapatkan semua pertanyaan asesmen.
func (h *AssessmentHandler) GetAssessmentQuestions(c *gin.Context) {
	questions, err := h.repo.GetAllQuestionsWithOptions()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil pertanyaan asesmen"})
		return
	}

	if len(questions) == 0 {
		c.JSON(http.StatusOK, []models.AssessmentQuestion{})
		return
	}

	c.JSON(http.StatusOK, questions)
}

// SubmitAssessment menangani pengiriman jawaban dan mengembalikan rekomendasi.
func (h *AssessmentHandler) SubmitAssessment(c *gin.Context) {
	var input models.SubmitAssessmentInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userIDString, _ := c.Get("userID")
	userID, err := uuid.Parse(userIDString.(string))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "ID pengguna tidak valid"})
		return
	}

	// Panggil service untuk memproses jawaban
	recommendations, err := h.service.ProcessAssessmentSubmission(userID, input.Answers)
	if err != nil {
		log.Printf("Gagal memproses asesmen untuk user %s: %v", userID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memproses hasil asesmen Anda"})
		return
	}

	// Kirim hasil rekomendasi sebagai respons
	c.JSON(http.StatusOK, recommendations)
}