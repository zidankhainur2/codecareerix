package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/zidankhainur2/codecareerix/backend/internal/repositories"
)

type JobHandler struct {
	repo *repositories.JobRepository
}

func NewJobHandler(repo *repositories.JobRepository) *JobHandler {
	return &JobHandler{repo: repo}
}

// GetAllJobs menangani pengambilan semua lowongan kerja.
func (h *JobHandler) GetAllJobs(c *gin.Context) {
	// Paginasi sederhana
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset := (page - 1) * limit

	jobs, err := h.repo.GetAllJobs(limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data lowongan kerja"})
		return
	}

	c.JSON(http.StatusOK, jobs)
}