package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

// Config captures backend runtime settings sourced from environment variables.
type Config struct {
	DatabaseURL     string
	JWTSecret       string
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
	RunMigrations   bool
	MigrationsPath  string
	DBMaxOpenConns  int
	DBMaxIdleConns  int
	DBConnMaxIdle   time.Duration
	DBConnMaxLife   time.Duration
}

func Load() (Config, error) {
	accessTTL, err := parseDuration("APP_ACCESS_TOKEN_TTL", "15m")
	if err != nil {
		return Config{}, fmt.Errorf("parse APP_ACCESS_TOKEN_TTL: %w", err)
	}

	refreshTTL, err := parseDuration("APP_REFRESH_TOKEN_TTL", "168h")
	if err != nil {
		return Config{}, fmt.Errorf("parse APP_REFRESH_TOKEN_TTL: %w", err)
	}

	maxIdleTime, err := parseDuration("APP_DB_CONN_MAX_IDLE_TIME", "5m")
	if err != nil {
		return Config{}, fmt.Errorf("parse APP_DB_CONN_MAX_IDLE_TIME: %w", err)
	}

	maxLifetime, err := parseDuration("APP_DB_CONN_MAX_LIFETIME", "30m")
	if err != nil {
		return Config{}, fmt.Errorf("parse APP_DB_CONN_MAX_LIFETIME: %w", err)
	}

	cfg := Config{
		DatabaseURL:     os.Getenv("APP_DB_URL"),
		JWTSecret:       os.Getenv("APP_JWT_SECRET"),
		AccessTokenTTL:  accessTTL,
		RefreshTokenTTL: refreshTTL,
		RunMigrations:   parseBool("APP_RUN_MIGRATIONS", true),
		MigrationsPath:  parseString("APP_MIGRATIONS_PATH", "backend/migrations"),
		DBMaxOpenConns:  parseInt("APP_DB_MAX_OPEN_CONNS", 25),
		DBMaxIdleConns:  parseInt("APP_DB_MAX_IDLE_CONNS", 25),
		DBConnMaxIdle:   maxIdleTime,
		DBConnMaxLife:   maxLifetime,
	}

	if cfg.DatabaseURL == "" {
		return Config{}, errors.New("APP_DB_URL is required")
	}
	if cfg.JWTSecret == "" {
		return Config{}, errors.New("APP_JWT_SECRET is required")
	}

	if !filepath.IsAbs(cfg.MigrationsPath) {
		cfg.MigrationsPath = filepath.Clean(cfg.MigrationsPath)
	}

	return cfg, nil
}

func parseDuration(key, fallback string) (time.Duration, error) {
	value := parseString(key, fallback)
	d, err := time.ParseDuration(value)
	if err != nil {
		return 0, err
	}
	return d, nil
}

func parseBool(key string, fallback bool) bool {
	raw := os.Getenv(key)
	if raw == "" {
		return fallback
	}
	b, err := strconv.ParseBool(raw)
	if err != nil {
		return fallback
	}
	return b
}

func parseInt(key string, fallback int) int {
	raw := os.Getenv(key)
	if raw == "" {
		return fallback
	}
	n, err := strconv.Atoi(raw)
	if err != nil {
		return fallback
	}
	if n <= 0 {
		return fallback
	}
	return n
}

func parseString(key, fallback string) string {
	raw := os.Getenv(key)
	if raw == "" {
		return fallback
	}
	return raw
}
