package middleware

import (
	"context"
	"log"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jokosaputro95/cms-go/internal/modules/auth/services"
	"github.com/jokosaputro95/cms-go/internal/pkg/api"
)

// key untuk menyimpan UserID di context
type contextKey string

const UserIDContextKey contextKey = "userID"

// AuthMiddleware adalah middleware untuk memvalidasi JWT
func AuthMiddleware(jwtService services.JWTService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				api.SendError(w, http.StatusUnauthorized, "Authorization header is required")
				return
			}

			// Format header harus "Bearer <token>"
			tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
			if len(tokenStr) == len(authHeader) { // Tidak ada prefix "Bearer "
				api.SendError(w, http.StatusUnauthorized, "Invalid token format")
				return
			}

			// Validasi token
			token, err := jwtService.ValidateAccessToken(tokenStr)
			if err != nil {
				log.Printf("Gagal memvalidasi token: %v", err)
				api.SendError(w, http.StatusUnauthorized, "Invalid or expired token")
				return
			}

			if !token.Valid {
				api.SendError(w, http.StatusUnauthorized, "Invalid token")
				return
			}

			// Ambil klaim dari token dan tambahkan ke context
			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				api.SendError(w, http.StatusUnauthorized, "Invalid token claims")
				return
			}

			userID, ok := claims["user_id"].(string)
			if !ok {
				api.SendError(w, http.StatusUnauthorized, "User ID not found in token")
				return
			}
			
			// Tambahkan UserID ke context permintaan
			ctx := context.WithValue(r.Context(), UserIDContextKey, userID)
			
			// Lanjutkan ke handler berikutnya dengan context yang baru
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}