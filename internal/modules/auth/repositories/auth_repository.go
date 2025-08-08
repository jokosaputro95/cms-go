package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jokosaputro95/cms-go/internal/modules/auth/models"
	profiles "github.com/jokosaputro95/cms-go/internal/modules/profile/models"

	"github.com/google/uuid"
)

// AuthRepositoryInterface mendefinisikan kontrak untuk interaksi database otentikasi
type AuthRepositoryInterface interface {
	SaveUser(ctx context.Context, user *models.User, profile *profiles.UserProfile) error
	FindUserByEmail(ctx context.Context, email string) (*models.User, error)
	FindUserByUsername(ctx context.Context, username string) (*models.User, error)
	SaveVerificationToken(ctx context.Context, token *models.EmailVerificationToken) error
	FindVerificationToken(ctx context.Context, tokenID string) (*models.EmailVerificationToken, error)
	UpdateUserStatus(ctx context.Context, userID string, status string, emailVerified bool, verifiedAt time.Time) error
	UpdateUserLoginStatus(ctx context.Context, userID string, ip string, failedAttempts int, lockUntil *time.Time) error
	RevokeRefreshToken(ctx context.Context, token string, expiresAt time.Time) error
	IsTokenRevoked(ctx context.Context, token string) (bool, error)
}

// AuthRepository adalah implementasi dari AuthRepositoryInterface
type AuthRepository struct {
	db *sql.DB
}

// NewAuthRepository membuat instance baru dari AuthRepository
func NewAuthRepository(db *sql.DB) *AuthRepository {
	return &AuthRepository{db: db}
}

// SaveUser menyimpan pengguna baru dan profilnya dalam satu transaksi
func (r *AuthRepository) SaveUser(ctx context.Context, user *models.User, profile *profiles.UserProfile) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("gagal memulai transaksi: %w", err)
	}
	defer tx.Rollback() // Rollback jika ada error

	// Simpan User
	user.ID = uuid.New().String()
	user.CreatedAt = time.Now().UTC()
	user.UpdatedAt = user.CreatedAt
	user.Status = "pending"
	user.RegistrationMethod = "manual"

	userQuery := `
		INSERT INTO users (
			id, username, email, password_hash, registration_method, status, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	_, err = tx.ExecContext(ctx, userQuery,
		user.ID,
		user.Username,
		user.Email,
		user.PasswordHash,
		user.RegistrationMethod,
		user.Status,
		user.CreatedAt,
		user.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("gagal menyimpan pengguna: %w", err)
	}

	// Simpan User Profile
	profile.ID = uuid.New().String()
	profile.UserID = user.ID
	profile.CreatedAt = time.Now().UTC()
	profile.UpdatedAt = profile.CreatedAt

	profileQuery := `
		INSERT INTO user_profiles (
			id, user_id, first_name, last_name, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err = tx.ExecContext(ctx, profileQuery,
		profile.ID,
		profile.UserID,
		profile.FirstName,
		profile.LastName,
		profile.CreatedAt,
		profile.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("gagal menyimpan profil pengguna: %w", err)
	}

	return tx.Commit()
}

// FindUserByEmail mencari pengguna berdasarkan email
func (r *AuthRepository) FindUserByEmail(ctx context.Context, email string) (*models.User, error) {
	query := `SELECT id, username, email, password_hash, registration_method, status, email_verified, email_verified_at, failed_login_attempts, locked_until, current_login_ip, created_at, updated_at FROM users WHERE email = $1`
	user := &models.User{}
	err := r.db.QueryRowContext(ctx, query, email).Scan(
        &user.ID,
		&user.Username,
        &user.Email,
        &user.PasswordHash,
		&user.RegistrationMethod,
        &user.Status,
		&user.EmailVerified,
		&user.EmailVerifiedAt,
        &user.FailedLoginAttempts,
        &user.LockedUntil,
        &user.CurrentLoginIP,
        &user.CreatedAt,
        &user.UpdatedAt,
    )
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Mengembalikan nil tanpa error jika tidak ada data
		}
		return nil, fmt.Errorf("gagal mencari pengguna berdasarkan email: %w", err)
	}
	return user, nil
}

// FindUserByUsername mencari pengguna berdasarkan username
func (r *AuthRepository) FindUserByUsername(ctx context.Context, username string) (*models.User, error) {
	query := `SELECT id, username, email, password_hash, registration_method, status, email_verified, email_verified_at, failed_login_attempts, locked_until, current_login_ip, created_at, updated_at FROM users WHERE username = $1`
	user := &models.User{}
	err := r.db.QueryRowContext(ctx, query, username).Scan(
		&user.ID,
		&user.Username,
        &user.Email,
        &user.PasswordHash,
		&user.RegistrationMethod,
        &user.Status,
		&user.EmailVerified,
		&user.EmailVerifiedAt,
        &user.FailedLoginAttempts,
        &user.LockedUntil,
        &user.CurrentLoginIP,
        &user.CreatedAt,
        &user.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("gagal mencari pengguna berdasarkan username: %w", err)
	}
	return user, nil
}

// SaveVerificationToken menyimpan token verifikasi email baru
func (r *AuthRepository) SaveVerificationToken(ctx context.Context, token *models.EmailVerificationToken) error {
	query := `
		INSERT INTO email_verification_tokens 
		(id, user_id, email, token, token_type, expires_at) 
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := r.db.ExecContext(ctx, query,
		token.ID,
		token.UserID, 
		token.Email, 
		token.Token, 
		token.TokenType, 
		token.ExpiresAt,
	)
	if err != nil {
		return fmt.Errorf("gagal menyimpan token verifikasi: %w", err)
	}
	return nil
}

// FindVerificationToken mencari token verifikasi berdasarkan tokenID
func (r *AuthRepository) FindVerificationToken(ctx context.Context, tokenStr string) (*models.EmailVerificationToken, error) {
	query := `
		SELECT user_id, email, token, token_type, expires_at, used_at
		FROM email_verification_tokens
		WHERE token = $1
		LIMIT 1
	`
	token := &models.EmailVerificationToken{}
	err := r.db.QueryRowContext(ctx, query, tokenStr).Scan(
		&token.UserID,
		&token.Email,
		&token.Token,
		&token.TokenType,
		&token.ExpiresAt,
		&token.UsedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Token tidak ditemukan
		}
		return nil, fmt.Errorf("gagal mencari token verifikasi: %w", err)
	}
	return token, nil
}

// UpdateUserStatus memperbarui status pengguna setelah verifikasi email
func (r *AuthRepository) UpdateUserStatus(ctx context.Context, userID, tokenStr string) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("gagal memulai transaksi: %w", err)
	}
	defer tx.Rollback() // Rollback jika ada error setelahnya

	// Query 1: Mengaktifkan pengguna dan memperbarui status verifikasi email
	userQuery := `
		UPDATE users 
		SET status = 'active', email_verified = true, email_verified_at = NOW() 
		WHERE id = $1
	`
	res, err := tx.ExecContext(ctx, userQuery, userID)
	if err != nil {
		return fmt.Errorf("gagal mengaktifkan pengguna: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("gagal mendapatkan jumlah baris yang terpengaruh: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("pengguna tidak ditemukan dengan ID: %s", userID)
	}

	// Query 2: Menandai token sebagai sudah digunakan
	tokenQuery := `
		UPDATE email_verification_tokens
		SET used_at = NOW() 
		WHERE token = $1
	`
	_, err = tx.ExecContext(ctx, tokenQuery, tokenStr)
	if err != nil {
		return fmt.Errorf("gagal menandai token sebagai sudah digunakan: %w", err)
	}

	// Commit transaksi
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("gagal commit transaksi: %w", err)
	}

	return nil
}

// UpdateUserLoginStatus memperbarui status login pengguna
func (r *AuthRepository) UpdateUserLoginStatus(ctx context.Context, userID string, ip string, failedAttempts int, lockUntil *time.Time) error {
	query := `
		UPDATE users
		SET 
			current_login_ip = $1,
			failed_login_attempts = $2,
			locked_until = $3
		WHERE id = $4
	`
	_, err := r.db.ExecContext(ctx, query,
		ip,
		failedAttempts,
		lockUntil,
		userID,
	)
	if err != nil {
		return fmt.Errorf("gagal memperbarui status login pengguna: %w", err)
	}
	return nil
}

// RevokeRefreshToken mencatat token yang dicabut ke dalam database
func (r *AuthRepository) RevokeRefreshToken(ctx context.Context, token string, expiresAt time.Time) error {
	query := `
		INSERT INTO revoked_tokens 
		(token, expires_at) 
		VALUES ($1, $2)
	`
	_, err := r.db.ExecContext(ctx, query, token, expiresAt)
	if err != nil {
		return fmt.Errorf("gagal mencabut refresh token: %w", err)
	}
	return nil
}

// IsTokenRevoked memeriksa apakah token sudah dicabut
func (r *AuthRepository) IsTokenRevoked(ctx context.Context, token string) (bool, error) {
	query := `
		SELECT EXISTS (
			SELECT 1 FROM revoked_tokens 
			WHERE token = $1
		)
	`
	var isRevoked bool
	err := r.db.QueryRowContext(ctx, query, token).Scan(&isRevoked)
	if err != nil {
		return false, fmt.Errorf("gagal memeriksa revoked token: %w", err)
	}
	return isRevoked, nil
}