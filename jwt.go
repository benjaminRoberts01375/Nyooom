package main

import (
	"context"
	"nyooom/logging"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const CookieName = "nyooom-session-token"
const LoginDuration = time.Hour*24*6 + time.Hour*12 // 6.5 days

type Claims struct {
	jwt.RegisteredClaims
}

// TimeFunc allows mocking time in tests
type TimeFunc func() time.Time

// JWTService handles JWT operations with dependency injection
type JWTService struct {
	secret   []byte
	timeFunc TimeFunc
}

// NewJWTService creates a new JWT service with the provided secret
func NewJWTService(secret string) *JWTService {
	return &JWTService{
		secret:   []byte(secret),
		timeFunc: time.Now,
	}
}

// withTimeFunc allows overriding the time function (useful for testing)
func (s *JWTService) withTimeFunc(timeFunc TimeFunc) *JWTService {
	s.timeFunc = timeFunc
	return s
}

// generateJWT creates a new JWT token with the specified duration
func (s *JWTService) generateJWT(duration time.Duration) (string, error) {
	now := s.timeFunc()
	claims := Claims{}
	claims.ExpiresAt = jwt.NewNumericDate(now.Add(duration))
	claims.IssuedAt = jwt.NewNumericDate(now)
	claims.NotBefore = jwt.NewNumericDate(now)
	claims.Issuer = "Backend API"
	claims.Subject = "Session Token"
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(s.secret)
}

// validateJWT checks if a token string is valid and returns the claims
func (s *JWTService) validateJWT(tokenString string) (*Claims, bool) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (any, error) {
		return s.secret, nil
	})
	if err != nil {
		return nil, false
	}
	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, false
	}
	return claims, claims.ExpiresAt.After(s.timeFunc())
}

func loadJWTSecret(db AdvancedDB) *JWTService {
	// Check if the JWT secret was passed in as an environment variable
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret != "" {
		logging.Println("JWT provided as an environment variable")
		return NewJWTService(jwtSecret)
	}

	// Read the JWT secret from database
	jwtSecret, err := db.GetJWT(context.Background())
	if err == nil && jwtSecret != "" {
		logging.Println("JWT provided in database")
		return NewJWTService(jwtSecret)
	}

	// Generate a new JWT secret
	logging.Println("Generating JWT and storing in database")
	newJWT := generateRandomString(15)
	err = db.SetJWT(context.Background(), newJWT)
	if err != nil {
		panic("Failed to save JWT to database: " + err.Error())
	}
	return NewJWTService(newJWT)
}
