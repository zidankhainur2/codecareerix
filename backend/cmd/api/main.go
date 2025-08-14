package main

import (
	"log"

	"github.com/zidankhainur2/codecareerix/backend/internal/config"   // Ganti dengan path modul Anda
	"github.com/zidankhainur2/codecareerix/backend/internal/database" // Ganti dengan path modul Anda
	"github.com/zidankhainur2/codecareerix/backend/internal/server"   // Ganti dengan path modul Anda
)

func main() {
	// 1. Muat Konfigurasi
	cfg := config.New()

	// 2. Hubungkan ke Database
	db, err := database.Connect()
if err != nil {
	log.Fatalf("Database error: %v", err)
}
defer db.Close()

	// 3. Buat dan Jalankan Server
	srv := server.New(db, cfg) // Perbarui baris ini
	if err := srv.Run(cfg.Port); err != nil {
		log.Fatalf("Gagal menjalankan server: %v", err)
	}
}