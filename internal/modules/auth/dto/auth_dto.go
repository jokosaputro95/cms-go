package dto

import "time"

// RegisterRequestDTO digunakan untuk menerima data registrasi dari pengguna
type RegisterRequestDTO struct {
	Username  string `json:"username" validate:"required,min=3,max=50"`
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=8"`
}

// LoginRequestDTO digunakan untuk menerima data login dari pengguna
type LoginRequestDTO struct {
	Identifier string `json:"identifier" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// AuthResponseDTO digunakan untuk mengirimkan token ke klien setelah login atau refresh
type AuthResponseDTO struct {
	ID string `json:"id"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int64  `json:"expires_in"` // Waktu kadaluarsa
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type RefreshTokenRequestDTO struct {
	RefreshToken string `json:"refresh_token"`
}

type VerifyEmailRequestDTO struct {
	Token string `json:"token"`
}