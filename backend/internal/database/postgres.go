package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
)

// Connect membuat dan memverifikasi koneksi ke database PostgreSQL
func Connect() (*sql.DB, error) {
	// Ambil konfigurasi dari environment variables
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	if dbUser == "" || dbPassword == "" || dbHost == "" || dbPort == "" || dbName == "" {
		log.Fatal("Environment variables DB_USER, DB_PASSWORD, DB_HOST, DB_PORT, DB_NAME harus di-set")
	}

	// Format DSN untuk driver pgx
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		dbUser,
		dbPassword,
		dbHost,
		dbPort,
		dbName,
	)

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
