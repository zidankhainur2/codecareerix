package repositories

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/google/uuid"
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
	queryQuestions := `SELECT id, question_text, question_type FROM assessment_questions ORDER BY id`
	rows, err := r.db.Query(queryQuestions)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	questionsMap := make(map[int]*models.AssessmentQuestion)
	var questionsList []*models.AssessmentQuestion

	for rows.Next() {
		var q models.AssessmentQuestion
		if err := rows.Scan(&q.ID, &q.QuestionText, &q.QuestionType); err != nil {
			return nil, err
		}
		q.Options = []models.AssessmentOption{}
		questionsMap[q.ID] = &q
		questionsList = append(questionsList, &q)
	}

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

		if question, found := questionsMap[opt.QuestionID]; found {
			question.Options = append(question.Options, opt)
		}
	}

	var result []models.AssessmentQuestion
	for _, q := range questionsList {
		result = append(result, *q)
	}

	return result, nil
}

// SubmitAnswers menyimpan jawaban asesmen dan mengembalikan ID asesmen yang baru.
func (r *AssessmentRepository) SubmitAnswers(userID uuid.UUID, answers []models.UserAnswer) (uuid.UUID, error) {
	tx, err := r.db.BeginTx(context.Background(), nil)
	if err != nil {
		return uuid.Nil, err
	}
	defer tx.Rollback()

	var assessmentID uuid.UUID
	queryCreateAssessment := `
		INSERT INTO user_assessments (user_id, status, completed_at)
		VALUES ($1, 'completed', now())
		RETURNING id`

	err = tx.QueryRow(queryCreateAssessment, userID).Scan(&assessmentID)
	if err != nil {
		return uuid.Nil, err
	}

	stmt, err := tx.Prepare(`
		INSERT INTO user_assessment_answers (user_assessment_id, question_id, option_id)
		VALUES ($1, $2, $3)`)
	if err != nil {
		return uuid.Nil, err
	}
	defer stmt.Close()

	for _, answer := range answers {
		_, err := stmt.Exec(assessmentID, answer.QuestionID, answer.OptionID)
		if err != nil {
			return uuid.Nil, err
		}
	}

	return assessmentID, tx.Commit()
}

// CalculateScores menghitung skor untuk setiap jalur karier berdasarkan jawaban pengguna.
func (r *AssessmentRepository) CalculateScores(assessmentID uuid.UUID) ([]models.CareerRecommendation, error) {
	query := `
		SELECT
			cp.id,
			cp.name,
			cp.description,
			COALESCE(SUM(ao.weight), 0)::INTEGER AS total_score
		FROM
			career_paths cp
		LEFT JOIN
			assessment_options ao ON cp.id = ao.career_path_id
		LEFT JOIN
			user_assessment_answers uaa ON ao.id = uaa.option_id
		WHERE
			uaa.user_assessment_id = $1
		GROUP BY
			cp.id, cp.name, cp.description
		ORDER BY
			total_score DESC;
	`

	rows, err := r.db.Query(query, assessmentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var recommendations []models.CareerRecommendation
	for rows.Next() {
		var rec models.CareerRecommendation
		if err := rows.Scan(&rec.CareerPathID, &rec.CareerPathName, &rec.CareerDescription, &rec.MatchScore); err != nil {
			return nil, err
		}
		recommendations = append(recommendations, rec)
	}

	return recommendations, nil
}

func (r *AssessmentRepository) SaveRecommendations(assessmentID uuid.UUID, userID uuid.UUID, recommendations []models.CareerRecommendation) error {
	// Konversi slice of structs ke JSONB
	recJSON, err := json.Marshal(recommendations)
	if err != nil {
		return err
	}

	query := `
		INSERT INTO user_career_recommendations (user_assessment_id, user_id, recommendations)
		VALUES ($1, $2, $3)`

	_, err = r.db.Exec(query, assessmentID, userID, recJSON)
	return err
}