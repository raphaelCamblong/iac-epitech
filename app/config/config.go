package config

import (
	"os"
	"strconv"
	"time"
)

// Config holds application configuration from environment.
type Config struct {
	Port           string
	DatabaseURL    string
	JWTSecret      string
	JWTExpiry      time.Duration
	LogLevel       string
	LogJSON        bool
	RateLimitRPS   float64
	RateLimitBurst int
}

// Load reads configuration from environment variables.
func Load() *Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	jwtExpiry := 15 * time.Minute
	if s := os.Getenv("JWT_EXPIRY_MINUTES"); s != "" {
		if m, err := strconv.Atoi(s); err == nil && m > 0 {
			jwtExpiry = time.Duration(m) * time.Minute
		}
	}
	rateRPS := 100.0
	if s := os.Getenv("RATE_LIMIT_RPS"); s != "" {
		if r, err := strconv.ParseFloat(s, 64); err == nil && r > 0 {
			rateRPS = r
		}
	}
	rateBurst := 50
	if s := os.Getenv("RATE_LIMIT_BURST"); s != "" {
		if b, err := strconv.Atoi(s); err == nil && b > 0 {
			rateBurst = b
		}
	}
	return &Config{
		Port:           port,
		DatabaseURL:    os.Getenv("DATABASE_URL"),
		JWTSecret:      os.Getenv("JWT_SECRET"),
		JWTExpiry:      jwtExpiry,
		LogLevel:       getEnv("LOG_LEVEL", "info"),
		LogJSON:        getEnv("LOG_JSON", "true") == "true",
		RateLimitRPS:   rateRPS,
		RateLimitBurst: rateBurst,
	}
}

func getEnv(key, defaultVal string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultVal
}
