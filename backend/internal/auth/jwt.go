package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// GenerateJWT membuat token JWT baru untuk user ID tertentu.
func GenerateJWT(userID uuid.UUID, secretKey string) (string, error) {
	// Menetapkan claims untuk token
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(), // Token berlaku selama 24 jam
		"iat":     time.Now().Unix(),                      // Waktu token dibuat
	}

	// Membuat token dengan claims dan metode signing HS256
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Menandatangani token dengan secret key
	signedToken, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}