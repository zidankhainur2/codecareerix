package handlers

import (
	"database/sql"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

func (h *UserHandler) GetProfile(c *gin.Context) {
	// 1. Ambil userID dari context yang di-set oleh middleware
	userIDString, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Pengguna tidak terotentikasi"})
		return
	}

	// 2. Konversi userID ke tipe UUID
	userID, err := uuid.Parse(userIDString.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID pengguna tidak valid"})
		return
	}

	// 3. Panggil repository untuk mendapatkan data pengguna
	user, err := h.repo.GetByID(userID)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Pengguna tidak ditemukan"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil profil"})
		return
	}

	// 4. Kirim data pengguna sebagai respons
	c.JSON(http.StatusOK, user)
}

func (h *UserHandler) SelectCareerPath(c *gin.Context) {
	var input models.SelectCareerPathInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userIDString, _ := c.Get("userID")
	userID, err := uuid.Parse(userIDString.(string))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "ID pengguna tidak valid"})
		return
	}

	// Panggil repository untuk update data pengguna
	err = h.repo.UpdateActiveCareerPath(userID, input.CareerPathID)
	if err != nil {
		log.Printf("Gagal memilih jalur karier untuk user %s: %v", userID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memilih jalur karier"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Jalur karier berhasil dipilih"})
}