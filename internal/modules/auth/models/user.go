package models

import (
	"time"
)

// User merepresentasikan tabel 'users' di database
type User struct {
	ID                 string     `json:"id"`
	Username           string     `json:"username"`
	Email              string     `json:"email"`
	PasswordHash       *string    `json:"-"` // Tidak akan disertakan dalam JSON
	RegistrationMethod string     `json:"registration_method"`
	OAuthProvider      *string    `json:"oauth_provider"`
	Status             string     `json:"status"`
	EmailVerified      bool       `json:"email_verified"`
	EmailVerifiedAt    *time.Time `json:"email_verified_at"`
	LastActionBy       *string    `json:"last_action_by"`
	IssuedReason       *string    `json:"issued_reason"`
	IssuedAt           *time.Time `json:"issued_at"`
	CurrentLoginAt     *time.Time `json:"current_login_at"`
	CurrentLoginIP     *string    `json:"current_login_ip"`
	FailedLoginAttempts int        `json:"failed_login_attempts"`
	LockedUntil        *time.Time `json:"locked_until"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`
}