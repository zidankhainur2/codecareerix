package repositories

import (
	"database/sql"

	"github.com/zidankhainur2/codecareerix/backend/internal/models"
)

type JobRepository struct {
	db *sql.DB
}

func NewJobRepository(db *sql.DB) *JobRepository {
	return &JobRepository{db: db}
}

// GetAllJobs mengambil semua lowongan kerja dengan paginasi.
func (r *JobRepository) GetAllJobs(limit, offset int) ([]models.JobPosting, error) {
	query := `
		SELECT id, title, company_name, location, description_clean, source_url, job_type, scraped_at
		FROM job_postings
		ORDER BY scraped_at DESC
		LIMIT $1 OFFSET $2`

	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var jobs []models.JobPosting
	for rows.Next() {
		var j models.JobPosting
		if err := rows.Scan(&j.ID, &j.Title, &j.CompanyName, &j.Location, &j.DescriptionClean, &j.SourceURL, &j.JobType, &j.ScrapedAt); err != nil {
			return nil, err
		}
		jobs = append(jobs, j)
	}
	return jobs, nil
}