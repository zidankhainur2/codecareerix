package repositories

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
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

func (r *UserRepository) GetByID(id uuid.UUID) (*models.User, error) {
	query := `
		SELECT id, full_name, email, password_hash, profile_picture_url, active_career_path_id, created_at, updated_at
		FROM users
		WHERE id = $1`

	var user models.User
	err := r.db.QueryRow(query, id).Scan(
		&user.ID,
		&user.FullName,
		&user.Email,
		&user.PasswordHash,
		&user.ProfilePictureURL,
		&user.ActiveCareerPathID,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) UpdateActiveCareerPath(userID uuid.UUID, careerPathID int) error {
	query := `
		UPDATE users
		SET active_career_path_id = $1, updated_at = now()
		WHERE id = $2`

	_, err := r.db.Exec(query, careerPathID, userID)
	return err
}