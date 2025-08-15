package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/zidankhainur2/codecareerix/backend/internal/config"
	"github.com/zidankhainur2/codecareerix/backend/internal/database"
	"github.com/zidankhainur2/codecareerix/backend/internal/server"
)

func main() {

	 _ = godotenv.Load()
	// 1. Muat Konfigurasi dari Environment
	cfg, err := config.New()
	if err != nil {
		log.Fatalf("Gagal memuat konfigurasi: %v", err)
	}

	// 2. Hubungkan ke Database menggunakan DSN dari config
	db, err := database.Connect(cfg.DB.DSN)
	if err != nil {
		log.Fatalf("Gagal terhubung ke database: %v", err)
	}
	defer db.Close()

	// 3. Buat dan Jalankan Server
	srv := server.New(db, cfg)
	if err := srv.Run(cfg.Port); err != nil {
		log.Fatalf("Gagal menjalankan server: %v", err)
	}
}