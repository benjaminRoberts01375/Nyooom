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

type JWTService struct {
	secret   []byte
	timeFunc TimeFunc
}

func NewJWTService(secret string, timeGenerator TimeFunc) JWTService {
	return JWTService{
		secret:   []byte(secret),
		timeFunc: timeGenerator,
	}
}

func (s *JWTService) GenerateJWT(duration time.Duration) (string, error) {
	now := s.timeFunc()
	claims := Claims{}
	claims.ExpiresAt = jwt.NewNumericDate(now.Add(duration))
	claims.IssuedAt = jwt.NewNumericDate(now)
	claims.NotBefore = jwt.NewNumericDate(now)
	claims.Issuer = "Backend API"
	claims.Subject = "Session Token"
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(s.secret)
}

func (s *JWTService) ValidateJWT(tokenString string) (*Claims, bool) {
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

func loadJWTSecret(db AdvancedDB) JWTService {
	// Check if the JWT secret was passed in as an environment variable
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret != "" {
		logging.Println("JWT provided as an environment variable")
		return NewJWTService(jwtSecret, time.Now)
	}

	// Read the JWT secret from database
	jwtSecret, err := db.GetJWTSecret(context.Background())
	if err == nil && jwtSecret != "" {
		logging.Println("JWT secret provided in database")
		return NewJWTService(jwtSecret, time.Now)
	}

	// Generate a new JWT secret
	logging.Println("Generating JWT secret and storing in database")
	newJWT := generateRandomString(15)
	err = db.SetJWTSecret(context.Background(), newJWT)
	if err != nil {
		panic("Failed to save JWT secret to database: " + err.Error())
	}
	return NewJWTService(newJWT, time.Now)
}
