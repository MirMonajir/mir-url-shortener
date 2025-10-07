package model

import (
    "errors"
    "net/url"
    "strings"
)

type URL struct {
    OriginalUrl  string
    ShortenedUrl string
}

// NewURL constructor validates the original URL
func NewURL(originalUrl string) (*URL, error) {
    original = strings.TrimSpace(originalUrl)
    if original == "" {
        return nil, errors.New("The provided URL is empty, please provide a valid url")
    }
    // Basic validation
    parsed, err := url.Parse(original)
    if err != nil || parsed.Scheme == "" || parsed.Host == "" {
        return nil, errors.New("invalid URL format")
    }
    return &URL{Original: orig}, nil
}
