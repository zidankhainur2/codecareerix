package config

// Config menampung semua konfigurasi aplikasi.
// Nantinya ini bisa dibaca dari file .env atau environment variables.
type Config struct {
	Port string
	DB   struct {
		DSN string // Data Source Name
	}
	JWTSecret string
}

// New mengembalikan instance konfigurasi dengan nilai default.
func New() *Config {
	cfg := &Config{}
	cfg.Port = "8080"
	cfg.DB.DSN = "postgres://user:password@localhost:5432/codecareerix_db?sslmode=disable"
	cfg.JWTSecret = "c80c86cf88c3deea1bab18e46cb97a1a00205a70c832e56044b6704c59e5eb70"
	return cfg
}