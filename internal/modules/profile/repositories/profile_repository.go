package repositories

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jokosaputro95/cms-go/internal/modules/profile/models"
)

// ProfileRepositoryInterface mendefinisikan kontrak untuk interaksi database profil
type ProfileRepositoryInterface interface {
	FindProfileByUserID(ctx context.Context, userID string) (*models.UserProfile, error)
}

// ProfileRepository adalah implementasi dari ProfileRepositoryInterface
type ProfileRepository struct {
	db *sql.DB
}

// NewProfileRepository membuat instance baru dari ProfileRepository
func NewProfileRepository(db *sql.DB) *ProfileRepository {
	return &ProfileRepository{db: db}
}

// FindProfileByUserID mencari profil berdasarkan user ID
func (r *ProfileRepository) FindProfileByUserID(ctx context.Context, userID string) (*models.UserProfile, error) {
    query := `
        SELECT 
            id, user_id, first_name, last_name, phone, bio, avatar_url, nik, date_of_birth, gender, address, village, district, city, province, postal_code, country, created_at, updated_at
        FROM user_profiles 
        WHERE user_id = $1
    `
    profile := &models.UserProfile{}
    err := r.db.QueryRowContext(ctx, query, userID).Scan(
        &profile.ID,
        &profile.UserID,
        &profile.FirstName,
        &profile.LastName,
        &profile.Phone,
        &profile.Bio,
        &profile.AvatarURL,
        &profile.NIK,
        &profile.DateOfBirth,
        &profile.Gender,
        &profile.Address,
        &profile.Village,
        &profile.District,
        &profile.City,
        &profile.Province,
        &profile.PostalCode,
        &profile.Country,
        &profile.CreatedAt,
        &profile.UpdatedAt,
    )
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, nil
        }
        return nil, fmt.Errorf("gagal mencari profil pengguna: %w", err)
    }
    return profile, nil
}