package models

import "time"

// RevokedToken merepresentasikan token yang telah dicabut di database
type RevokedToken struct {
	ID        string    `json:"id"`
	Token     string    `json:"token"`
	RevokedAt time.Time `json:"revoked_at"`
	ExpiresAt time.Time `json:"expires_at"`
}