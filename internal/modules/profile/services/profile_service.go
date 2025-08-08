package services

import (
	"context"
	"fmt"

	"github.com/jokosaputro95/cms-go/internal/modules/profile/models"
	"github.com/jokosaputro95/cms-go/internal/modules/profile/repositories"
)

// ProfileServiceInterface mendefinisikan kontrak untuk service profil
type ProfileServiceInterface interface {
	GetProfile(ctx context.Context, userID string) (*models.UserProfile, error)
}

// ProfileService adalah implementasi dari ProfileServiceInterface
type ProfileService struct {
	profileRepo repositories.ProfileRepositoryInterface
}

// NewProfileService membuat instance baru dari ProfileService
func NewProfileService(profileRepo repositories.ProfileRepositoryInterface) *ProfileService {
	return &ProfileService{profileRepo: profileRepo}
}

// GetProfile mengambil data profil pengguna dari database
func (s *ProfileService) GetProfile(ctx context.Context, userID string) (*models.UserProfile, error) {
	profile, err := s.profileRepo.FindProfileByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if profile == nil {
		return nil, fmt.Errorf("profil pengguna tidak ditemukan")
	}
	return profile, nil
}