package models

import (
	"time"

	"github.com/google/uuid"
)

// User merepresentasikan struktur data di tabel 'users'
type User struct {
	ID                uuid.UUID `json:"id"`
	FullName          string    `json:"full_name"`
	Email             string    `json:"email"`
	PasswordHash      string    `json:"-"` // Tanda '-' artinya jangan pernah kirim field ini dalam response JSON
	ProfilePictureURL string    `json:"profile_picture_url"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

// RegisterUserInput adalah struct untuk validasi input dari request registrasi
type RegisterUserInput struct {
	FullName string `json:"full_name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type LoginUserInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}