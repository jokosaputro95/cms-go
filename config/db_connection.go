package config

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

type Database struct {
	DB *sql.DB
}

func SetUpDatabase(cfg DatabaseConfig) (*Database, error) {
	// Build connection string
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBUser,
		cfg.DBPass,
		cfg.DBName,
		cfg.DBSSLMode,
	)

	// Open Connection
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("error opening database connection: %w", err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(cfg.MaxOpenConn)
	db.SetMaxIdleConns(cfg.MaxIdleConn)
	db.SetConnMaxLifetime(cfg.MaxLifetimeConn)
	db.SetConnMaxIdleTime(cfg.MaxIdleTimeConn)

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), cfg.ConnectTimeout)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return nil, fmt.Errorf("error pinging database: %w", err)
	}

	log.Printf("Database connected successfully to %s:%s/%s", cfg.DBHost, cfg.DBPort, cfg.DBName)

	return &Database{DB: db}, nil
}

// Close closes the database connection
func (db *Database) Close() error {
	return db.DB.Close()
}

// HealthCheck checks the database connection
func (db *Database) HealthCheck(ctx context.Context) error {
    // Test ping
    if err := db.DB.PingContext(ctx); err != nil {
        return fmt.Errorf("database ping failed: %w", err)
    }

    // Test simple query
    var result int
    query := "SELECT 1"
    if err := db.DB.QueryRowContext(ctx, query).Scan(&result); err != nil {
        return fmt.Errorf("database query test failed: %w", err)
    }

    return nil
}

// GetStats returns database connection pool statistics
func (db *Database) GetStats() sql.DBStats {
    return db.DB.Stats()
}