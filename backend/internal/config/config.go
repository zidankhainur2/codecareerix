package config

import (
	"fmt"
	"os"
)

// Config menampung semua konfigurasi aplikasi yang dimuat dari environment.
type Config struct {
	Port      string
	JWTSecret string
	DB        struct {
		DSN string // Data Source Name
	}
}

// New memuat konfigurasi dari environment variables dan mengembalikannya.
func New() (*Config, error) {
	// Ambil port, berikan nilai default "8080" jika tidak ada
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Ambil JWT Secret, ini wajib ada
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		return nil, fmt.Errorf("environment variable JWT_SECRET tidak boleh kosong")
	}

	// Ambil konfigurasi database dari environment
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	// Pastikan semua variabel database ada
	if dbUser == "" || dbPassword == "" || dbHost == "" || dbPort == "" || dbName == "" {
		return nil, fmt.Errorf("semua environment variable database (DB_USER, DB_PASSWORD, DB_HOST, DB_PORT, DB_NAME) harus di-set")
	}

	// Buat DSN (Data Source Name)
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		dbUser,
		dbPassword,
		dbHost,
		dbPort,
		dbName,
	)

	// Buat dan kembalikan struct Config
	cfg := &Config{
		Port:      port,
		JWTSecret: jwtSecret,
	}
	cfg.DB.DSN = dsn

	return cfg, nil
}