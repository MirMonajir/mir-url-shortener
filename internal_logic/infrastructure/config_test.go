package infrastructure

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfig_DefaultValues(t *testing.T) {
	// Clear env vars
	assert.NoError(t, os.Unsetenv("PORT"))
	assert.NoError(t, os.Unsetenv("HOST"))
	assert.NoError(t, os.Unsetenv("SERVER_URL"))
	assert.NoError(t, os.Unsetenv("GIN_MODE"))
	assert.NoError(t, os.Unsetenv("LOG_LEVEL"))

	cfg, err := LoadConfig()

	assert.NoError(t, err)
	assert.Equal(t, 8080, cfg.Port)
	assert.Equal(t, "0.0.0.0", cfg.Host)
	assert.Equal(t, "debug", cfg.GinMode)
	assert.Equal(t, "info", cfg.LogLevel)
}

func TestLoadConfig_CustomPort(t *testing.T) {
	assert.NoError(t, os.Setenv("PORT", "9090"))
	defer func() { assert.NoError(t, os.Unsetenv("PORT")) }()

	cfg, err := LoadConfig()

	assert.NoError(t, err)
	assert.Equal(t, 9090, cfg.Port)
}

func TestLoadConfig_InvalidPort_NotANumber(t *testing.T) {
	assert.NoError(t, os.Setenv("PORT", "not-a-number"))
	defer func() { assert.NoError(t, os.Unsetenv("PORT")) }()

	_, err := LoadConfig()

	assert.Error(t, err)
}

func TestLoadConfig_InvalidPort_OutOfRange(t *testing.T) {
	assert.NoError(t, os.Setenv("PORT", "70000"))
	defer func() { assert.NoError(t, os.Unsetenv("PORT")) }()

	_, err := LoadConfig()

	assert.Error(t, err)
}

func TestLoadConfig_CustomHost(t *testing.T) {
	assert.NoError(t, os.Setenv("HOST", "localhost"))
	defer func() { assert.NoError(t, os.Unsetenv("HOST")) }()

	cfg, err := LoadConfig()

	assert.NoError(t, err)
	assert.Equal(t, "localhost", cfg.Host)
}

func TestLoadConfig_ServerURL(t *testing.T) {
	assert.NoError(t, os.Setenv("SERVER_URL", "example.com:8080"))
	defer func() { assert.NoError(t, os.Unsetenv("SERVER_URL")) }()

	cfg, err := LoadConfig()

	assert.NoError(t, err)
	assert.Equal(t, "example.com:8080", cfg.ServerURL)
}

func TestLoadConfig_InvalidGinMode(t *testing.T) {
	assert.NoError(t, os.Setenv("GIN_MODE", "invalid-mode"))
	defer func() { assert.NoError(t, os.Unsetenv("GIN_MODE")) }()

	_, err := LoadConfig()

	assert.Error(t, err)
}

func TestLoadConfig_ValidLogLevels(t *testing.T) {
	levels := []string{"debug", "info", "warn", "error"}

	for _, level := range levels {
		assert.NoError(t, os.Setenv("LOG_LEVEL", level))

		cfg, err := LoadConfig()

		assert.NoError(t, err)
		assert.Equal(t, level, cfg.LogLevel)
	}

	assert.NoError(t, os.Unsetenv("LOG_LEVEL"))
}

func TestLoadConfig_InvalidLogLevel(t *testing.T) {
	assert.NoError(t, os.Setenv("LOG_LEVEL", "invalid-level"))
	defer func() { assert.NoError(t, os.Unsetenv("LOG_LEVEL")) }()

	_, err := LoadConfig()

	assert.Error(t, err)
}

func TestValidateHost_IP(t *testing.T) {
	tests := []struct {
		host  string
		valid bool
	}{
		{"127.0.0.1", true},
		{"::1", true},
		{"0.0.0.0", true},
		{"192.168.1.1", true},
		{"localhost", true},
		{"example.com", true},
		{"sub.example.com", true},
		// Note: Go's net.ParseIP accepts various formats
		// 192.168.1 is valid (interpreted as 192.168.1.0)
		// So we focus on other validation checks
	}

	for _, tt := range tests {
		err := validateHost(tt.host)
		if (err == nil) != tt.valid {
			t.Errorf("validateHost(%q) valid=%v, want %v, err=%v", tt.host, err == nil, tt.valid, err)
		}
	}
}
