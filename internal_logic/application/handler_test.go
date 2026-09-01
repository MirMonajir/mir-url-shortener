package application

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/MirMonajir/mir-url-shortener/internal_logic/infrastructure"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestShortenURL_Success(t *testing.T) {
	os.Setenv("SERVER_URL", "localhost:8080")
	defer os.Unsetenv("SERVER_URL")

	store := infrastructure.NewInMemoryStore()
	service := NewShortenerService(store)
	handler := NewHTTPHandler(service)

	// Create a mock Gin context
	w := httptest.NewRecorder()
	body := bytes.NewBufferString(`{"url": "https://www.google.com"}`)
	req, _ := http.NewRequest("POST", "/shortenurl", body)
	req.Header.Set("Content-Type", "application/json")

	c, _ := gin.CreateTestContext(w)
	c.Request = req
	handler.ShortenURL(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var response shortenResp
	_ = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NotEmpty(t, response.ShortURL)
	// Short URL should include the server URL prefix
	assert.Contains(t, response.ShortURL, "localhost:8080/")
}

func TestShortenURL_EmptyURL(t *testing.T) {
	store := infrastructure.NewInMemoryStore()
	service := NewShortenerService(store)
	handler := NewHTTPHandler(service)

	w := httptest.NewRecorder()
	body := bytes.NewBufferString(`{"url": ""}`)
	req, _ := http.NewRequest("POST", "/shortenurl", body)
	req.Header.Set("Content-Type", "application/json")

	c, _ := gin.CreateTestContext(w)
	c.Request = req
	handler.ShortenURL(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "error_type")
}

func TestShortenURL_InvalidURL_NoScheme(t *testing.T) {
	store := infrastructure.NewInMemoryStore()
	service := NewShortenerService(store)
	handler := NewHTTPHandler(service)

	w := httptest.NewRecorder()
	body := bytes.NewBufferString(`{"url": "www.example.com"}`)
	req, _ := http.NewRequest("POST", "/shortenurl", body)
	req.Header.Set("Content-Type", "application/json")

	c, _ := gin.CreateTestContext(w)
	c.Request = req
	handler.ShortenURL(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	var response map[string]interface{}
	_ = json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "invalid_url", response["error_type"])
}

func TestShortenURL_Localhost(t *testing.T) {
	store := infrastructure.NewInMemoryStore()
	service := NewShortenerService(store)
	handler := NewHTTPHandler(service)

	w := httptest.NewRecorder()
	body := bytes.NewBufferString(`{"url": "http://localhost:8080/api"}`)
	req, _ := http.NewRequest("POST", "/shortenurl", body)
	req.Header.Set("Content-Type", "application/json")

	c, _ := gin.CreateTestContext(w)
	c.Request = req
	handler.ShortenURL(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	var response map[string]interface{}
	_ = json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "invalid_url", response["error_type"])
}

func TestShortenAndResolve_EndToEnd(t *testing.T) {
	os.Setenv("SERVER_URL", "localhost:8080")
	defer os.Unsetenv("SERVER_URL")

	store := infrastructure.NewInMemoryStore()
	service := NewShortenerService(store)

	originalURL := "https://www.example.com"

	// Shorten the URL
	short, err := service.Shorten(originalURL)
	assert.NoError(t, err)
	assert.NotEmpty(t, short)
	// Extract just the code from the full URL (e.g., from "http://localhost:8080/abc123")
	assert.Contains(t, short, "localhost:8080/")

	// Extract the code from the URL
	code := short[len("http://localhost:8080/"):]
	assert.Len(t, code, 6)

	// Resolve it back
	resolved, err := service.Resolve(code)
	assert.NoError(t, err)
	assert.Equal(t, originalURL, resolved)
}

func TestRedirect_NotFound(t *testing.T) {
	store := infrastructure.NewInMemoryStore()
	service := NewShortenerService(store)

	// Try to resolve a URL that doesn't exist
	_, err := service.Resolve("nonexistent")
	assert.Error(t, err)
}

func TestMetrics(t *testing.T) {
	store := infrastructure.NewInMemoryStore()
	service := NewShortenerService(store)
	handler := NewHTTPHandler(service)

	// Shorten multiple URLs from different domains
	if _, err := service.Shorten("https://www.example.com/path1"); err != nil {
		t.Fatalf("Shorten failed: %v", err)
	}
	if _, err := service.Shorten("https://www.example.com/path2"); err != nil {
		t.Fatalf("Shorten failed: %v", err)
	}
	if _, err := service.Shorten("https://www.google.com/search"); err != nil {
		t.Fatalf("Shorten failed: %v", err)
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/appmetrics", nil)
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	handler.Metrics(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	_ = json.Unmarshal(w.Body.Bytes(), &response)
	assert.Contains(t, response, "top_domains")
}
