package jwtutil

import (
	"log"
	"os"
	"time"

	AppError "remedymate-backend/domain/AppError"
	"remedymate-backend/domain/entities"

	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Claims struct {
	TokenID  string        `json:"id"`
	UserID   string        `json:"user_id"`
	Username string        `json:"username"`
	Email    string        `json:"email"`
	Role     entities.Role `json:"role"`
	jwt.RegisteredClaims
}

func getJWTExpiry(isAccessToken bool) time.Duration {
	var timeSpan = os.Getenv("REFRESH_EXPIRY_DAYS")
	if isAccessToken {
		timeSpan = os.Getenv("ACCESS_EXPIRY_MINUTES")
	}
	if timeSpan != "" && isAccessToken {
		h, err := time.ParseDuration(timeSpan + "m")
		if err == nil {
			return h
		} else {
			log.Println("❌ Invalid ACCESS_EXPIRY_MINUTES format, using default 30 minutes")
		}
	} else if timeSpan != "" && !isAccessToken {
		d, err := time.ParseDuration(timeSpan + "h")
		if err == nil {
			return d * 24
		} else {
			log.Println("❌ Invalid REFRESH_EXPIRY_DAYS format, using default 7 days")
		}
	}

	if !isAccessToken {
		return 7 * 24 * time.Hour // default to 7 days for refresh tokens
	}
	return 30 * time.Minute // default to 24 hours
}

func getJWTSecret(isAccessToken bool) string {
	var jwt_secret = os.Getenv("JWT_REFRESH_SECRET_KEY")
	if isAccessToken {
		jwt_secret = os.Getenv("JWT_ACCESS_SECRET_KEY")
	}
	if jwt_secret == "" {
		log.Println("❌ JWT_SECRET_KEY is not set")
		jwt_secret = "default_secret" // default secret for development
	}
	return jwt_secret
}

func GenerateAccessToken(user *entities.User) (string, error) {
	claims := &Claims{
		TokenID:  primitive.NewObjectID().Hex(),
		UserID:   user.ID,
		Username: user.Username,
		Email:    user.Email,
		Role:     user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(getJWTExpiry(true))),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "remedymate",
			Subject:   user.ID,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(getJWTSecret(true)))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func GenerateRefreshToken(user *entities.User) (*entities.RefreshToken, error) {
	claims := &Claims{
		TokenID:  primitive.NewObjectID().Hex(),
		UserID:   user.ID,
		Username: user.Username,
		Email:    user.Email,
		Role:     user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(getJWTExpiry(false))),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "remedymate",
			Subject:   user.ID,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(getJWTSecret(false)))
	if err != nil {
		return nil, err
	}

	return &entities.RefreshToken{
		ID:        claims.TokenID,
		Token:     tokenString,
		UserID:    user.ID,
		ExpiresAt: claims.ExpiresAt.Time,
		CreatedAt: claims.IssuedAt.Time,
	}, nil
}

func ValidateToken(tokenString string, isAccessToken bool) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			log.Printf("unexpected signing method: %v", token.Header["alg"])
			return nil, AppError.ErrInternalServer
		}
		return []byte(getJWTSecret(isAccessToken)), nil
	})

	if err != nil {
		log.Printf("error parsing token: %v", err)
		return nil, AppError.ErrInvalidToken
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, AppError.ErrInvalidToken
}
