package server

import (
	"database/sql"
	"log"
	"net/http"

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
	s.router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Selamat datang di CodeCareerix API!"})
	})
	healthHandler := handlers.NewHealthHandler(s.db)
	s.router.GET("/health", healthHandler.Check)

	userRepo := repositories.NewUserRepository(s.db)
	userHandler := handlers.NewUserHandler(userRepo, s.cfg.JWTSecret)
	s.router.POST("/users/register", userHandler.Register)
	s.router.POST("/users/login", userHandler.Login)

	// --- RUTE TERPROTEKSI (Butuh login) ---
	authRoutes := s.router.Group("/")
	authRoutes.Use(auth.AuthMiddleware(s.cfg.JWTSecret))
	{
		// Rute User terproteksi
		authRoutes.GET("/users/profile", userHandler.GetProfile)
		
		// Rute Asesmen terproteksi
		assessmentRepo := repositories.NewAssessmentRepository(s.db)
		assessmentHandler := handlers.NewAssessmentHandler(assessmentRepo)
		authRoutes.GET("/assessments", assessmentHandler.GetAssessmentQuestions)
	}
}

// Run menjalankan server HTTP pada port yang diberikan.
func (s *Server) Run(port string) error {
	log.Printf("🚀 Server berjalan di http://localhost:%s\n", port)
	return s.router.Run(":" + port)
}