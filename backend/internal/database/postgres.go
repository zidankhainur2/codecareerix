package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/jackc/pgx/v5/stdlib"
)

// Connect membuat dan memverifikasi koneksi ke database menggunakan DSN yang diberikan.
func Connect(dsn string) (*sql.DB, error) {
	// Buka koneksi
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("gagal membuka koneksi database: %w", err)
	}

	// Cek koneksi
	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("gagal terhubung ke database: %w", err)
	}

	log.Println("✅ Koneksi ke database PostgreSQL berhasil!")
	return db, nil
}