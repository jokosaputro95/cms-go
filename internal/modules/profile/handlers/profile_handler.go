package handlers

import (
	"net/http"

	"github.com/jokosaputro95/cms-go/internal/modules/profile/services"
	"github.com/jokosaputro95/cms-go/internal/modules/security/middleware"
	"github.com/jokosaputro95/cms-go/internal/pkg/api"
)

// ProfileHandler menangani permintaan HTTP untuk profil pengguna
type ProfileHandler struct {
	profileService services.ProfileServiceInterface
}

// NewProfileHandler membuat instance baru dari ProfileHandler
func NewProfileHandler(profileService services.ProfileServiceInterface) *ProfileHandler {
	return &ProfileHandler{profileService: profileService}
}

// GetProfile menangani permintaan untuk mendapatkan data profil pengguna
func (h *ProfileHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	// Ambil userID dari context
	userID, ok := r.Context().Value(middleware.UserIDContextKey).(string)
	if !ok {
		api.SendError(w, http.StatusInternalServerError, "User ID not found in context")
		return
	}

	profile, err := h.profileService.GetProfile(r.Context(), userID)
	if err != nil {
		api.SendError(w, http.StatusInternalServerError, "Failed to get user profile")
		return
	}

	api.SendSuccess(w, http.StatusOK, "Profile fetched successfully", profile, nil)
}