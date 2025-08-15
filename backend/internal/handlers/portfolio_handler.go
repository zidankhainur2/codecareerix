package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/zidankhainur2/codecareerix/backend/internal/models"
	"github.com/zidankhainur2/codecareerix/backend/internal/repositories"
)

type PortfolioHandler struct {
	repo *repositories.PortfolioRepository
}

func NewPortfolioHandler(repo *repositories.PortfolioRepository) *PortfolioHandler {
	return &PortfolioHandler{repo: repo}
}

// CreateProject menangani pembuatan proyek baru.
func (h *PortfolioHandler) CreateProject(c *gin.Context) {
	var input models.CreateProjectInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, _ := c.Get("userID")
	
	project, err := h.repo.CreateProject(userID.(uuid.UUID), input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat proyek"})
		return
	}

	c.JSON(http.StatusCreated, project)
}

// GetMyProjects menangani pengambilan semua proyek pengguna.
func (h *PortfolioHandler) GetMyProjects(c *gin.Context) {
	userID, _ := c.Get("userID")

	projects, err := h.repo.GetProjectsByUserID(userID.(uuid.UUID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil proyek"})
		return
	}

	c.JSON(http.StatusOK, projects)
}