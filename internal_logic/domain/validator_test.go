package domain

import (
	"testing"
)

func TestURLValidator_ValidateEmptyURL(t *testing.T) {
	v := NewURLValidator()
	err := v.Validate("")
	if err == nil {
		t.Error("Expected error for empty URL, got nil")
	}
	if err.Type != ErrEmptyURL {
		t.Errorf("Expected ErrEmptyURL, got %v", err.Type)
	}
}

func TestURLValidator_ValidateURLTooShort(t *testing.T) {
	v := NewURLValidator()
	err := v.Validate("http://a")
	if err == nil {
		t.Error("Expected error for short URL, got nil")
	}
	if err.Type != ErrInvalidURL {
		t.Errorf("Expected ErrInvalidURL, got %v", err.Type)
	}
}

func TestURLValidator_ValidateURLTooLong(t *testing.T) {
	v := NewURLValidator()
	longURL := "https://example.com/" + string(make([]byte, MaxURLLength))
	err := v.Validate(longURL)
	if err == nil {
		t.Error("Expected error for long URL, got nil")
	}
	if err.Type != ErrInvalidURL {
		t.Errorf("Expected ErrInvalidURL, got %v", err.Type)
	}
}

func TestURLValidator_ValidateMissingScheme(t *testing.T) {
	v := NewURLValidator()
	err := v.Validate("www.example.com/path")
	if err == nil {
		t.Error("Expected error for URL without scheme, got nil")
	}
	if err.Type != ErrInvalidURL {
		t.Errorf("Expected ErrInvalidURL, got %v", err.Type)
	}
}

func TestURLValidator_ValidateInvalidScheme(t *testing.T) {
	v := NewURLValidator()
	err := v.Validate("ftp://example.com/path")
	if err == nil {
		t.Error("Expected error for invalid scheme, got nil")
	}
	if err.Type != ErrInvalidURL {
		t.Errorf("Expected ErrInvalidURL, got %v", err.Type)
	}
}

func TestURLValidator_ValidateMissingHost(t *testing.T) {
	v := NewURLValidator()
	err := v.Validate("https:///path")
	if err == nil {
		t.Error("Expected error for URL without host, got nil")
	}
	if err.Type != ErrInvalidURL {
		t.Errorf("Expected ErrInvalidURL, got %v", err.Type)
	}
}

func TestURLValidator_ValidateLocalhost(t *testing.T) {
	v := NewURLValidator()
	err := v.Validate("http://localhost:8080/path")
	if err == nil {
		t.Error("Expected error for localhost URL, got nil")
	}
	if err.Type != ErrInvalidURL {
		t.Errorf("Expected ErrInvalidURL, got %v", err.Type)
	}
}

func TestURLValidator_ValidateLoopbackIP(t *testing.T) {
	v := NewURLValidator()
	err := v.Validate("http://127.0.0.1/path")
	if err == nil {
		t.Error("Expected error for loopback IP, got nil")
	}
	if err.Type != ErrInvalidURL {
		t.Errorf("Expected ErrInvalidURL, got %v", err.Type)
	}
}

func TestURLValidator_ValidatePrivateIP(t *testing.T) {
	v := NewURLValidator()
	err := v.Validate("http://192.168.1.1/path")
	if err == nil {
		t.Error("Expected error for private IP, got nil")
	}
	if err.Type != ErrInvalidURL {
		t.Errorf("Expected ErrInvalidURL, got %v", err.Type)
	}
}

func TestURLValidator_ValidateValidURL(t *testing.T) {
	v := NewURLValidator()
	tests := []string{
		"https://www.example.com",
		"https://example.com/path/to/resource",
		"https://example.com/path?query=value",
		"http://example.com:8080/path",
	}

	for _, testURL := range tests {
		err := v.Validate(testURL)
		if err != nil {
			t.Errorf("Expected no error for valid URL %q, got %v", testURL, err)
		}
	}
}

func TestURLValidator_ValidateHTTPScheme(t *testing.T) {
	v := NewURLValidator()
	err := v.Validate("http://example.com/path")
	if err != nil {
		t.Errorf("Expected no error for http URL, got %v", err)
	}
}
