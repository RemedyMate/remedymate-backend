package auth

import (
	"errors"
	"fmt"
	"os"
	"time"

	"remedymate-backend/domain/entities"

	"github.com/golang-jwt/jwt/v5"
)

type JWTService struct {
	secretKey string
	expiry    time.Duration
}

type Claims struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	jwt.RegisteredClaims
}

func NewJWTService() *JWTService {
	expiryHour := 24
	if hours := os.Getenv("JWT_EXPIRY_HOURS"); hours != "" {
		if h, err := time.ParseDuration(hours + "h"); err == nil {
			expiryHour = int(h.Hours())
		}
	}

	return &JWTService{
		secretKey: os.Getenv("JWT_SECRET_KEY"),
		expiry:    time.Duration(expiryHour) * time.Hour,
	}
}

func (s *JWTService) GenerateToken(user entities.User) (string, error) {
	if s.secretKey == "" {
		return "", errors.New("JWT_SECRET_KEY is not set")
	}

	claims := &Claims{
		UserID:   user.ID,
		Username: user.Username,
		Email:    user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.expiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "remedymate",
			Subject:   user.ID,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.secretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (s *JWTService) ValidateToken(tokenString string) (*Claims, error) {
	if s.secretKey == "" {
		return nil, errors.New("JWT_SECRET_KEY is not set")
	}

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

func (s *JWTService) RefreshToken(tokenString string) (string, error) {
	claims, err := s.ValidateToken(tokenString)
	if err != nil {
		return "", err
	}

	newClaims := &Claims{
		UserID:   claims.UserID,
		Username: claims.Username,
		Email:    claims.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.expiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "remedymate",
			Subject:   claims.UserID,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, newClaims)
	newTokenString, err := token.SignedString([]byte(s.secretKey))
	if err != nil {
		return "", err
	}

	return newTokenString, nil
}
