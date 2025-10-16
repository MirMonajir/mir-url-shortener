package application

import (
	"errors"
	"os"
	"testing"

	"github.com/MirMonajir/mir-url-shortener/internal_logic/domain"
	"github.com/stretchr/testify/assert"
)

func TestShorten_ReturnSuccess(t *testing.T) {
	mockRepo := new(MockStorage)
	svc := NewShortenerService(mockRepo)

	originalURL := "https://google.com"
	urlObj, _ := domain.NewURL(originalURL)

	mockRepo.On("Save", urlObj).Return("xyz123", nil)

	// Set BASE_URL env variable for the test
	os.Setenv("SERVER_URL", "localhost:8080")
	defer os.Unsetenv("SERVER_URL") // clean up after test

	shortURL, err := svc.Shorten(originalURL)

	assert.NoError(t, err)
	assert.Equal(t, "http://localhost:8080/xyz123", shortURL)
	mockRepo.AssertExpectations(t)
}

func TestShorten_ReturnInvalidURL(t *testing.T) {
	mockRepo := new(MockStorage)
	svc := NewShortenerService(mockRepo)

	invalidURL := "ht!tp:/bad-url"

	shortURL, err := svc.Shorten(invalidURL)

	assert.Error(t, err)
	assert.Empty(t, shortURL)
}

func TestShorten_ReturnSaveError(t *testing.T) {
	mockRepo := new(MockStorage)
	svc := NewShortenerService(mockRepo)

	originalURL := "https://example.com"
	urlObj, _ := domain.NewURL(originalURL)

	mockRepo.On("Save", urlObj).Return("", errors.New("save failed"))

	shortURL, err := svc.Shorten(originalURL)

	assert.Error(t, err)
	assert.Empty(t, shortURL)
	mockRepo.AssertExpectations(t)
}

func TestResolve_ReturnSuccess(t *testing.T) {
	mockRepo := new(MockStorage)
	svc := NewShortenerService(mockRepo)

	mockRepo.On("Get", "abc123").Return("https://google.com", nil)

	result, err := svc.Resolve("abc123")

	assert.NoError(t, err)
	assert.Equal(t, "https://google.com", result)
}

func TestResolve_ReturnNotFound(t *testing.T) {
	mockRepo := new(MockStorage)
	svc := NewShortenerService(mockRepo)

	mockRepo.On("Get", "notfound").Return("", errors.New("not found"))

	result, err := svc.Resolve("notfound")

	assert.Error(t, err)
	assert.Empty(t, result)
}

func TestTopDomains(t *testing.T) {
	mockRepo := new(MockStorage)
	svc := NewShortenerService(mockRepo)

	expected := map[string]int{
		"google.com": 5,
		"openai.com": 3,
	}

	mockRepo.On("TopDomains", 2).Return(expected)

	result := svc.TopDomains(2)

	assert.Equal(t, expected, result)
}
