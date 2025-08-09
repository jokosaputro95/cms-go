package app

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/jokosaputro95/cms-go/config"
	auth_hendlers "github.com/jokosaputro95/cms-go/internal/modules/auth/handlers"
	auth_repositories "github.com/jokosaputro95/cms-go/internal/modules/auth/repositories"
	auth_routes "github.com/jokosaputro95/cms-go/internal/modules/auth/routes"
	auth_services "github.com/jokosaputro95/cms-go/internal/modules/auth/services"
	profile_handlers "github.com/jokosaputro95/cms-go/internal/modules/profile/handlers"
	profile_repositories "github.com/jokosaputro95/cms-go/internal/modules/profile/repositories"
	profile_services "github.com/jokosaputro95/cms-go/internal/modules/profile/services"
	"github.com/jokosaputro95/cms-go/internal/modules/security/middleware"
	"github.com/jokosaputro95/cms-go/internal/pkg/email"
)

// App adalah struktur utama yang menampung server dan dependensi
type App struct {
	Config      *config.Config
	Server      *http.Server
	DB          *config.Database
	AuthRoutes  *auth_routes.AuthRoutes
	AuthMiddleware func(http.Handler) http.Handler
	ProfileHandler *profile_handlers.ProfileHandler
	
}

// StartServer adalah fungsi entry point untuk inisialisasi aplikasi
func StartServer(isProd bool, envFile string) (*App, error) {
	// 1. Muat konfigurasi
	log.Println("Loading server configuration...")
	cfg, err := config.LoadConfig(isProd, envFile)

	if err != nil {
        return nil, fmt.Errorf("failed to load configuration: %w", err)
    }

	// 2. Buat koneksi database
	db, err := config.SetUpDatabase(cfg.Database)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	
	// 3. Inisialisasi container dengan semua dependensi
	jwtService := auth_services.NewJWTService(cfg)
	authRepo := auth_repositories.NewAuthRepository(db.DB)
	emailSvc := email.NewEmailService(cfg)
	authService := auth_services.NewAuthService(authRepo, jwtService, emailSvc)
	authHandler := auth_hendlers.NewAuthHandler(authService)

	// Inisialisasi service dan repository untuk profile
	profileRepo := profile_repositories.NewProfileRepository(db.DB)
	profileService := profile_services.NewProfileService(profileRepo)
	profileHandler := profile_handlers.NewProfileHandler(profileService)
	
	// Inisialisasi rute dan middleware
	authRoutes := auth_routes.NewAuthRoutes(authHandler)
	authMiddleware := middleware.AuthMiddleware(jwtService, authService)

	// 4. Daftarkan rute ke router
	router := http.NewServeMux()
	authRoutes.RegisterRoutes(router)

	// Contoh pendaftaran rute yang dilindungi
	protectedRouter := http.NewServeMux()
	protectedRouter.HandleFunc("/profile", profileHandler.GetProfile)
	router.Handle("/profile", authMiddleware(protectedRouter))

	// 5. Buat instance server
	server := &http.Server{
		Addr:    ":" + cfg.Server.ServerPort,
		Handler: router,
		ReadTimeout:  cfg.Server.ServerReadTimeout,
        WriteTimeout: cfg.Server.ServerWriteTimeout,
        IdleTimeout:  cfg.Server.ServerIdleTimeout,
	}

	return &App{
		Config: cfg,
		Server: server,
		DB: db,
		AuthRoutes: authRoutes,
		AuthMiddleware: authMiddleware,
	}, nil
}

// Start memulai server HTTP
func (a *App) Start() error {
	log.Println("Starting server...")
	address := fmt.Sprintf("%s:%s", a.Config.Server.ServerHost, a.Config.Server.ServerPort)
	
	log.Printf("INFO: %-16s: %s", "APP_NAME", a.Config.Server.AppName)
	log.Printf("INFO: %-16s: %s", "APP_VERSION", a.Config.Server.AppVersion)
	log.Printf("INFO: %-16s: %s", "APP_ENV", a.Config.Server.AppEnv)
	log.Printf("🚀 Server running on http://%s", address)

	if err := a.Server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}

// Shutdown menutup server secara bertahap dan melepaskan sumber daya
func (a *App) Shutdown(ctx context.Context) error {
	log.Println("Shutting down server...")
	// Tutup koneksi database
	if err := a.DB.Close(); err != nil {
		log.Printf("Error closing database: %v", err)
	} else {
		log.Println("✅ Database connection closed")
	}

	return a.Server.Shutdown(ctx)
}