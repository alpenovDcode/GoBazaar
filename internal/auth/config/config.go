package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

// Config contains auth service configuration
type Config struct {
	// HTTP server
	HTTPHost string
	HTTPPort string

	// Database
	DatabaseURL  string
	DatabaseHost string
	DatabasePort string
	DatabaseUser string
	DatabasePass string
	DatabaseName string

	// JWT
	JWTSecret           string
	JWTExpirationHours  int
	RefreshTokenExpDays int

	// Redis
	RedisHost     string
	RedisPort     string
	RedisPassword string
	RedisDB       int

	// NATS
	NATSHost string
	NATSPort string

	// Security
	BCryptCost        int
	MaxLoginAttempts  int
	LockoutDuration   time.Duration
	PasswordMinLength int

	// Environment
	Environment string
	LogLevel    string
	Debug       bool
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	cfg := &Config{
		// Default values
		HTTPHost:            getEnv("AUTH_HTTP_HOST", "0.0.0.0"),
		HTTPPort:            getEnv("AUTH_HTTP_PORT", "8080"),
		DatabaseHost:        getEnv("DB_HOST", "localhost"),
		DatabasePort:        getEnv("DB_PORT", "5432"),
		DatabaseUser:        getEnv("DB_USER", "postgres"),
		DatabasePass:        getEnv("DB_PASSWORD", "postgres"),
		DatabaseName:        getEnv("DB_NAME", "gobazaar_auth"),
		JWTSecret:           getEnv("JWT_SECRET", "your-secret-key-change-in-production"),
		JWTExpirationHours:  getEnvAsInt("JWT_EXPIRATION_HOURS", 24),
		RefreshTokenExpDays: getEnvAsInt("REFRESH_TOKEN_EXP_DAYS", 30),
		RedisHost:           getEnv("REDIS_HOST", "localhost"),
		RedisPort:           getEnv("REDIS_PORT", "6379"),
		RedisPassword:       getEnv("REDIS_PASSWORD", ""),
		RedisDB:             getEnvAsInt("REDIS_DB", 0),
		NATSHost:            getEnv("NATS_HOST", "localhost"),
		NATSPort:            getEnv("NATS_PORT", "4222"),
		BCryptCost:          getEnvAsInt("BCRYPT_COST", 12),
		MaxLoginAttempts:    getEnvAsInt("MAX_LOGIN_ATTEMPTS", 5),
		LockoutDuration:     getEnvAsDuration("LOCKOUT_DURATION", "15m"),
		PasswordMinLength:   getEnvAsInt("PASSWORD_MIN_LENGTH", 8),
		Environment:         getEnv("ENVIRONMENT", "development"),
		LogLevel:            getEnv("LOG_LEVEL", "info"),
		Debug:               getEnvAsBool("DEBUG", true),
	}

	// Build DATABASE_URL if not provided
	if cfg.DatabaseURL == "" {
		cfg.DatabaseURL = fmt.Sprintf(
			"postgresql://%s:%s@%s:%s/%s?sslmode=disable",
			cfg.DatabaseUser,
			cfg.DatabasePass,
			cfg.DatabaseHost,
			cfg.DatabasePort,
			cfg.DatabaseName,
		)
	} else {
		cfg.DatabaseURL = getEnv("DATABASE_URL", cfg.DatabaseURL)
	}

	// Validate critical parameters
	if err := cfg.validate(); err != nil {
		return nil, fmt.Errorf("configuration validation failed: %w", err)
	}

	return cfg, nil
}

// validate checks configuration correctness
func (c *Config) validate() error {
	if c.JWTSecret == "your-secret-key-change-in-production" {
		return fmt.Errorf("JWT_SECRET must be changed from default value")
	}

	if len(c.JWTSecret) < 32 {
		return fmt.Errorf("JWT_SECRET must be at least 32 characters long")
	}

	if c.JWTExpirationHours < 1 {
		return fmt.Errorf("JWT_EXPIRATION_HOURS must be at least 1")
	}

	if c.RefreshTokenExpDays < 1 {
		return fmt.Errorf("REFRESH_TOKEN_EXP_DAYS must be at least 1")
	}

	if c.BCryptCost < 10 || c.BCryptCost > 15 {
		return fmt.Errorf("BCRYPT_COST must be between 10 and 15")
	}

	if c.MaxLoginAttempts < 1 {
		return fmt.Errorf("MAX_LOGIN_ATTEMPTS must be at least 1")
	}

	if c.PasswordMinLength < 8 {
		return fmt.Errorf("PASSWORD_MIN_LENGTH must be at least 8")
	}

	return nil
}

// IsProduction returns true if environment is production
func (c *Config) IsProduction() bool {
	return c.Environment == "production"
}

// IsDevelopment returns true if environment is development
func (c *Config) IsDevelopment() bool {
	return c.Environment == "development"
}

// GetAddress returns HTTP server address
func (c *Config) GetAddress() string {
	return c.HTTPHost + ":" + c.HTTPPort
}

// GetDatabaseURL returns database URL
func (c *Config) GetDatabaseURL() string {
	return c.DatabaseURL
}

// GetRedisAddress returns Redis address
func (c *Config) GetRedisAddress() string {
	return c.RedisHost + ":" + c.RedisPort
}

// GetNATSURL returns NATS URL
func (c *Config) GetNATSURL() string {
	return "nats://" + c.NATSHost + ":" + c.NATSPort
}

// Helper functions for reading environment variables

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvAsBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}

func getEnvAsDuration(key string, defaultValue string) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	duration, _ := time.ParseDuration(defaultValue)
	return duration
}
