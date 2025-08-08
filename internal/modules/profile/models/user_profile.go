package models

import (
	"time"
)

// UserProfile represents the user_profiles table (UPDATED)
type UserProfile struct {
    ID           string     `json:"id" db:"id"`
    UserID       string     `json:"user_id" db:"user_id"`
    FirstName    string     `json:"first_name" db:"first_name"`
    LastName     *string    `json:"last_name,omitempty" db:"last_name"`
    Phone        *string    `json:"phone,omitempty" db:"phone"` // Now nullable
    Bio          *string    `json:"bio,omitempty" db:"bio"`
    AvatarURL    *string    `json:"avatar_url,omitempty" db:"avatar_url"`
    NIK          *string    `json:"nik,omitempty" db:"nik"`
    DateOfBirth  *time.Time `json:"date_of_birth,omitempty" db:"date_of_birth"`
    Gender       *string    `json:"gender,omitempty" db:"gender"`
    Address      *string    `json:"address,omitempty" db:"address"`
    Village      *string    `json:"village,omitempty" db:"village"`
    District     *string    `json:"district,omitempty" db:"district"`
    City         *string    `json:"city,omitempty" db:"city"`
    Province     *string    `json:"province,omitempty" db:"province"`
    PostalCode   *string    `json:"postal_code,omitempty" db:"postal_code"`
    Country      string     `json:"country" db:"country"`
    CreatedAt    time.Time  `json:"created_at" db:"created_at"`
    UpdatedAt    time.Time  `json:"updated_at" db:"updated_at"`
}

// ProfileSummary represents basic profile info for API responses
type ProfileSummary struct {
    ID          string  `json:"id"`
    FirstName   string  `json:"first_name"`
    LastName    *string `json:"last_name,omitempty"`
    Phone       *string `json:"phone,omitempty"`
    Bio         *string `json:"bio,omitempty"`
    AvatarURL   *string `json:"avatar_url,omitempty"`
    City        *string `json:"city,omitempty"`
    Province    *string `json:"province,omitempty"`
    Country     string  `json:"country"`
}

// ToSummary converts UserProfile to ProfileSummary
func (p *UserProfile) ToSummary() ProfileSummary {
    return ProfileSummary{
        ID:        p.ID,
        FirstName: p.FirstName,
        LastName:  p.LastName,
        Phone:     p.Phone,
        Bio:       p.Bio,
        AvatarURL: p.AvatarURL,
        City:      p.City,
        Province:  p.Province,
        Country:   p.Country,
    }
}