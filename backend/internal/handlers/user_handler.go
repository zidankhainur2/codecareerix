package handlers

import (
	"database/sql"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	"github.com/zidankhainur2/codecareerix/backend/internal/auth"
	"github.com/zidankhainur2/codecareerix/backend/internal/models"       // Ganti dengan path modul Anda
	"github.com/zidankhainur2/codecareerix/backend/internal/repositories" // Ganti dengan path modul Anda
	"golang.org/x/crypto/bcrypt"
)

type UserHandler struct {
	repo *repositories.UserRepository
	jwtSecret  string
}

func NewUserHandler(repo *repositories.UserRepository, jwtSecret string) *UserHandler {
	return &UserHandler{
		repo:      repo,
		jwtSecret: jwtSecret, // Tambahkan ini
	}
}

func (h *UserHandler) Register(c *gin.Context) {
	var input models.RegisterUserInput

	// 1. Bind & Validasi Input
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 2. Hash Password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memproses password"})
		return
	}

	// 3. Buat User baru
	newUser := models.User{
		FullName:     input.FullName,
		Email:        strings.ToLower(input.Email),
		PasswordHash: string(hashedPassword),
	}

	// 4. Simpan ke Database
	if err := h.repo.Create(&newUser); err != nil {
		// Cek apakah error karena email duplikat
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			c.JSON(http.StatusConflict, gin.H{"error": "Email sudah terdaftar"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat pengguna"})
		return
	}

	// 5. Kirim Respons Sukses
	c.JSON(http.StatusCreated, newUser)
}

func (h *UserHandler) Login(c *gin.Context) {
	var input models.LoginUserInput

	// 1. Bind & Validasi Input
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 2. Cari pengguna berdasarkan email
	user, err := h.repo.GetByEmail(strings.ToLower(input.Email))
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Email atau password salah"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mencari pengguna"})
		return
	}

	// 3. Bandingkan password
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password))
	if err != nil {
		// Password tidak cocok
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Email atau password salah"})
		return
	}

	// 4. Buat JWT
	token, err := auth.GenerateJWT(user.ID, h.jwtSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat token"})
		return
	}

	// 5. Kirim Respons Sukses
	c.JSON(http.StatusOK, gin.H{"token": token})
}