package models

import (
	"time"
)

// UserProfile merepresentasikan tabel 'user_profiles' di database
type UserProfile struct {
	ID        string     `json:"id"`
	UserID    string     `json:"user_id"`
	FirstName string     `json:"first_name"`
	LastName  *string    `json:"last_name"`
	Phone     *string    `json:"phone"`
	Bio       *string    `json:"bio"`
	AvatarURL *string    `json:"avatar_url"`
	NIK       *string    `json:"nik"`
	DateOfBirth *time.Time `json:"date_of_birth"`
	Gender    *string    `json:"gender"`
	Address   *string    `json:"address"`
	Village   *string    `json:"village"`
	District  *string    `json:"district"`
	City      *string    `json:"city"`
	Province  *string    `json:"province"`
	PostalCode *string    `json:"postal_code"`
	Country   string     `json:"country"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}