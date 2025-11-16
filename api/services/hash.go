package services

import (
	"crypto/rand"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Claims struct {
	jwt.RegisteredClaims
	ID uuid.UUID `json:"id"`
}

type TokenKey struct{}

func GenerateOTP(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	var otp string
	for i := 0; i < length; i++ {
		digit := bytes[i] % 10
		otp += fmt.Sprintf("%d", digit)
	}

	return otp, nil
}

func GenerateJWT(id uuid.UUID) (string, error) {
	jwtSecret, ok := os.LookupEnv("JWT_SECRET")
	if !ok {
		return "", fmt.Errorf("JWT_SECRET environment variable not set")
	}

	expirationTime := jwt.NewNumericDate(time.Now().Add(24 * time.Hour))
	claims := &Claims{
		ID: id,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: expirationTime,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtSecret))
}

func ValidateJWT(tokenString string) (*Claims, error) {
	jwtKey := os.Getenv("JWT_SECRET")
	if jwtKey == "" {
		return nil, fmt.Errorf("JWT_SECRET environment variable not set")
	}

	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
		return []byte(jwtKey), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("token is invalid")
	}

	return claims, nil
}

func GetToken(r *http.Request) (string, bool) {
	token, ok := r.Context().Value(TokenKey{}).(string)
	return token, ok
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func ValidateHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
