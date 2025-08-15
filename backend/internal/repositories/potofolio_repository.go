package repositories

import (
	"database/sql"

	"github.com/google/uuid"
	"github.com/zidankhainur2/codecareerix/backend/internal/models"
)

type PortfolioRepository struct {
	db *sql.DB
}

func NewPortfolioRepository(db *sql.DB) *PortfolioRepository {
	return &PortfolioRepository{db: db}
}

// CreateProject membuat proyek baru untuk seorang pengguna.
func (r *PortfolioRepository) CreateProject(userID uuid.UUID, input models.CreateProjectInput) (*models.UserProject, error) {
	query := `
		INSERT INTO user_projects (user_id, resource_id, title, description, project_url, cover_image_url)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at, updated_at`

	var project models.UserProject
	err := r.db.QueryRow(
		query,
		userID,
		input.ResourceID,
		input.Title,
		input.Description,
		input.ProjectURL,
		input.CoverImageURL,
	).Scan(&project.ID, &project.CreatedAt, &project.UpdatedAt)

	if err != nil {
		return nil, err
	}

	project.UserID = userID
	project.Title = input.Title
	project.Description = input.Description
	project.ProjectURL = input.ProjectURL
	project.CoverImageURL = input.CoverImageURL
	project.ResourceID = input.ResourceID

	return &project, nil
}

// GetProjectsByUserID mengambil semua proyek milik seorang pengguna.
func (r *PortfolioRepository) GetProjectsByUserID(userID uuid.UUID) ([]models.UserProject, error) {
	query := `
		SELECT id, resource_id, title, description, project_url, cover_image_url, created_at, updated_at
		FROM user_projects
		WHERE user_id = $1
		ORDER BY created_at DESC`
	
	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []models.UserProject
	for rows.Next() {
		var p models.UserProject
		p.UserID = userID
		if err := rows.Scan(&p.ID, &p.ResourceID, &p.Title, &p.Description, &p.ProjectURL, &p.CoverImageURL, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		projects = append(projects, p)
	}

	return projects, nil
}