package models

import (
	"time"
)

// EmailVerificationToken merepresentasikan tabel 'email_verification_tokens' di database
type EmailVerificationToken struct {
	ID        string     `json:"id"`
	UserID    string     `json:"user_id"`
	Token     string     `json:"token"`
	Email     string     `json:"email"`
	TokenType string     `json:"token_type"`
	ExpiresAt time.Time  `json:"expires_at"`
	UsedAt    *time.Time `json:"used_at"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}