package infrastructure

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

// Config holds all application configuration
type Config struct {
	Port      int
	Host      string
	ServerURL string
	GinMode   string
	LogLevel  string
}

// LoadConfig loads configuration from environment variables with validation
func LoadConfig() (*Config, error) {
	cfg := &Config{
		Port:     8080,
		Host:     "0.0.0.0",
		GinMode:  "debug",
		LogLevel: "info",
	}

	// Load PORT
	if portStr := os.Getenv("PORT"); portStr != "" {
		port, err := strconv.Atoi(portStr)
		if err != nil {
			return nil, fmt.Errorf("invalid PORT: %s - must be a number", portStr)
		}
		if port < 1 || port > 65535 {
			return nil, fmt.Errorf("invalid PORT: %d - must be between 1 and 65535", port)
		}
		cfg.Port = port
	}

	// Load HOST
	if host := os.Getenv("HOST"); host != "" {
		if err := validateHost(host); err != nil {
			return nil, fmt.Errorf("invalid HOST: %w", err)
		}
		cfg.Host = host
	}

	// Load SERVER_URL
	if serverURL := os.Getenv("SERVER_URL"); serverURL != "" {
		if err := validateServerURL(serverURL); err != nil {
			return nil, fmt.Errorf("invalid SERVER_URL: %w", err)
		}
		cfg.ServerURL = serverURL
	}

	// Load GIN_MODE
	if ginMode := os.Getenv("GIN_MODE"); ginMode != "" {
		if !isValidGinMode(ginMode) {
			return nil, fmt.Errorf("invalid GIN_MODE: %s - must be 'debug', 'release', or 'test'", ginMode)
		}
		cfg.GinMode = ginMode
	}

	// Load LOG_LEVEL
	if logLevel := os.Getenv("LOG_LEVEL"); logLevel != "" {
		if !isValidLogLevel(logLevel) {
			return nil, fmt.Errorf("invalid LOG_LEVEL: %s - must be 'debug', 'info', 'warn', or 'error'", logLevel)
		}
		cfg.LogLevel = logLevel
	}

	return cfg, nil
}

// validateHost checks if host is valid
func validateHost(host string) error {
	// Check if it's localhost or valid IP
	if host == "localhost" || host == "127.0.0.1" || host == "::1" {
		return nil
	}

	// Try to parse as IP address
	if ip := net.ParseIP(host); ip != nil {
		return nil
	}

	// Try to parse as hostname
	if !isValidHostname(host) {
		return fmt.Errorf("invalid hostname or IP: %s", host)
	}

	return nil
}

// isValidHostname validates a hostname
func isValidHostname(host string) bool {
	if len(host) > 253 {
		return false
	}

	labels := strings.Split(host, ".")
	for _, label := range labels {
		if len(label) == 0 || len(label) > 63 {
			return false
		}
		// Label must start and end with alphanumeric
		if !isAlphaNumeric(label[0]) || !isAlphaNumeric(label[len(label)-1]) {
			return false
		}
		// Label can contain hyphens in the middle
		for _, c := range label {
			if !isAlphaNumeric(byte(c)) && c != '-' {
				return false
			}
		}
	}
	return true
}

// isAlphaNumeric checks if a character is alphanumeric
func isAlphaNumeric(b byte) bool {
	return (b >= 'a' && b <= 'z') || (b >= 'A' && b <= 'Z') || (b >= '0' && b <= '9')
}

// validateServerURL validates the server URL format
func validateServerURL(url string) error {
	// Should be in format: hostname:port or just hostname
	parts := strings.Split(url, ":")
	if len(parts) > 2 {
		return fmt.Errorf("invalid format: too many colons")
	}

	hostname := parts[0]
	if err := validateHost(hostname); err != nil {
		return err
	}

	// If port is specified, validate it
	if len(parts) == 2 {
		port, err := strconv.Atoi(parts[1])
		if err != nil {
			return fmt.Errorf("invalid port in SERVER_URL: %s", parts[1])
		}
		if port < 1 || port > 65535 {
			return fmt.Errorf("port must be between 1 and 65535, got %d", port)
		}
	}

	return nil
}

// isValidGinMode checks if the Gin mode is valid
func isValidGinMode(mode string) bool {
	validModes := map[string]bool{
		"debug":   true,
		"release": true,
		"test":    true,
	}
	return validModes[mode]
}

// isValidLogLevel checks if the log level is valid
func isValidLogLevel(level string) bool {
	validLevels := map[string]bool{
		"debug": true,
		"info":  true,
		"warn":  true,
		"error": true,
	}
	return validLevels[level]
}
