package services

import (
	"fmt"
	"time"

	"github.com/jokosaputro95/cms-go/config"
	"github.com/jokosaputro95/cms-go/internal/modules/auth/dto"

	"github.com/golang-jwt/jwt/v5"
)

// JWTService mendefinisikan kontrak untuk layanan JWT
type JWTService interface {
	GenerateTokenPair(userID, email string) (*dto.AuthResponseDTO, error)
	ValidateAccessToken(tokenStr string) (*jwt.Token, error)
	ValidateRefreshToken(tokenStr string) (*jwt.Token, error)
}

// jwtCustomClaims menyimpan data custom yang akan dimasukkan ke dalam JWT
type jwtCustomClaims struct {
	UserID    string `json:"user_id"`
	Email     string `json:"email"`
	TokenType string `json:"token_type"`
	jwt.RegisteredClaims
}

// jwtService adalah implementasi dari JWTService
type jwtService struct {
	cfg *config.Config
}

// NewJWTService membuat instance baru dari jwtService
func NewJWTService(cfg *config.Config) JWTService {
	return &jwtService{cfg: cfg}
}

// GenerateTokenPair membuat access token dan refresh token
func (s *jwtService) GenerateTokenPair(userID, email string) (*dto.AuthResponseDTO, error) {
	// Buat access token
	accessClaims := &jwtCustomClaims{
		UserID:    userID,
		Email:     email,
		TokenType: "access",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * time.Duration(s.cfg.JWT.JWTExpiresIn))),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "cms-go",
		},
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	signedAccessToken, err := accessToken.SignedString([]byte(s.cfg.JWT.JWTSecret))
	if err != nil {
		return nil, fmt.Errorf("gagal menandatangani access token: %w", err)
	}

	// Buat refresh token
	refreshClaims := &jwtCustomClaims{
		UserID:    userID,
		Email:     email,
		TokenType: "refresh",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * time.Duration(s.cfg.JWT.JWTRefreshExpiresIn))),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "cms-go",
		},
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	signedRefreshToken, err := refreshToken.SignedString([]byte(s.cfg.JWT.JWTRefreshSecret))
	if err != nil {
		return nil, fmt.Errorf("gagal menandatangani refresh token: %w", err)
	}

	return &dto.AuthResponseDTO{
		AccessToken:  signedAccessToken,
		RefreshToken: signedRefreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    int64(s.cfg.JWT.JWTExpiresIn.Seconds()),
	}, nil
}

// ValidateAccessToken memvalidasi access token menggunakan secret key
func (s *jwtService) ValidateAccessToken(tokenStr string) (*jwt.Token, error) {
	return jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("metode penandatanganan tidak valid: %v", token.Header["alg"])
		}
		return []byte(s.cfg.JWT.JWTSecret), nil
	})
}

// ValidateRefreshToken memvalidasi refresh token menggunakan refresh secret key
func (s *jwtService) ValidateRefreshToken(tokenStr string) (*jwt.Token, error) {
	return jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("metode penandatanganan tidak valid: %v", token.Header["alg"])
		}
		return []byte(s.cfg.JWT.JWTRefreshSecret), nil
	})
}