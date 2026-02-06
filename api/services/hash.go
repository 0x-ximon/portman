package services

import (
	"context"
	"crypto/rand"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/0x-ximon/portman/api/repositories"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type ClaimsKey struct{}
type UserKey struct{}

type Claims struct {
	jwt.RegisteredClaims
	ID            uuid.UUID `json:"id"`
	EmailAddress  string    `json:"email"`
	WalletAddress string    `json:"wallet"`
}

func GetIDFromContext(ctx context.Context) (uuid.UUID, bool) {
	if user, ok := ctx.Value(UserKey{}).(*repositories.User); ok {
		return user.ID, true
	}

	if claims, ok := ctx.Value(ClaimsKey{}).(*Claims); ok {
		return claims.ID, true
	}

	return uuid.Nil, false
}

func GenerateOTP(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	var otp strings.Builder
	for i := range length {
		digit := bytes[i] % 10
		fmt.Fprintf(&otp, "%d", digit)
	}

	return otp.String(), nil
}

func GenerateJWT(id uuid.UUID) (string, error) {
	jwtSecret, ok := os.LookupEnv("JWT_SECRET")
	if !ok {
		return "", fmt.Errorf("JWT_SECRET environment variable not set")
	}

	t := time.Now().Add(24 * time.Hour)
	if os.Getenv("ENV") == "dev" {
		t.Add(10e3 * time.Hour)
	}

	expirationTime := jwt.NewNumericDate(t)
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

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func ValidateHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
