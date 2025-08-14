package repositories

import (
	"database/sql"
	"time"

	"github.com/zidankhainur2/codecareerix/backend/internal/models" // Ganti dengan path modul Anda
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

// Create menyisipkan user baru ke dalam database
func (r *UserRepository) Create(user *models.User) error {
	query := `
		INSERT INTO users (full_name, email, password_hash, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at`

	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	err := r.db.QueryRow(
		query,
		user.FullName,
		user.Email,
		user.PasswordHash,
		user.CreatedAt,
		user.UpdatedAt,
	).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)

	return err
}

func (r *UserRepository) GetByEmail(email string) (*models.User, error) {
	query := `
		SELECT id, full_name, email, password_hash, profile_picture_url, created_at, updated_at
		FROM users
		WHERE email = $1`

	var user models.User
	err := r.db.QueryRow(query, email).Scan(
		&user.ID,
		&user.FullName,
		&user.Email,
		&user.PasswordHash,
		&user.ProfilePictureURL,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err // Akan mengembalikan sql.ErrNoRows jika tidak ditemukan
	}

	return &user, nil
}