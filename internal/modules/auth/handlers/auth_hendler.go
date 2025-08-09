package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/jokosaputro95/cms-go/internal/modules/auth/dto"
	"github.com/jokosaputro95/cms-go/internal/modules/auth/services"
	"github.com/jokosaputro95/cms-go/internal/pkg/api"
)

// AuthHandler menangani permintaan HTTP untuk otentikasi
type AuthHandler struct {
	authService services.AuthServiceInterface
}

// NewAuthHandler membuat instance baru dari AuthHandler
func NewAuthHandler(authService services.AuthServiceInterface) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// Register menangani permintaan registrasi pengguna
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		api.SendError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var req dto.RegisterRequestDTO
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		api.SendError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	err = h.authService.RegisterUser(r.Context(), &req)
	if err != nil {
		log.Printf("Gagal registrasi pengguna: %v", err)
		api.SendError(w, http.StatusInternalServerError, "Failed to register user")
		return
	}

	api.SendSuccess(w, http.StatusCreated, "User registered successfully. Please check your email for verification.", nil, nil)
}

// Login menangani permintaan login pengguna
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		api.SendError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var req dto.LoginRequestDTO
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		api.SendError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	ip := r.RemoteAddr

	tokenPair, err := h.authService.LoginUser(r.Context(), &req, ip)
    if err != nil {
        log.Printf("Gagal login pengguna: %v", err)
        
        // Error types
        // switch {
		// case errors.Is(err, services.AuthServiceError(services.ErrUserLocked.Error())):
        //     api.SendError(w, http.StatusTooManyRequests, "Account is temporarily locked")
        // case errors.Is(err, services.AuthServiceError(services.ErrInvalidCredentials.Error())):
        //     api.SendError(w, http.StatusUnauthorized, "Invalid credentials")
        // default:
        //     api.SendError(w, http.StatusInternalServerError, "Login failed")
        // }
        // return

		switch err := err.(type) {
		case *services.LockoutError:
			retryAfter := int(time.Until(err.UnlockAt).Seconds())
			details := map[string]interface{}{
				"unlock_at": err.UnlockAt,
				"retry_after": retryAfter,
			}
			w.Header().Set("Retry-After", fmt.Sprintf("%d", retryAfter))
			api.SendDetailedError(
				w, http.StatusTooManyRequests,
				err.Error(),
				"account_locked",
				details,
			)
		case services.AuthServiceError:
			api.SendDetailedError(w, http.StatusUnauthorized, err.Error(), "invalid_credentials", nil)

		default:
			api.SendError(w, http.StatusInternalServerError, "Login failed")
		}
		return
    }
    
    api.SendSuccess(w, http.StatusOK, "Login successful", tokenPair, nil)
}

// VerifyEmail menangani verifikasi email dari tautan
func (h *AuthHandler) VerifyEmail(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	if token == "" {
		api.SendError(w, http.StatusBadRequest, "Token is missing")
		return
	}

	err := h.authService.VerifyEmail(r.Context(), token)
	if err != nil {
		log.Printf("Gagal verifikasi email: %v", err)
		api.SendError(w, http.StatusInternalServerError, "Failed to verify email")
		return
	}

	api.SendSuccess(w, http.StatusOK, "Email verified successfully, you can now login", nil, nil)
}

// RefreshToken menangani refresh token untuk mendapatkan access token baru
func (h *AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		api.SendError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var req dto.RefreshTokenRequestDTO
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		api.SendError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	tokenPair, err := h.authService.RefreshToken(r.Context(), req.RefreshToken)
	if err != nil {
		log.Printf("Gagal refresh token: %v", err)
		api.SendError(w, http.StatusUnauthorized, "Invalid refresh token")
		return
	}

	api.SendSuccess(w, http.StatusOK, "Token refreshed successfully", tokenPair, nil)
}

// Logout menangani permintaan logout pengguna
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		api.SendError(w, http.StatusUnauthorized, "Authorization header is missing")
		return
	}

	accessToken := strings.Replace(authHeader, "Bearer ", "", 1)
	if accessToken == "" {
		api.SendError(w, http.StatusUnauthorized, "Access token is missing")
		return
	}

	err := h.authService.LogoutUser(r.Context(), accessToken)
	if err != nil {
		log.Printf("Gagal logout pengguna: %v", err)
		api.SendError(w, http.StatusInternalServerError, "Failed to logout user")
		return
	}

	api.SendSuccess(w, http.StatusOK, "Logout successful", nil, nil)
}