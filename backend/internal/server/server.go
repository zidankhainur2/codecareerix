package server

import (
	"database/sql"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/zidankhainur2/codecareerix/backend/internal/auth"
	"github.com/zidankhainur2/codecareerix/backend/internal/config"
	"github.com/zidankhainur2/codecareerix/backend/internal/handlers"
	"github.com/zidankhainur2/codecareerix/backend/internal/repositories"
)

type Server struct {
	router *gin.Engine
	db     *sql.DB
	cfg    *config.Config
}

func New(db *sql.DB, cfg *config.Config) *Server {
	router := gin.Default()
	s := &Server{
		router: router,
		db:     db,
		cfg:    cfg,
	}
	s.registerRoutes()
	return s
}

func (s *Server) registerRoutes() {
	// --- RUTE PUBLIK (Tidak butuh login) ---
	s.router.GET("/", func(c *gin.Context) { /* ... */ })
	healthHandler := handlers.NewHealthHandler(s.db)
	s.router.GET("/health", healthHandler.Check)

	userRepo := repositories.NewUserRepository(s.db)
	userHandler := handlers.NewUserHandler(userRepo, s.cfg.JWTSecret)
	s.router.POST("/users/register", userHandler.Register)
	s.router.POST("/users/login", userHandler.Login)

	// --- RUTE TERPROTEKSI (Butuh login) ---
	// Buat grup baru untuk rute yang membutuhkan otentikasi
	authRoutes := s.router.Group("/")
	authRoutes.Use(auth.AuthMiddleware(s.cfg.JWTSecret))
	{
		// Semua rute di dalam blok ini akan dilindungi oleh middleware
		authRoutes.GET("/users/profile", userHandler.GetProfile)
		// Tambahkan rute terproteksi lainnya di sini...
	}
}

// Run menjalankan server HTTP pada port yang diberikan.
func (s *Server) Run(port string) error {
	log.Printf("🚀 Server berjalan di http://localhost:%s\n", port)
	return s.router.Run(":" + port)
}