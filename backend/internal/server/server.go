package server

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zidankhainur2/codecareerix/backend/internal/config"
	"github.com/zidankhainur2/codecareerix/backend/internal/handlers"
	"github.com/zidankhainur2/codecareerix/backend/internal/repositories"
)

// Server adalah struct untuk server HTTP kita.
type Server struct {
	router *gin.Engine
	db     *sql.DB
	cfg    *config.Config
}

// New membuat instance Server baru.
func New(db *sql.DB) *Server {
	router := gin.Default()
	s := &Server{
		router: router,
		db:     db,
		cfg:    cfg,
	}
	s.registerRoutes()
	return s
}

// registerRoutes mendaftarkan semua rute untuk aplikasi.
func (s *Server) registerRoutes() {
	// Rute root
	s.router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Selamat datang di CodeCareerix API!"})
	})

	// Membuat instance health handler dan mendaftarkan rutenya
	healthHandler := handlers.NewHealthHandler(s.db)
	s.router.GET("/health", healthHandler.Check)

	// --- RUTE PENGGUNA ---
	userRepo := repositories.NewUserRepository(s.db)
	// Teruskan JWT secret dari config ke handler
	userHandler := handlers.NewUserHandler(userRepo, s.cfg.JWTSecret)

	s.router.POST("/users/register", userHandler.Register)
	s.router.POST("/users/login", userHandler.Login) // Daftarkan rute baru
}

// Run menjalankan server HTTP pada port yang diberikan.
func (s *Server) Run(port string) error {
	log.Printf("Server berjalan di http://localhost:%s\n", port)
	return s.router.Run(":" + port)
}