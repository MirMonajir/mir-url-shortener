package domain

import (
	"net"
	"net/url"
	"strings"
)

const (
	MaxURLLength = 2048
	MinURLLength = 10
)

// URLValidator validates URLs
type URLValidator struct{}

// NewURLValidator creates a new URL validator
func NewURLValidator() *URLValidator {
	return &URLValidator{}
}

// Validate performs comprehensive URL validation
func (v *URLValidator) Validate(originalURL string) *AppError {
	// Check for empty URL
	if strings.TrimSpace(originalURL) == "" {
		return NewAppError(ErrEmptyURL, "URL cannot be empty", 400)
	}

	// Check length
	if len(originalURL) < MinURLLength {
		return NewAppErrorWithDetails(
			ErrInvalidURL,
			"URL is too short",
			"Minimum length is "+string(rune(MinURLLength))+" characters",
			400,
		)
	}

	if len(originalURL) > MaxURLLength {
		return NewAppErrorWithDetails(
			ErrInvalidURL,
			"URL exceeds maximum length",
			"Maximum length is "+string(rune(MaxURLLength))+" characters",
			400,
		)
	}

	// Parse URL
	parsedURL, err := url.Parse(originalURL)
	if err != nil {
		return NewAppErrorWithDetails(
			ErrInvalidURL,
			"Invalid URL format",
			err.Error(),
			400,
		)
	}

	// Validate scheme
	if parsedURL.Scheme == "" {
		return NewAppError(
			ErrInvalidURL,
			"URL must include a scheme (http or https)",
			400,
		)
	}

	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return NewAppErrorWithDetails(
			ErrInvalidURL,
			"Only http and https schemes are allowed",
			"Provided scheme: "+parsedURL.Scheme,
			400,
		)
	}

	// Validate host
	if parsedURL.Host == "" {
		return NewAppError(
			ErrInvalidURL,
			"URL must include a host",
			400,
		)
	}

	// Prevent SSRF by blocking private/local IPs and reserved hosts
	if err := v.validateHost(parsedURL.Host); err != nil {
		return err
	}

	return nil
}

// validateHost checks if host is private or reserved
func (v *URLValidator) validateHost(host string) *AppError {
	// Remove port if present
	hostOnly := strings.Split(host, ":")[0]

	// Check for localhost
	if hostOnly == "localhost" {
		return NewAppError(
			ErrInvalidURL,
			"Localhost URLs are not allowed",
			400,
		)
	}

	// Try to parse as IP
	ip := net.ParseIP(hostOnly)
	if ip != nil {
		// Check for private/reserved IPs
		if ip.IsLoopback() {
			return NewAppError(ErrInvalidURL, "Loopback IP addresses are not allowed", 400)
		}
		if ip.IsPrivate() {
			return NewAppError(ErrInvalidURL, "Private IP addresses are not allowed", 400)
		}
		if ip.IsUnspecified() {
			return NewAppError(ErrInvalidURL, "Unspecified IP addresses are not allowed", 400)
		}
		if ip.IsMulticast() {
			return NewAppError(ErrInvalidURL, "Multicast IP addresses are not allowed", 400)
		}
		if ip.IsLinkLocalUnicast() {
			return NewAppError(ErrInvalidURL, "Link-local IP addresses are not allowed", 400)
		}
	}

	return nil
}
