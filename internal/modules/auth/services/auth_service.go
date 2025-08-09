package services

import (
	"context"
	"errors"
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jokosaputro95/cms-go/internal/modules/auth/dto"
	"github.com/jokosaputro95/cms-go/internal/modules/auth/models"
	"github.com/jokosaputro95/cms-go/internal/modules/auth/repositories"
	profiles "github.com/jokosaputro95/cms-go/internal/modules/profile/models"
	"github.com/jokosaputro95/cms-go/internal/pkg/email"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// AuthServiceError adalah tipe error kustom
type AuthServiceError string

func (e AuthServiceError) Error() string {
	return string(e)
}

const (
	ErrUserAlreadyExists = AuthServiceError("email atau username sudah terdaftar")
	ErrInvalidToken      = AuthServiceError("token tidak valid atau kedaluwarsa")
	ErrTokenAlreadyUsed  = AuthServiceError("token sudah digunakan sebelumnya")
	ErrInvalidCredentials = AuthServiceError("kredensial tidak valid")
	ErrUserLocked        = AuthServiceError("Account is temporarily locked, please try again later")
	maxFailedAttempts = 5
	lockoutDuration   = 30 * time.Minute
)

type LockoutError struct {
	Message string
	UnlockAt time.Time
}

func (e LockoutError) Error() string {
	return e.Message

}

// AuthServiceInterface mendefinisikan kontrak untuk service otentikasi
type AuthServiceInterface interface {
	RegisterUser(ctx context.Context, req *dto.RegisterRequestDTO) error
    LoginUser(ctx context.Context, req *dto.LoginRequestDTO, ip string) (*dto.AuthResponseDTO, error)
    LogoutUser(ctx context.Context, tokenStr string) error
    RefreshToken(ctx context.Context, refreshTokenStr string) (*dto.AuthResponseDTO, error)
    VerifyEmail(ctx context.Context, token string) error
    IsTokenRevoked(ctx context.Context, token string) (bool, error)
}

// AuthService adalah implementasi dari AuthServiceInterface
type AuthService struct {
	authRepo *repositories.AuthRepository
	jwtSvc JWTService
	emailSvc email.EmailService
	validate *validator.Validate
}

// NewAuthService membuat instance baru dari AuthService
func NewAuthService(authRepo *repositories.AuthRepository, jwtSvc JWTService, emailSvc email.EmailService) *AuthService {
	return &AuthService{
		authRepo: authRepo,
		jwtSvc: jwtSvc,
		emailSvc: emailSvc,
		validate: validator.New(),
	}
}

// RegisterUser memproses logika registrasi pengguna baru
func (s *AuthService) RegisterUser(ctx context.Context, req *dto.RegisterRequestDTO) error {
	// 1. Validasi input menggunakan DTO
	err := s.validate.Struct(req)
	if err != nil {
		return fmt.Errorf("validasi input gagal: %w", err)
	}

	// 2. Cek apakah email atau username sudah ada
	existingUser, err := s.authRepo.FindUserByEmail(ctx, req.Email)
	if err != nil {
		return err
	}
	if existingUser != nil {
		return ErrUserAlreadyExists
	}
	existingUser, err = s.authRepo.FindUserByUsername(ctx, req.Username)
	if err != nil {
		return err
	}
	if existingUser != nil {
		return ErrUserAlreadyExists
	}

	// 3. Hashing password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("gagal melakukan hashing password: %w", err)
	}
	hashedPasswordStr := string(hashedPassword)

	// 4. Siapkan model user dan profile
	user := &models.User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: &hashedPasswordStr,
	}
	
	var emptyString = ""

	profile := &profiles.UserProfile{
		FirstName: "",
		LastName:  &emptyString,
	}

	// 5. Simpan user dan profile ke database
	err = s.authRepo.SaveUser(ctx, user, profile)
	if err != nil {
		return err
	}

	// 6. Buat dan simpan token verifikasi email, lalu kirim email
	verificationToken := &models.EmailVerificationToken{
		ID: uuid.New().String(),
		UserID: user.ID,
		Email:  user.Email,
		Token:  uuid.New().String(),
		TokenType: "email_verification",
		ExpiresAt: time.Now().Add(time.Minute * 30),
	}
	
	// Panggil metode repository untuk menyimpan token
	err = s.authRepo.SaveVerificationToken(ctx, verificationToken)
	if err != nil {
		return err
	}	

	// 7. Kirim email verifikasi di goroutine
	go func(to, token, username string) {
		err := s.emailSvc.SendVerificationEmail(user.Email, verificationToken.Token, user.Username)
		if err != nil {
			log.Printf("Gagal mengirim email verifikasi ke %s: %v", to, err)
		} else {
			log.Printf("Email verifikasi berhasil dikirim ke %s", to)
		}
	}(user.Email, verificationToken.Token, user.Username)

	return nil
}

var emailRegex = regexp.MustCompile(`^[^\s]+$`)

// LoginUser memproses logika login pengguna
func (s *AuthService) LoginUser(ctx context.Context, req *dto.LoginRequestDTO, ip string) (*dto.AuthResponseDTO, error) {
	// 1. Validasi input
	if err := s.validate.Struct(req); err != nil {
		return nil, fmt.Errorf("validasi input gagal: %w", err)
	}

	// 2. Cari user berdasarkan Identifier
	var user *models.User
	var err error
	if emailRegex.MatchString(req.Identifier) {
		user, err = s.authRepo.FindUserByEmail(ctx, req.Identifier)
	} else {
		user, err = s.authRepo.FindUserByUsername(ctx, req.Identifier)
	}

	// 🔍 DEBUG: Log nilai dari database
    if user != nil {
        log.Printf("DEBUG - User found: ID=%s, FailedAttempts=%d, LockedUntil=%v", 
            user.ID, user.FailedLoginAttempts, user.LockedUntil)
    }

	if err != nil {
		return nil, err
	}
	if user == nil || user.PasswordHash == nil {
		return nil, ErrInvalidCredentials
	}

	// Cek apakah akun terkunci
	if user.LockedUntil != nil && user.LockedUntil.After(time.Now().UTC()) {
    	return nil, &LockoutError{
			Message: string(ErrUserLocked),
			UnlockAt: *user.LockedUntil,
		}
	}

	// 3. Bandingkan/Cek password
	if err := bcrypt.CompareHashAndPassword([]byte(*user.PasswordHash), []byte(req.Password)); err != nil {
    // DEBUG: Log sebelum increment
	log.Printf("DEBUG - Password salah. Current FailedAttempts: %d", user.FailedLoginAttempts)
	
	// Hanya increment sekali
    newFailedAttempts := user.FailedLoginAttempts + 1

	// // 🔍 DEBUG: Log nilai yang akan di-update
    //     log.Printf("DEBUG - Will update to: newFailedAttempts=%d", newFailedAttempts)
    
    // Cek apakah perlu lock
    var lockUntil *time.Time
    
	if newFailedAttempts >= maxFailedAttempts {
        lockedTime := time.Now().Add(lockoutDuration)
        lockUntil = &lockedTime
        log.Printf("SECURITY: Account %s locked from IP %s after %d failed attempts", user.ID, ip, newFailedAttempts)
		log.Printf("SECURITY: Account locked until: %v", lockedTime)
    }

	if newFailedAttempts == 3 {
		log.Printf("WARNING: User %s has %d failed attempts from IP %s", user.ID, newFailedAttempts, ip)
	}
    
    // Update ke database
    if errUpd := s.authRepo.UpdateUserLoginStatus(ctx, user.ID, ip, newFailedAttempts, lockUntil); errUpd != nil {
        log.Printf("ERROR: UpdateUserLoginStatus gagal: %v", errUpd)
        return nil, fmt.Errorf("failed to update login status: %w", errUpd)
    }

	// // 🔍 DEBUG: Konfirmasi update berhasil
    //     log.Printf("DEBUG - UpdateUserLoginStatus SUCCESS for user %s with attempts=%d", user.ID, newFailedAttempts)
    
    // Return error yang sesuai
    if lockUntil != nil {
        return nil, &LockoutError{
            Message: string(ErrUserLocked),
            UnlockAt: *lockUntil,
        }
    }
    	return nil, ErrInvalidCredentials
	}
	// Login berhasil: reset failed attempts
	if err := s.authRepo.UpdateUserLoginStatus(ctx, user.ID, ip, 0, nil); err != nil {
		return nil, err
	}

	// 4. Periksa status pengguna
	if user.Status != "active" {
		return nil, errors.New("akun tidak aktif, silakan verifikasi email")
	}

	// 5. Buat access token dan refresh token
	tokenPair, err := s.jwtSvc.GenerateTokenPair(user.ID, user.Email)
	if err != nil {
		return nil, fmt.Errorf("gagal membuat token: %w", err)
	}

	data := &dto.AuthResponseDTO{
		ID:           user.ID,
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		TokenType:    tokenPair.TokenType,
		ExpiresIn:    tokenPair.ExpiresIn,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
	}

	log.Printf("Pengguna %s:%s berhasil login.", user.Username, user.Email)
	return data, nil
}

// VerifyEmail memproses logika verifikasi email
func (s *AuthService) VerifyEmail(ctx context.Context, tokenStr string) error {
	// 1. Cari token di database
	token, err := s.authRepo.FindVerificationToken(ctx, tokenStr)
	if err != nil {
		return err
	}
	if token == nil {
		return ErrInvalidToken // Token tidak ditemukan
	}

	// 2. Cek apakah token sudah kedaluwarsa
	if token.ExpiresAt.Before(time.Now()) {
		return ErrInvalidToken
	}
	
	// 3. Cek apakah token sudah digunakan
	if token.UsedAt != nil {
		return ErrTokenAlreadyUsed
	}

	// 4. Update status pengguna di database
	// Kita akan mengaktifkan pengguna dan menandai token sebagai sudah digunakan dalam satu operasi
	err = s.authRepo.UpdateUserStatus(ctx, token.UserID, token.Token)
	if err != nil {
		return err
	}
	
	// Kirim email selamat datang di goroutine
	// Kita memanggil layanan email yang baru kita buat
	go func(to, username string) {
		err := s.emailSvc.SendWelcomeEmail(to, username)
		if err != nil {
			log.Printf("Gagal mengirim email selamat datang ke %s: %v", to, err)
		} else {
			log.Printf("Email selamat datang berhasil dikirim ke %s", to)
		}
	}(token.Email, token.Email) // Menggunakan email sebagai username sementara

	return nil
}

// RefreshToken memproses permintaan untuk mendapatkan access token baru
func (s *AuthService) RefreshToken(ctx context.Context, refreshTokenStr string) (*dto.AuthResponseDTO, error) {
	// 1. Validasi refresh token secara sintaksis dan cek kedaluwarsa
	token, err := s.jwtSvc.ValidateRefreshToken(refreshTokenStr)
	if err != nil {
		return nil, ErrInvalidToken
	}

	// 2. Periksa apakah refresh token sudah dicabut (di-blacklist)
	isRevoked, err := s.authRepo.IsTokenRevoked(ctx, refreshTokenStr)
	if err != nil {
		return nil, err
	}
	if isRevoked {
		return nil, ErrInvalidToken // Tolak jika token sudah ada di daftar hitam
	}

	// 3. Ambil klaim dari token
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}

	// 4. Periksa jenis token
	tokenType, ok := claims["token_type"].(string)
	if !ok || tokenType != "refresh" {
		return nil, ErrInvalidToken
	}
	
	// 5. Ambil data user
	userID, ok := claims["user_id"].(string)
	if !ok {
		return nil, ErrInvalidToken
	}
	email, ok := claims["email"].(string)
	if !ok {
		return nil, ErrInvalidToken
	}
	
	// 6. Dapatkan waktu kedaluwarsa token lama dari klaim
	expiresAt, err := claims.GetExpirationTime()
	if err != nil {
		return nil, ErrInvalidToken
	}

	// 7. Revoke (cabut) refresh token lama
	err = s.authRepo.RevokeToken(ctx, refreshTokenStr, expiresAt.Time)
	if err != nil {
		return nil, err
	}

	// 8. Buat pasangan token baru
	tokenPair, err := s.jwtSvc.GenerateTokenPair(userID, email)
	if err != nil {
		return nil, fmt.Errorf("gagal membuat token baru: %w", err)
	}


	user, err := s.authRepo.FindUserByEmail(ctx, email)
	if err != nil || user == nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	responseDTO := &dto.AuthResponseDTO{
		ID: userID,
		AccessToken: tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		TokenType: tokenPair.TokenType,
		ExpiresIn: tokenPair.ExpiresIn,
		// CreatedAt: time.Now().UTC(),
		// UpdatedAt: time.Now().UTC(),
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	return responseDTO, nil
}

// LogoutUser mencabut access token yang sedang digunakan
func (s *AuthService) LogoutUser(ctx context.Context, tokenStr string) error {
	// 1. Validasi token dan cek kadaluarsa
	token, err := s.jwtSvc.ValidateAccessToken(tokenStr)
	if err != nil {
		return ErrInvalidToken
	}

	// 2. Periksa jenis token
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return ErrInvalidToken
	}

	tokenType, ok := claims["token_type"].(string)
	if !ok || tokenType != "access" {
		return ErrInvalidToken
	
	}

	expiresAt, err := claims.GetExpirationTime()
	if err != nil {
		return ErrInvalidToken
	}

	// 3. Cabut token dan store ke database
	err = s.authRepo.RevokeToken(ctx, tokenStr, expiresAt.Time)
	if err != nil {
		return fmt.Errorf("gagal cabut token: %w", err)
	}

	return nil
}

// IsTokenRevoked adalah implementasi untuk service yang akan dipanggil oleh middleware
func (s *AuthService) IsTokenRevoked(ctx context.Context, token string) (bool, error) {
    return s.authRepo.IsTokenRevoked(ctx, token)
}