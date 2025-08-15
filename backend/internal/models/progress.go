package models

import (
	"time"

	"github.com/google/uuid"
)

// UserProgress merepresentasikan data di tabel 'user_progress'.
type UserProgress struct {
	ID          int64
	UserID      uuid.UUID
	ResourceID  int
	Status      string
	CompletedAt *time.Time // Pointer agar bisa NULL
}