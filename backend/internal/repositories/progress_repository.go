package repositories

import (
	"database/sql"

	"github.com/google/uuid"
)

type ProgressRepository struct {
	db *sql.DB
}

func NewProgressRepository(db *sql.DB) *ProgressRepository {
	return &ProgressRepository{db: db}
}

// MarkResourceAsComplete menandai sebuah resource sebagai 'completed' untuk seorang user.
func (r *ProgressRepository) MarkResourceAsComplete(userID uuid.UUID, resourceID int) error {
	query := `
		INSERT INTO user_progress (user_id, resource_id, status, completed_at)
		VALUES ($1, $2, 'completed', now())
		ON CONFLICT (user_id, resource_id) 
		DO UPDATE SET 
			status = EXCLUDED.status, 
			completed_at = EXCLUDED.completed_at;
	`
	// ON CONFLICT menangani kasus jika user sudah pernah memulai resource ini.
	// EXCLUDED merujuk pada nilai yang akan dimasukkan jika tidak ada konflik.

	_, err := r.db.Exec(query, userID, resourceID)
	return err
}