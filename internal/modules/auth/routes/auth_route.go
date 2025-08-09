package routes

import (
	"net/http"

	"github.com/jokosaputro95/cms-go/internal/modules/auth/handlers"
)

// AuthRoutes mengelola pendaftaran rute untuk modul otentikasi
type AuthRoutes struct {
	authHandler *handlers.AuthHandler
}

// NewAuthRoutes membuat instance baru dari AuthRoutes
func NewAuthRoutes(authHandler *handlers.AuthHandler) *AuthRoutes {
	return &AuthRoutes{authHandler: authHandler}
}

// RegisterRoutes mendaftarkan rute-rute otentikasi ke router yang diberikan
func (r *AuthRoutes) RegisterRoutes(router *http.ServeMux) {
	router.HandleFunc("/auth/register", r.authHandler.Register)
	router.HandleFunc("/auth/login", r.authHandler.Login)
	router.HandleFunc("/auth/logout", r.authHandler.Logout)
	router.HandleFunc("/auth/verify-email", r.authHandler.VerifyEmail)
	router.HandleFunc("/auth/refresh-token", r.authHandler.RefreshToken)
}