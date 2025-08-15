package models

import (
	"time"

	"github.com/google/uuid"
)

// UserProject merepresentasikan data di tabel 'user_projects'.
type UserProject struct {
	ID            uuid.UUID `json:"id"`
	UserID        uuid.UUID `json:"user_id"`
	ResourceID    *int      `json:"resource_id,omitempty"`
	Title         string    `json:"title"`
	Description   *string   `json:"description,omitempty"`
	ProjectURL    *string   `json:"project_url,omitempty"`
	CoverImageURL *string   `json:"cover_image_url,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// CreateProjectInput adalah struct untuk input saat membuat proyek baru.
type CreateProjectInput struct {
	Title         string  `json:"title" binding:"required"`
	Description   *string `json:"description"`
	ProjectURL    *string `json:"project_url" binding:"omitempty,url"`
	CoverImageURL *string `json:"cover_image_url" binding:"omitempty,url"`
	ResourceID    *int    `json:"resource_id"` // Jika proyek ini terhubung ke materi belajar
}