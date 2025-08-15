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
	"github.com/zidankhainur2/codecareerix/backend/internal/services" // Import services
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
	// --- RUTE PUBLIK ---
	s.router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Selamat datang di CodeCareerix API!"})
	})
	healthHandler := handlers.NewHealthHandler(s.db)
	s.router.GET("/health", healthHandler.Check)

	userRepo := repositories.NewUserRepository(s.db)
	userHandler := handlers.NewUserHandler(userRepo, s.cfg.JWTSecret)
	s.router.POST("/users/register", userHandler.Register)
	s.router.POST("/users/login", userHandler.Login)

	// --- RUTE TERPROTEKSI ---
	authRoutes := s.router.Group("/")
	authRoutes.Use(auth.AuthMiddleware(s.cfg.JWTSecret))
	{
		// Rute User terproteksi
		authRoutes.GET("/users/profile", userHandler.GetProfile)
		authRoutes.POST("/users/career-path", userHandler.SelectCareerPath)

		// Rute Asesmen terproteksi
		assessmentRepo := repositories.NewAssessmentRepository(s.db)
		assessmentService := services.NewAssessmentService(assessmentRepo)
		assessmentHandler := handlers.NewAssessmentHandler(assessmentRepo, assessmentService)
		authRoutes.GET("/assessments", assessmentHandler.GetAssessmentQuestions)
		authRoutes.POST("/assessments/submit", assessmentHandler.SubmitAssessment)

		// RUTE LEARNING/ROADMAP BARU
		learningRepo := repositories.NewLearningRepository(s.db)
		learningHandler := handlers.NewLearningHandler(learningRepo, userRepo) // Gunakan userRepo yang sudah ada
		authRoutes.GET("/roadmap", learningHandler.GetMyRoadmap)
	
		// RUTE PROGRESS BARU
		progressRepo := repositories.NewProgressRepository(s.db)
		progressHandler := handlers.NewProgressHandler(progressRepo)
		authRoutes.POST("/progress/resource/:resource_id", progressHandler.MarkAsComplete)
	
		// RUTE PORTFOLIO BARU
		portfolioRepo := repositories.NewPortfolioRepository(s.db)
		portfolioHandler := handlers.NewPortfolioHandler(portfolioRepo)
		authRoutes.POST("/portfolio/projects", portfolioHandler.CreateProject)
		authRoutes.GET("/portfolio/projects", portfolioHandler.GetMyProjects)
	
		// RUTE LOWONGAN KERJA BARU
		jobRepo := repositories.NewJobRepository(s.db)
		jobHandler := handlers.NewJobHandler(jobRepo)
		authRoutes.GET("/jobs", jobHandler.GetAllJobs)
	}
}

func (s *Server) Run(port string) error {
	log.Printf("🚀 Server berjalan di http://localhost:%s\n", port)
	return s.router.Run(":" + port)
}