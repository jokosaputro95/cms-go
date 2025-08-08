package config

import (
	"log"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/joho/godotenv"
)

type ServerConfig struct {
	AppName string
	AppVersion string
	AppEnv string

	ServerHost string
	ServerPort string
	ServerReadTimeout time.Duration
	ServerWriteTimeout time.Duration
	ServerIdleTimeout time.Duration
}

type DatabaseConfig struct {
	DBHost string
	DBPort string
	DBUser string
	DBPass string
	DBName string
	DBSSLMode string

	// Connection Pool
	MaxOpenConn int
	MaxIdleConn int
	MaxLifetimeConn time.Duration
	MaxIdleTimeConn time.Duration

	// Timeouts
	QueryTimeout time.Duration
	ConnectTimeout time.Duration
	
}

type JWTConfig struct {
	JWTSecret string
	JWTRefreshSecret string
	JWTExpiresIn time.Duration
	JWTRefreshExpiresIn time.Duration
	JWTIssuer string
}

type EmailConfig struct {
	EmailSMTPHost string
	EmailSMTPPort string
	EmailSMTPUsername string
	EmailSMTPPassword string
}

type Config struct {
	Server ServerConfig
	Database DatabaseConfig
	JWT JWTConfig
	Email EmailConfig
}

var (
	cfg *Config
	once sync.Once
)

func LoadConfig(isProd bool, envPath string) (*Config, error) {
	once.Do(func() {
		// Load .env file
		if err := godotenv.Load(envPath); err != nil {
			log.Fatalf("Error loading .env file: %v", err)
		}

		var (
			dbHost, dbPort, dbUser, dbPass, dbName, dbSSLMode string
		)

		if isProd {
			dbHost = GetEnv("DB_HOST", "localhost")
			dbPort = GetEnv("DB_PORT", "5432")
			dbUser = GetEnv("DB_USER", "developer")
			dbPass = GetEnv("DB_PASS", "password")
			dbName = GetEnv("DB_NAME", "cms_news")
			dbSSLMode = GetEnv("DB_SSLMODE", "disable")
		} else {
			dbHost = GetEnv("DB_HOST_TEST", "localhost")
			dbPort = GetEnv("DB_PORT_TEST", "5432")
			dbUser = GetEnv("DB_USER_TEST", "developer")
			dbPass = GetEnv("DB_PASS_TEST", "password")
			dbName = GetEnv("DB_NAME_TEST", "cms_news_test")
			dbSSLMode = GetEnv("DB_SSLMODE_TEST", "disable")
		}
		
		cfg = &Config{
			Server: ServerConfig{
				AppName: GetEnv("APP_NAME", "CMS GO"),
				AppVersion: GetEnv("APP_VERSION", "1.0.0"),
				AppEnv: GetEnv("APP_ENV", "development"),
				ServerHost: GetEnv("SERVER_HOST", "localhost"),
				ServerPort: GetEnv("SERVER_PORT", "8080"),
				ServerReadTimeout: GetEnvAsDuration("SERVER_READ_TIMEOUT", "10s"),
				ServerWriteTimeout: GetEnvAsDuration("SERVER_WRITE_TIMEOUT", "10s"),
				ServerIdleTimeout: GetEnvAsDuration("SERVER_IDLE_TIMEOUT", "15s"),
			},
			Database: DatabaseConfig{
				DBHost: dbHost,
				DBPort: dbPort,
				DBUser: dbUser,
				DBPass: dbPass,
				DBName: dbName,
				DBSSLMode: dbSSLMode,
				MaxOpenConn: GetEnvAsInt("DB_MAX_OPEN_CONNS", 20),
				MaxIdleConn: GetEnvAsInt("DB_MAX_IDLE_CONNS", 5),
				MaxLifetimeConn: GetEnvAsDuration("DB_CONN_MAX_LIFETIME", "5m"),
				MaxIdleTimeConn: GetEnvAsDuration("DB_CONN_MAX_IDLE_TIME", "30s"),
				QueryTimeout: GetEnvAsDuration("DB_QUERY_TIMEOUT", "10s"),
				ConnectTimeout: GetEnvAsDuration("DB_CONNECT_TIMEOUT", "15s"),
			},
			JWT: JWTConfig{
				JWTSecret: GetEnv("JWT_SECRET", "your-super-secret-jwt-key"),
				JWTRefreshSecret: GetEnv("JWT_REFRESH_SECRET", "your-super-secret-jwt-refresh-key"),
				JWTExpiresIn: GetEnvAsDuration("JWT_ACCESS_TOKEN_TTL", "15m"),
				JWTRefreshExpiresIn: GetEnvAsDuration("JWT_REFRESH_TOKEN_TTL", "168h"),
				JWTIssuer: GetEnv("JWT_ISSUER", "cms-go"),
			},
			Email: EmailConfig{
				EmailSMTPHost: GetEnv("EMAIL_SMTP_HOST", "smtp.gmail.com"),
				EmailSMTPPort: GetEnv("EMAIL_SMTP_PORT", "587"),
				EmailSMTPUsername: GetEnv("EMAIL_SMTP_USERNAME", ""),
				EmailSMTPPassword: GetEnv("EMAIL_SMTP_PASSWORD", ""),
			},
		}
	})
	
	return cfg, nil
}

// Herlper functions
func GetEnv(key string, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	} else {
		log.Printf("Warning: %s environment variable is not set, using default value: %s", key, defaultValue)
	}

	return defaultValue
}

func GetEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		} else {
			log.Printf("Warning: invalid int for %s: %v, using default %d", key, err, defaultValue)
		}
	}
	return defaultValue
}

func GetEnvAsDuration(key string, defaultValue string) time.Duration {
	if value := os.Getenv(key); value != "" {
		if durationValue, err := time.ParseDuration(value); err == nil {
			return durationValue
		} else {
			log.Printf("Warning: invalid duration for %s: %v, using default %s", key, err, defaultValue)
		}
	}

	duration, err := time.ParseDuration(defaultValue)
	if err != nil {
		log.Fatalf("Error parsing default duration: %v", err)
	}

	return duration
}