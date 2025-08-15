package services

import (
	"github.com/google/uuid"
	"github.com/zidankhainur2/codecareerix/backend/internal/models"
	"github.com/zidankhainur2/codecareerix/backend/internal/repositories"
)

type AssessmentService struct {
	repo *repositories.AssessmentRepository
}

func NewAssessmentService(repo *repositories.AssessmentRepository) *AssessmentService {
	return &AssessmentService{repo: repo}
}

// ProcessAssessmentSubmission adalah alur lengkap dari submit hingga mendapat rekomendasi.
func (s *AssessmentService) ProcessAssessmentSubmission(userID uuid.UUID, answers []models.UserAnswer) ([]models.CareerRecommendation, error) {
	// Langkah 1: Simpan jawaban dan dapatkan ID asesmen
	assessmentID, err := s.repo.SubmitAnswers(userID, answers)
	if err != nil {
		return nil, err
	}

	// Langkah 2: Hitung skor menggunakan ID asesmen
	recommendations, err := s.repo.CalculateScores(assessmentID)
	if err != nil {
		return nil, err
	}

	// Ambil hanya 3 teratas sesuai brief proyek
	if len(recommendations) > 3 {
		return recommendations[:3], nil
	}

	return recommendations, nil
}