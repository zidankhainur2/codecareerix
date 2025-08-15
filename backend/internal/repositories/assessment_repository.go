package repositories

import (
	"database/sql"

	"github.com/zidankhainur2/codecareerix/backend/internal/models"
)

type AssessmentRepository struct {
	db *sql.DB
}

func NewAssessmentRepository(db *sql.DB) *AssessmentRepository {
	return &AssessmentRepository{db: db}
}

// GetAllQuestionsWithOptions mengambil semua pertanyaan beserta pilihan jawabannya.
func (r *AssessmentRepository) GetAllQuestionsWithOptions() ([]models.AssessmentQuestion, error) {
	// Query untuk mengambil semua pertanyaan
	queryQuestions := `SELECT id, question_text, question_type FROM assessment_questions ORDER BY id`
	rows, err := r.db.Query(queryQuestions)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Gunakan map untuk menyusun pertanyaan agar mudah dicari
	questionsMap := make(map[int]*models.AssessmentQuestion)
	var questionsList []*models.AssessmentQuestion

	for rows.Next() {
		var q models.AssessmentQuestion
		if err := rows.Scan(&q.ID, &q.QuestionText, &q.QuestionType); err != nil {
			return nil, err
		}
		q.Options = []models.AssessmentOption{} // Inisialisasi slice options
		questionsMap[q.ID] = &q
		questionsList = append(questionsList, &q)
	}

	// Query untuk mengambil semua pilihan jawaban
	queryOptions := `SELECT id, question_id, option_text FROM assessment_options ORDER BY id`
	rowsOptions, err := r.db.Query(queryOptions)
	if err != nil {
		return nil, err
	}
	defer rowsOptions.Close()

	for rowsOptions.Next() {
		var opt models.AssessmentOption
		if err := rowsOptions.Scan(&opt.ID, &opt.QuestionID, &opt.OptionText); err != nil {
			return nil, err
		}

		// Jika pertanyaan untuk opsi ini ada di map, tambahkan opsi ke pertanyaan tersebut
		if question, found := questionsMap[opt.QuestionID]; found {
			question.Options = append(question.Options, opt)
		}
	}

	// Konversi dari slice of pointers ke slice of values
	var result []models.AssessmentQuestion
	for _, q := range questionsList {
		result = append(result, *q)
	}

	return result, nil
}